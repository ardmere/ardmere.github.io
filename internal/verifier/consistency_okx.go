package verifier

import (
	"time"

	"github.com/shopspring/decimal"
)

func (v InternalConsistency) runOKX() Verification {
	const (
		tolRel = 1e-5
		tolAbs = 1e-2
	)

	out := Verification{
		VerifierID:     "internal-consistency",
		Version:        "1.0",
		SnapshotID:     v.Snapshot.ID,
		VerifiedAt:     time.Now().UTC(),
		InputArtifacts: []string{v.SummarySha256, v.WalletZipSha256},
		Coverage:       1.0,
	}

	allPass := true
	for _, c := range v.Snapshot.CoinSummaries {
		exchActual := v.Aggregate.Exchange[c.Coin]
		exchClaim := c.ExchangeReserve
		exchOk := approxEqual(exchClaim, exchActual, tolRel, tolAbs)
		st := VerdictPass
		if !exchOk {
			st = VerdictFail
			allPass = false
			out.Findings = append(out.Findings, Finding{
				Subject: c.Coin,
				Field:   "exchangeReserveBalances",
				Claim:   exchClaim.String(),
				Actual:  exchActual.String(),
				Status:  VerdictFail,
				Note:    "signed address CSV aggregate (+ ETH staking) != summary exchangeReserveBalances",
			})
		} else {
			out.Findings = append(out.Findings, Finding{
				Subject: c.Coin,
				Field:   "exchangeReserveBalances",
				Claim:   exchClaim.String(),
				Actual:  exchActual.String(),
				Status:  st,
			})
		}

		if c.ThirdPartyReserve.GreaterThan(decimal.Zero) {
			out.Findings = append(out.Findings, Finding{
				Subject: c.Coin,
				Field:   "custodyReserveBalances",
				Claim:   c.ThirdPartyReserve.String(),
				Status:  VerdictWarn,
				Note:    "custody balances are summary-only; public wallet CSV does not break out third-party addresses",
			})
		}
	}

	if allPass {
		out.Verdict = VerdictPass
	} else {
		out.Verdict = VerdictFail
	}
	return out
}
