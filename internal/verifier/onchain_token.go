package verifier

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/rpc"
	"github.com/ardmere/ardmere/internal/walletzip"
)

// OnchainBalanceToken verifies ERC20/BEP20 balanceOf for supported HotCold rows.
type OnchainBalanceToken struct {
	Exchange        string
	WalletCSV       walletzip.File
	WalletZipPath   string
	WalletZipSha256 string
	BapiSha256      string
	SnapshotID      string
	RPC             *rpc.Client
	Ledger          *rpc.LedgerClient
	Concurrency     int // default 4
}

func (v OnchainBalanceToken) ledgerClient() *rpc.LedgerClient {
	if v.Ledger != nil {
		return v.Ledger
	}
	return rpc.NewLedger()
}

type decimalsCache struct {
	mu    sync.Mutex
	byKey map[string]int
}

func newDecimalsCache() *decimalsCache {
	return &decimalsCache{byKey: map[string]int{}}
}

func (v OnchainBalanceToken) Run(ctx context.Context) Verification {
	if v.Concurrency <= 0 {
		v.Concurrency = 4
	}

	out := Verification{
		VerifierID:     "onchain-balance-token",
		Version:        "2.0",
		SnapshotID:     v.SnapshotID,
		VerifiedAt:     time.Now().UTC(),
		InputArtifacts: []string{v.BapiSha256, v.WalletZipSha256},
	}

	var allRows []walletzip.Row
	csvFile := v.WalletCSV
	if csvFile == 0 {
		csvFile = walletzip.HotCold
	}
	_, err := walletzip.ForEachRow(v.WalletZipPath, csvFile, func(r walletzip.Row) error {
		allRows = append(allRows, r)
		return nil
	})
	if err != nil {
		out.Verdict = VerdictFail
		out.Reason = fmt.Sprintf("read HotCold csv: %v", err)
		return out
	}

	var supported []walletzip.Row
	unsupportedPairs := map[string]int{}
	for _, r := range allRows {
		key := r.Coin + "|" + r.Network
		if _, ok := loadTokenSupportedFor(v.Exchange)[key]; !ok {
			if isEVMTokenCandidate(r.Coin) {
				unsupportedPairs[key]++
			}
			continue
		}
		supported = append(supported, r)
	}

	if len(supported) == 0 {
		out.Verdict = VerdictPartial
		out.Coverage = 0
		out.Reason = "no HotCold rows in supported token (coin,network) set"
		out.Findings = unsupportedFindings(unsupportedPairs)
		return out
	}

	var sampleNote string
	supported, sampleNote = maybeSubsampleOnchainRows(v.Exchange, supported)
	if sampleNote != "" {
		out.Findings = append(out.Findings, Finding{
			Subject: "onchain-sample",
			Field:   "token_rows",
			Status:  VerdictPass,
			Note:    sampleNote,
		})
	}

	type result struct {
		row            walletzip.Row
		actual         decimal.Decimal
		used           string
		matchedHeight  int64
		boundaryNote   string
		components     map[string]string
		err            error
	}
	jobs := make(chan walletzip.Row, len(supported))
	results := make(chan result, len(supported))
	decCache := newDecimalsCache()
	sonicCache := newSonicSFCCache()
	var wg sync.WaitGroup
	for i := 0; i < v.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for r := range jobs {
				actual, used, comps, err := v.tokenBalanceAtHeight(ctx, r, decCache, sonicCache)
				res := result{row: r, actual: actual, used: used, matchedHeight: r.Height, components: comps, err: err}
				if err == nil {
					if !balanceWithinTolerance(actual, r.Balance) {
						if alt, u, h, note, ok := v.tryTokenHeightOffsets(ctx, r, decCache, sonicCache); ok {
							res.actual = alt
							res.used = u
							res.matchedHeight = h
							res.boundaryNote = note
						}
					}
				}
				results <- res
			}
		}()
	}
	for _, r := range supported {
		jobs <- r
	}
	close(jobs)
	go func() { wg.Wait(); close(results) }()

	var passCount, warnCount, failCount, unverifiableCount, errCount int
	for res := range results {
		if res.err != nil {
			errCount++
			out.Findings = append(out.Findings, Finding{
				Subject: res.row.Address,
				Field:   res.row.Coin + "@" + res.row.Network,
				Status:  VerdictWarn,
				Note:    fmt.Sprintf("rpc error (provider=%s): %v", res.used, res.err),
			})
			continue
		}

		claim := res.row.Balance
		diff := res.actual.Sub(claim)
		within, surplus := classifyBalanceMismatch(res.actual, claim)
		if within {
			passCount++
			continue
		}

		if unv, note := btcEthNativeCustodyUnverifiable(res.row, res.actual); unv {
			unverifiableCount++
			components := res.components
			if components == nil {
				components = map[string]string{}
			}
			spec := loadTokenSupportedFor(v.Exchange)[res.row.Coin+"|"+res.row.Network]
			if spec.Contract != "" {
				components["token"] = spec.Contract
			}
			out.Findings = append(out.Findings, Finding{
				Subject:    res.row.Address,
				Field:      res.row.Coin + "@" + res.row.Network + "#" + fmt.Sprintf("%d", res.matchedHeight),
				Claim:      claim.String(),
				Actual:     res.actual.String(),
				Status:     VerdictUnverifiable,
				Note:       note,
				Components: components,
			})
			continue
		}

		if omnibus, note := stablecoinEthOmnibusMismatch(res.row, res.actual); omnibus {
			warnCount++
			spec := loadTokenSupportedFor(v.Exchange)[res.row.Coin+"|"+res.row.Network]
			components := res.components
			if components == nil {
				components = map[string]string{}
			}
			if spec.Contract != "" {
				components["token"] = spec.Contract
				components["mode"] = "erc20"
			}
			if res.boundaryNote != "" {
				note = res.boundaryNote + "; " + note
			}
			out.Findings = append(out.Findings, Finding{
				Subject:    res.row.Address,
				Field:      res.row.Coin + "@" + res.row.Network + "#" + fmt.Sprintf("%d", res.matchedHeight),
				Claim:      claim.String(),
				Actual:     res.actual.String(),
				Status:     VerdictWarn,
				Note:       note + fmt.Sprintf("; delta=%s (provider=%s)", diff.Abs().String(), res.used),
				Components: components,
			})
			continue
		}

		spec := loadTokenSupportedFor(v.Exchange)[res.row.Coin+"|"+res.row.Network]
		note := fmt.Sprintf("on-chain balance != csv claim by %s (provider=%s)", diff.Abs().String(), res.used)
		if surplus {
			note = "chain observed > CSV row allocation; " + note
		}
		if res.boundaryNote != "" {
			note = res.boundaryNote + "; " + note
		}
		if res.row.Coin == "S" && res.row.Network == "SONIC" {
			if diag := sonicBalanceNote(claim, res.actual); diag != "" {
				note = diag + "; " + note
			}
		}
		if res.row.Coin == "POL" && res.row.Network == "ETH" {
			if diag := polEthBalanceNote(claim, res.actual); diag != "" {
				note = diag + "; " + note
			}
		}
		if diag := tronSnapshotMismatchNote(res.row.Coin, res.row.Network); diag != "" {
			note = diag + "; " + note
		}
		if diag := fevmSnapshotNote(res.row.Coin, res.row.Network); diag != "" {
			note = diag + "; " + note
		}
		components := res.components
		if components == nil {
			components = map[string]string{}
		}
		if res.row.Coin == "S" && res.row.Network == "SONIC" {
			// mode/liquid/staked already set by sonicAccountedBalance
		} else if res.row.Coin == "POL" && res.row.Network == "ETH" {
			// pol/matic components already set by polEthAccounted
		} else if spec.Native {
			components["mode"] = "native_balance"
		} else if spec.Net == rpc.NetTron && spec.Contract != "" {
			components["token"] = spec.Contract
			components["mode"] = "trc20"
		} else if spec.Contract != "" {
			components["token"] = spec.Contract
		}

		tronSnap := tronSnapshotMismatchNote(res.row.Coin, res.row.Network) != ""
		fevmSnap := fevmSnapshotNote(res.row.Coin, res.row.Network) != ""
		status := VerdictFail
		if surplus || tronSnap || fevmSnap {
			status = VerdictWarn
		}
		if surplus || tronSnap || fevmSnap {
			warnCount++
		} else {
			failCount++
		}
		out.Findings = append(out.Findings, Finding{
			Subject:    res.row.Address,
			Field:      res.row.Coin + "@" + res.row.Network + "#" + fmt.Sprintf("%d", res.matchedHeight),
			Claim:      claim.String(),
			Actual:     res.actual.String(),
			Status:     status,
			Note:       note,
			Components: components,
		})
	}

	out.Coverage = float64(passCount+failCount) / float64(len(allRows))
	out.Findings = append(out.Findings, unsupportedFindings(unsupportedPairs)...)
	out.Findings = append(out.Findings, Finding{
		Subject: "summary",
		Field:   "verified_pairs",
		Actual: fmt.Sprintf("%d pass / %d warn / %d fail / %d unverifiable / %d rpc-error -- of %d supported token rows / %d total HotCold rows",
			passCount, warnCount, failCount, unverifiableCount, errCount, len(supported), len(allRows)),
		Status: VerdictPass,
	})

	switch {
	case failCount > 0:
		out.Verdict = VerdictFail
	case unverifiableCount > 0:
		out.Verdict = VerdictWarn
		out.Reason = "some supported rows cannot be verified on-chain (e.g. native BTC custody on ETH network row)"
	case warnCount > 0:
		out.Verdict = VerdictWarn
		out.Reason = "some rows show chain-observed balance above CSV (surplus), not reserve shortfall"
	case errCount > len(supported)/4:
		out.Verdict = VerdictWarn
		out.Reason = "more than 25% of supported token rows failed to query rpc"
	case len(supported) < len(allRows):
		out.Verdict = VerdictPartial
	default:
		out.Verdict = VerdictPass
	}
	return out
}

func (v OnchainBalanceToken) tokenBalanceAtHeight(ctx context.Context, r walletzip.Row, decCache *decimalsCache, sonicCache *sonicSFCCache) (decimal.Decimal, string, map[string]string, error) {
	spec := loadTokenSupportedFor(v.Exchange)[r.Coin+"|"+r.Network]
	if spec.Net == rpc.Network("STARKNET") && spec.Native {
		bal, used, err := v.ledgerClient().StarknetNativeETHAtBlock(ctx, r.Address, r.Height)
		if err != nil {
			return decimal.Zero, used, nil, err
		}
		return decimal.NewFromBigInt(bal, -int32(spec.Decimals)), used, map[string]string{"mode": "starknet_weth"}, nil
	}
	if r.Coin == "S" && r.Network == "SONIC" {
		accounted, err := v.sonicAccountedBalance(ctx, r.Address, r.Height, sonicCache)
		if err != nil {
			return decimal.Zero, accounted.used, nil, err
		}
		total := accounted.liquid.Add(accounted.staked)
		components := map[string]string{
			"mode":       "native_plus_sfc",
			"liquid":     accounted.liquid.String(),
			"staked":     accounted.staked.String(),
			"validators": fmt.Sprintf("%d", accounted.scanned),
		}
		for k, v := range accounted.byValidator {
			components[k] = v
		}
		return total, accounted.used, components, nil
	}
	if r.Coin == "POL" && r.Network == "ETH" {
		total, used, components, err := polEthAccounted(ctx, v.RPC, r.Address, r.Height)
		return total, used, components, err
	}
	if spec.Native {
		if spec.Net == rpc.NetTron {
			bal, used, err := v.RPC.TronNativeBalance(ctx, r.Address, r.Height)
			if err != nil {
				return decimal.Zero, used, nil, err
			}
			return decimal.NewFromBigInt(bal, -int32(spec.Decimals)), used, map[string]string{"mode": "native_balance"}, nil
		}
		actual, used, err := nativeBalanceDecimal(ctx, v.RPC, spec.Net, r.Address, r.Height, spec.Decimals)
		return actual, used, map[string]string{"mode": "native_balance"}, err
	}
	if spec.Net == rpc.NetTron && spec.Contract != "" {
		bal, used, err := v.RPC.TRC20BalanceOf(ctx, spec.Contract, r.Address, r.Height)
		if err != nil {
			return decimal.Zero, used, nil, err
		}
		return decimal.NewFromBigInt(bal, -int32(spec.Decimals)), used, map[string]string{"token": spec.Contract, "mode": "trc20"}, nil
	}

	decimals := spec.Decimals
	if decimals == 0 {
		var err error
		decimals, err = v.tokenDecimals(ctx, spec, r.Height, decCache)
		if err != nil {
			return decimal.Zero, "", nil, err
		}
	}
	bal, used, err := erc20BalanceOf(ctx, v.RPC, spec.Net, spec.Contract, r.Address, r.Height)
	if err != nil {
		return decimal.Zero, used, nil, err
	}
	return decimal.NewFromBigInt(bal, -int32(decimals)), used, map[string]string{"token": spec.Contract}, nil
}

func (v OnchainBalanceToken) tokenBalance(ctx context.Context, r walletzip.Row, decCache *decimalsCache) (decimal.Decimal, string, error) {
	actual, used, _, err := v.tokenBalanceAtHeight(ctx, r, decCache, newSonicSFCCache())
	return actual, used, err
}

func (v OnchainBalanceToken) tokenDecimals(ctx context.Context, spec tokenSpec, height int64, decCache *decimalsCache) (int, error) {
	key := string(spec.Net) + "|" + spec.Contract
	decCache.mu.Lock()
	if d, ok := decCache.byKey[key]; ok {
		decCache.mu.Unlock()
		return d, nil
	}
	decCache.mu.Unlock()

	raw, used, err := v.RPC.CallContract(ctx, spec.Net, spec.Contract, selector("decimals()"), height)
	if err != nil {
		return 0, fmt.Errorf("decimals %s (provider=%s): %w", spec.Contract, used, err)
	}
	if len(raw) < 32 {
		return 0, fmt.Errorf("short decimals response for %s", spec.Contract)
	}
	d := int(new(big.Int).SetBytes(raw[len(raw)-32:]).Int64())
	decCache.mu.Lock()
	decCache.byKey[key] = d
	decCache.mu.Unlock()
	return d, nil
}

func erc20BalanceOf(ctx context.Context, c *rpc.Client, net rpc.Network, token, holder string, height int64) (*big.Int, string, error) {
	addr, err := rpc.EncodeHexAddress(holder)
	if err != nil {
		return nil, "", err
	}
	data := append(selector("balanceOf(address)"), rpc.To32Bytes(addr)...)
	raw, used, err := c.CallContract(ctx, net, token, data, height)
	if err != nil {
		return nil, used, err
	}
	if len(raw) < 32 {
		return nil, used, fmt.Errorf("short balanceOf response: %d bytes", len(raw))
	}
	return new(big.Int).SetBytes(raw[:32]), used, nil
}

func isEVMTokenCandidate(coin string) bool {
	switch coin {
	case "USDT", "USDC", "FDUSD", "USD1", "DAI", "TUSD", "BUSD", "USDE", "U",
		"ETH", "BTC", "SHIB", "LINK", "UNI", "AAVE", "PEPE", "ARB", "OP", "S", "POL",
		"PENDLE", "CRV", "ENA", "MASK", "GRT", "1INCH", "CAKE", "HFT", "SSV", "WLFI", "ASTER", "FORM",
		"PAXG", "CHR", "CHZ", "ENJ", "DOGE", "XRP", "RLUSD",
		"DOT", "NEAR", "HBAR", "SOL", "LTC", "BCH", "ZEC", "BNB":
		return true
	default:
		return false
	}
}

// StakeHubSummary is the aggregated BSC Stake Hub position for one delegator.
type StakeHubSummary struct {
	Staked  decimal.Decimal
	Locked  decimal.Decimal
	Scanned int
	Used    string
	Note    string
}

// StakeHubAccounted exposes BSC Stake Hub aggregation for diagnostics.
func StakeHubAccounted(ctx context.Context, rpcClient *rpc.Client, delegator string, height int64) (StakeHubSummary, error) {
	v := OnchainBalanceHot{RPC: rpcClient, Concurrency: 1}
	out, err := v.bscStakeHubAccounted(ctx, delegator, height, newBSCStakeHubCache())
	if err != nil {
		return StakeHubSummary{}, err
	}
	return StakeHubSummary{
		Staked:  out.staked,
		Locked:  out.locked,
		Scanned: out.scanned,
		Used:    out.used,
		Note:    out.scanNote,
	}, nil
}

// KeyAddresses from PR01JUN26 regression set.
var KeyAddresses = []struct {
	Label   string
	Address string
	Height  int64
	Coin    string
}{
	{"eth-deposit", "0x32e11a20337ebc79abd0eeab2d91bafbd9591149", 25218797, "ETH"},
	{"bnb-stakehub", "0x86523c87c8ec98c7539e2c58cd813ee9d1a08d96", 101590091, "BNB"},
	{"bnb-stakehub-2", "0xbf83d18a46325acb7d8f40a462d23a92f467ed7a", 101590091, "BNB"},
}

func SortFindingsByStatus(findings []Finding) []Finding {
	out := append([]Finding(nil), findings...)
	sort.Slice(out, func(i, j int) bool {
		return out[i].Subject < out[j].Subject
	})
	return out
}
