package verifier

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/ardmere/ardmere/internal/rpc"
	"github.com/ardmere/ardmere/internal/walletzip"
)

// OnchainBalanceHot verifies the accounted native (ETH on Ethereum, BNB on BSC)
// balance of every HotCold row whose (coin, network) pair we currently support
// against public RPC nodes, at the snapshot's declared block Height.
//
// V2 is staking-aware for the two native false-positive paths observed in
// PR01JUN26:
//   - BNB on BSC includes Stake Hub pooled and unbonding balances.
//   - ETH deposit-likely rows are WARN instead of liquid-only FAIL until an
//     execution-layer indexer / beacon API is wired in.
//
// Out of scope:
//   - ERC20 balanceOf calls (USDT/USDC/DAI/etc.)
//   - Non-EVM chains (BTC/SOL/TRX/APT/...)
//   - Deposit CSV (sampled separately by `OnchainBalanceDeposit`)
//
// As such, the verdict is intentionally PARTIAL when there are HotCold rows
// we can't verify (which there always are today). `coverage` reflects the
// fraction of HotCold rows we did verify.
type OnchainBalanceHot struct {
	Exchange        string
	WalletCSV       walletzip.File
	WalletZipPath   string
	WalletZipSha256 string
	BapiSha256      string // summary artifact sha256 (legacy field name)
	SnapshotID      string
	RPC             *rpc.Client
	ETHDeposits     ETHDepositIndexer
	Concurrency     int // default 8
}

// supported pairs we know how to verify natively (no contract calls needed).
// Loaded from config/exchanges/<exchange>/onchain.json (see onchain_native.go).

func (v OnchainBalanceHot) Run(ctx context.Context) Verification {
	if v.Concurrency <= 0 {
		v.Concurrency = 8
	}
	if v.ETHDeposits == nil {
		v.ETHDeposits = NewEtherscanETHDepositIndexer("")
	}

	out := Verification{
		VerifierID:     "onchain-balance-hot",
		Version:        "2.1",
		SnapshotID:     v.SnapshotID,
		VerifiedAt:     time.Now().UTC(),
		InputArtifacts: []string{v.BapiSha256, v.WalletZipSha256},
	}

	// Read all HotCold rows into memory (it's only ~10^3 rows / 110 KB).
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

	// Partition into supported vs unsupported.
	var supported []walletzip.Row
	unsupportedPairs := map[string]int{}
	for _, r := range allRows {
		key := r.Coin + "|" + r.Network
		if _, ok := loadNativeSupported(v.Exchange)[key]; ok {
			supported = append(supported, r)
		} else {
			unsupportedPairs[key]++
		}
	}

	if len(supported) == 0 {
		out.Verdict = VerdictPartial
		out.Coverage = 0
		out.Reason = "no HotCold rows in supported native (coin,network) set yet (V1: ETH|ETH, BNB|BSC)"
		out.Findings = unsupportedFindings(unsupportedPairs)
		return out
	}

	var sampleNote string
	supported, sampleNote = maybeSubsampleOnchainRows(v.Exchange, supported)
	if sampleNote != "" {
		out.Findings = append(out.Findings, Finding{
			Subject: "onchain-sample",
			Field:   "native_rows",
			Status:  VerdictPass,
			Note:    sampleNote,
		})
	}

	// Concurrent balance lookups.
	type result struct {
		row       walletzip.Row
		accounted accountedBalance
		used      string
		err       error
	}
	jobs := make(chan walletzip.Row, len(supported))
	results := make(chan result, len(supported))
	var wg sync.WaitGroup
	stakeCache := newBSCStakeHubCache()
	ethCache := newETHDepositCache()
	for i := 0; i < v.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for r := range jobs {
				accounted, err := v.accountedBalance(ctx, r, stakeCache, ethCache)
				if err != nil {
					results <- result{row: r, err: err, used: accounted.used}
					continue
				}
				results <- result{row: r, accounted: accounted, used: accounted.used}
			}
		}()
	}
	for _, r := range supported {
		jobs <- r
	}
	close(jobs)
	go func() { wg.Wait(); close(results) }()

	var (
		passCount           int
		toleranceNoiseCount int
		warnCount           int
		failCount           int
		errCount            int
	)

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
		matchedHeight := res.row.Height
		boundaryNote := ""
		if !balanceWithinTolerance(res.accounted.total, res.row.Balance) {
			if alt, h, note, ok := v.tryAccountedHeightWindow(ctx, res.row, stakeCache, ethCache); ok {
				res.accounted = alt
				matchedHeight = h
				boundaryNote = note
			}
		}
		diff := res.accounted.total.Sub(res.row.Balance).Abs()
		if balanceWithinTolerance(res.accounted.total, res.row.Balance) {
			passCount++
			if !diff.IsZero() {
				toleranceNoiseCount++ // tracked but no per-row finding emitted
			}
			continue
		}
		if res.accounted.incomplete {
			warnCount++
			out.Findings = append(out.Findings, Finding{
				Subject:    res.row.Address,
				Field:      res.row.Coin + "@" + res.row.Network + "#" + fmt.Sprintf("%d", res.row.Height),
				Claim:      res.row.Balance.String(),
				Actual:     res.accounted.total.String(),
				Status:     VerdictWarn,
				Note:       fmt.Sprintf("%s; observed delta=%s (provider=%s)", res.accounted.note, diff.String(), res.used),
				Components: res.accounted.components,
			})
			continue
		}
		if ok, note := ethHotInternalCustodyLikely(res.row.Balance, res.accounted.total); ok {
			warnCount++
			out.Findings = append(out.Findings, Finding{
				Subject:    res.row.Address,
				Field:      res.row.Coin + "@" + res.row.Network + "#" + fmt.Sprintf("%d", matchedHeight),
				Claim:      res.row.Balance.String(),
				Actual:     res.accounted.total.String(),
				Status:     VerdictWarn,
				Note:       note + fmt.Sprintf("; observed delta=%s (provider=%s)", diff.String(), res.used),
				Components: res.accounted.components,
			})
			continue
		}
		_, surplus := classifyBalanceMismatch(res.accounted.total, res.row.Balance)
		note := fmt.Sprintf("accounted on-chain balance != csv claim by %s (provider=%s)", diff.String(), res.used)
		if boundaryNote != "" {
			note = boundaryNote + "; " + note
		}
		status := VerdictFail
		if surplus {
			status = VerdictWarn
			note = "chain observed > CSV row allocation; " + note
			warnCount++
		} else {
			failCount++
		}
		out.Findings = append(out.Findings, Finding{
			Subject:    res.row.Address,
			Field:      res.row.Coin + "@" + res.row.Network + "#" + fmt.Sprintf("%d", matchedHeight),
			Claim:      res.row.Balance.String(),
			Actual:     res.accounted.total.String(),
			Status:     status,
			Note:       note,
			Components: res.accounted.components,
		})
	}

	out.Coverage = float64(passCount+failCount) / float64(len(allRows))
	out.Findings = append(out.Findings, unsupportedFindings(unsupportedPairs)...)
	out.Findings = append(out.Findings, Finding{
		Subject: "summary",
		Field:   "verified_pairs",
		Actual: fmt.Sprintf("%d pass (%d within tol but non-zero diff) / %d warn / %d fail / %d rpc-error  -- of %d supported / %d total HotCold rows",
			passCount, toleranceNoiseCount, warnCount, failCount, errCount, len(supported), len(allRows)),
		Status: VerdictPass,
	})

	switch {
	case failCount > 0:
		out.Verdict = VerdictFail
	case errCount > len(supported)/4:
		out.Verdict = VerdictWarn
		out.Reason = "more than 25% of supported rows failed to query rpc"
	case warnCount > 0:
		out.Verdict = VerdictWarn
		out.Reason = "some rows show chain-observed balance above CSV (surplus) or incomplete staking/deposit attribution — not reserve shortfall"
	case len(supported) < len(allRows):
		out.Verdict = VerdictPartial
	default:
		out.Verdict = VerdictPass
	}
	return out
}

func unsupportedFindings(pairs map[string]int) []Finding {
	keys := make([]string, 0, len(pairs))
	for k := range pairs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]Finding, 0, len(keys))
	for _, k := range keys {
		out = append(out, Finding{
			Subject: k,
			Field:   "row_count",
			Actual:  fmt.Sprintf("%d", pairs[k]),
			Status:  VerdictUnverifiable,
			Note:    "no native verifier for this (coin,network) pair yet",
		})
	}
	return out
}
