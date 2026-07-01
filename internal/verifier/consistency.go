package verifier

import (
	"fmt"
	"sort"
	"time"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/por"
	"github.com/ardmere/ardmere/internal/walletzip"
)

// InternalConsistency reconciles an exchange summary snapshot against the
// aggregated wallet-address CSVs. For each coin in the summary, both
// exchange reserve and third-party reserve must equal the sum of CSV
// balance rows split by custodian, within tolerance.
//
// Tolerance is a *relative* fraction (e.g. 1e-12) to absorb float<->decimal
// representation noise; absolute tolerance also applies as a floor for very
// small balances (< 1).
type InternalConsistency struct {
	SummarySha256   string
	WalletZipSha256 string
	Snapshot        por.Snapshot
	Aggregate       *walletzip.Aggregate
}

// Run produces a Verification.
func (v InternalConsistency) Run() Verification {
	if v.Snapshot.Exchange == "okx" {
		return v.runOKX()
	}
	return v.runDefault()
}

func (v InternalConsistency) runDefault() Verification {
	// BAPI returns float64 with ~9-10 significant figures of precision; the
	// per-row CSV is full decimal. Cumulative truncation noise scales with
	// (#rows × per-row precision) so larger pools have larger acceptable noise.
	// We accept up to 1e-5 relative *or* 1e-2 absolute (tiny-coin floor).
	const (
		tolRel = 1e-5
		tolAbs = 1e-2
	)

	out := Verification{
		VerifierID:     "internal-consistency",
		Version:        "1.1",
		SnapshotID:     v.Snapshot.ID,
		VerifiedAt:     time.Now().UTC(),
		InputArtifacts: []string{v.SummarySha256, v.WalletZipSha256},
		Coverage:       1.0,
	}

	if v.Aggregate != nil && len(v.Aggregate.BySource) > 0 {
		out.Findings = append(out.Findings, walletSourceSummary(v.Aggregate)...)
	}

	allPass := true
	for _, c := range v.Snapshot.CoinSummaries {
		exchActual := v.Aggregate.Exchange[c.Coin]
		tpActual := v.Aggregate.ThirdParty[c.Coin]

		exchClaim := c.ExchangeReserve
		tpClaim := c.ThirdPartyReserve

		exchOk := approxEqual(exchClaim, exchActual, tolRel, tolAbs)
		tpOk := approxEqual(tpClaim, tpActual, tolRel, tolAbs)

		st := VerdictPass
		if !exchOk || !tpOk {
			st = VerdictFail
			allPass = false
		}
		if !exchOk {
			out.Findings = append(out.Findings, Finding{
				Subject: c.Coin,
				Field:   "exchangeBalance",
				Claim:   exchClaim.String(),
				Actual:  exchActual.String(),
				Status:  VerdictFail,
				Note:    "csv aggregate (custodian=='') != summary exchangeBalance",
			})
		}
		if !tpOk {
			out.Findings = append(out.Findings, Finding{
				Subject: c.Coin,
				Field:   "thirdPartyCustody",
				Claim:   tpClaim.String(),
				Actual:  tpActual.String(),
				Status:  VerdictFail,
				Note:    "csv aggregate (custodian!='') != summary thirdPartyCustody",
			})
		}
		if exchOk && tpOk {
			out.Findings = append(out.Findings, Finding{
				Subject: c.Coin,
				Field:   "exchangeLiability",
				Claim:   c.ExchangeLiability.String(),
				Actual:  exchActual.Add(tpActual).String(),
				Status:  st,
			})
		}
	}

	// The summary may list fewer coins than the wallet CSV (long-tail).
	// We only assert reconciliation for summary-listed coins; extras are surfaced
	// as informational WARN findings.
	summaryCoins := map[string]bool{}
	for _, c := range v.Snapshot.CoinSummaries {
		summaryCoins[c.Coin] = true
	}
	var extras []string
	for c := range v.Aggregate.Exchange {
		if !summaryCoins[c] {
			extras = append(extras, c)
		}
	}
	for c := range v.Aggregate.ThirdParty {
		if !summaryCoins[c] && !contains(extras, c) {
			extras = append(extras, c)
		}
	}
	sort.Strings(extras)
	if len(extras) > 0 {
		out.Findings = append(out.Findings, Finding{
			Subject: "csv-extras",
			Field:   "longTailCoins",
			Actual:  fmt.Sprintf("%d coins not in summary snapshot: %v", len(extras), extras),
			Status:  VerdictWarn,
			Note:    fmt.Sprintf("informational — %s summary lists top coins only", v.Snapshot.Exchange),
		})
	}

	if allPass {
		out.Verdict = VerdictPass
	} else {
		out.Verdict = VerdictFail
	}
	return out
}

func approxEqual(a, b decimal.Decimal, relTol, absTol float64) bool {
	diff := a.Sub(b).Abs()
	if diff.LessThanOrEqual(decimal.NewFromFloat(absTol)) {
		return true
	}
	maxAbs := decimal.Max(a.Abs(), b.Abs())
	if maxAbs.IsZero() {
		return diff.IsZero()
	}
	rel := diff.Div(maxAbs)
	return rel.LessThanOrEqual(decimal.NewFromFloat(relTol))
}

func contains(ss []string, s string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}
	return false
}
