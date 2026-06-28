package verifier

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/rpc"
	"github.com/ardmere/ardmere/internal/walletzip"
)

// OnchainBalanceDeposit value-samples Deposit.csv rows and verifies on-chain
// balances using the same routes as hot/token/ledger verifiers.
type OnchainBalanceDeposit struct {
	Exchange        string
	WalletZipPath   string
	WalletZipSha256 string
	BapiSha256      string
	SnapshotID      string
	RPC             *rpc.Client
	Ledger          *rpc.LedgerClient
	ETHDeposits     ETHDepositIndexer
	Concurrency     int
	MaxSamples      int
	TopKPerCoin     int
	ValueCoverage   float64
}

func (v OnchainBalanceDeposit) ledgerClient() *rpc.LedgerClient {
	if v.Ledger != nil {
		return v.Ledger
	}
	return rpc.NewLedger()
}

func (v OnchainBalanceDeposit) sampleOpts() walletzip.DepositSampleOpts {
	topK := v.TopKPerCoin
	if topK <= 0 {
		topK = envInt("DEPOSIT_TOP_K_PER_COIN", 5000)
	}
	maxTotal := v.MaxSamples
	if maxTotal <= 0 {
		maxTotal = envInt("DEPOSIT_MAX_SAMPLES", 5000)
	}
	target := v.ValueCoverage
	if target <= 0 {
		target = envFloat("DEPOSIT_VALUE_COVERAGE", 0.99)
	}
	return walletzip.DepositSampleOpts{
		TopKPerCoin:         topK,
		MaxTotal:            maxTotal,
		ValueCoverageTarget: target,
	}
}

func (v OnchainBalanceDeposit) Run(ctx context.Context) Verification {
	if v.Concurrency <= 0 {
		v.Concurrency = 4
	}

	out := Verification{
		VerifierID:     "onchain-balance-deposit",
		Version:        "1.2",
		SnapshotID:     v.SnapshotID,
		VerifiedAt:     time.Now().UTC(),
		InputArtifacts: []string{v.BapiSha256, v.WalletZipSha256},
	}

	sample, err := walletzip.SampleDepositRows(v.WalletZipPath, func(r walletzip.Row) bool {
		_, ok := routeOnchainRow(v.Exchange, r.Coin, r.Network)
		return ok
	}, v.sampleOpts())
	if err != nil {
		out.Verdict = VerdictFail
		out.Reason = fmt.Sprintf("sample Deposit.csv: %v", err)
		return out
	}

	out.Findings = append(out.Findings, Finding{
		Subject: "sample",
		Field:   "value_coverage",
		Actual:  walletzip.FormatDepositSampleSummary(sample),
		Status:  VerdictPass,
		Note:    "value-weighted head sampling of verifiable exchange-owned deposit rows",
	})

	if len(sample.Rows) == 0 {
		out.Verdict = VerdictPartial
		out.Reason = "no verifiable Deposit rows in supported (coin,network) set"
		return out
	}

	out.Coverage = sample.ValueCoverage

	tokenV := OnchainBalanceToken{Exchange: v.Exchange, RPC: v.RPC, Ledger: v.ledgerClient(), Concurrency: v.Concurrency}
	hotV := OnchainBalanceHot{Exchange: v.Exchange, RPC: v.RPC, ETHDeposits: v.ETHDeposits, Concurrency: v.Concurrency}
	ledgerV := OnchainBalanceLedger{Ledger: v.ledgerClient(), Concurrency: v.Concurrency}
	stakeCache := newBSCStakeHubCache()
	ethCache := newETHDepositCache()
	decCache := newDecimalsCache()
	sonicCache := newSonicSFCCache()

	type result struct {
		row        walletzip.Row
		actual     decimal.Decimal
		used       string
		route      string
		mode       string
		components map[string]string
		err        error
		unv        bool
		unvNote    string
	}
	jobs := make(chan walletzip.Row, len(sample.Rows))
	results := make(chan result, len(sample.Rows))
	var wg sync.WaitGroup
	for i := 0; i < v.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for r := range jobs {
				route, _ := routeOnchainRow(v.Exchange, r.Coin, r.Network)
				var actual decimal.Decimal
				var used, unvNote string
				var err error
				var unv bool
				switch route {
				case routeNative:
					acc, e := hotV.accountedBalance(ctx, r, stakeCache, ethCache)
					actual, used, err = acc.total, acc.used, e
				case routeToken:
					actual, used, _, err = tokenV.tokenBalanceAtHeight(ctx, r, decCache, sonicCache)
				case routeLedger:
					spec := loadLedgerSupportedFor(v.Exchange)[r.Coin+"|"+r.Network]
					var mode string
					var comps map[string]string
					actual, used, mode, comps, unv, unvNote, err = ledgerV.balanceAtHeight(ctx, r, spec)
					_ = mode
					results <- result{row: r, actual: actual, used: used, route: string(routeLedger), components: comps, err: err, unv: unv, unvNote: unvNote}
					continue
				}
				results <- result{row: r, actual: actual, used: used, route: string(route), err: err, unv: unv, unvNote: unvNote}
			}
		}()
	}
	for _, r := range sample.Rows {
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
				Field:   res.row.Coin + "|" + res.row.Network,
				Status:  VerdictWarn,
				Note:    fmt.Sprintf("deposit rpc error (provider=%s): %v", res.used, res.err),
			})
			continue
		}
		if res.unv {
			unverifiableCount++
			out.Findings = append(out.Findings, Finding{
				Subject: res.row.Address,
				Field:   res.row.Coin + "|" + res.row.Network,
				Claim:   res.row.Balance.String(),
				Actual:  res.actual.String(),
				Status:  VerdictUnverifiable,
				Note:    res.unvNote,
			})
			continue
		}
		within, surplus := classifyBalanceMismatch(res.actual, res.row.Balance)
		if within {
			passCount++
			continue
		}
		diff := res.actual.Sub(res.row.Balance)
		note := fmt.Sprintf("deposit on-chain != csv by %s (provider=%s, route=%s)", diff.Abs().String(), res.used, res.route)
		status := VerdictFail
		if surplus {
			status = VerdictWarn
		} else if res.route == string(routeLedger) {
			if spec, ok := loadLedgerSupportedFor(v.Exchange)[res.row.Coin+"|"+res.row.Network]; ok {
				if snap := ledgerLiveSnapshotNote(spec, res.components); snap != "" {
					status = VerdictWarn
					note = snap + "; " + note
				}
			}
		}
		if status == VerdictWarn {
			warnCount++
		} else {
			failCount++
		}
		out.Findings = append(out.Findings, Finding{
			Subject:    res.row.Address,
			Field:      res.row.Coin + "|" + res.row.Network + "#" + fmt.Sprintf("%d", res.row.Height),
			Claim:      res.row.Balance.String(),
			Actual:     res.actual.String(),
			Status:     status,
			Note:       note,
			Components: res.components,
		})
	}

	out.Findings = append(out.Findings, Finding{
		Subject: "summary",
		Field:   "sampled_deposit_rows",
		Actual: fmt.Sprintf("%d pass / %d warn / %d fail / %d unverifiable / %d rpc-error — sampled %d rows (%.2f%% deposit value coverage, %d verifiable total)",
			passCount, warnCount, failCount, unverifiableCount, errCount,
			len(sample.Rows), sample.ValueCoverage*100, sample.VerifiableRows),
		Status: VerdictPass,
	})

	switch {
	case failCount > 0:
		out.Verdict = VerdictFail
	case errCount > len(sample.Rows)/4:
		out.Verdict = VerdictWarn
		out.Reason = "more than 25% of sampled deposit rows failed RPC"
	case warnCount > 0 || unverifiableCount > 0:
		out.Verdict = VerdictWarn
		out.Reason = "some sampled deposit rows warn or are unverifiable"
	default:
		out.Verdict = VerdictPass
	}
	return out
}

type onchainRoute string

const (
	routeNative onchainRoute = "native"
	routeToken  onchainRoute = "token"
	routeLedger onchainRoute = "ledger"
)

func routeOnchainRow(exchangeID, coin, network string) (onchainRoute, bool) {
	key := coin + "|" + network
	if _, ok := loadNativeSupported(exchangeID)[key]; ok {
		return routeNative, true
	}
	if _, ok := loadTokenSupportedFor(exchangeID)[key]; ok {
		return routeToken, true
	}
	if _, ok := loadLedgerSupportedFor(exchangeID)[key]; ok {
		return routeLedger, true
	}
	return "", false
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return def
}

func envFloat(key string, def float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f > 0 && f <= 1 {
			return f
		}
	}
	return def
}
