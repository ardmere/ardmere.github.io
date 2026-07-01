package verifier

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ardmere/ardmere/internal/por"
)

const solvencySelfReportedNote = "self-reported exchange summary only; does not prove reserve authenticity"

// SolvencyClaim checks exchange-declared solvency ratios from the summary artifact.
// Gate.io: total_reserve_rate and per-coin reserve_rate >= 100%.
// Binance: binanceLiability >= customerLiability per coin.
type SolvencyClaim struct {
	SummarySha256 string
	Snapshot      por.Snapshot
}

func (v SolvencyClaim) Run() Verification {
	out := Verification{
		VerifierID:     "solvency-claim",
		Version:        "1.0",
		SnapshotID:     v.Snapshot.ID,
		VerifiedAt:     time.Now().UTC(),
		InputArtifacts: []string{v.SummarySha256},
		Coverage:       1.0,
	}

	switch v.Snapshot.Exchange {
	case "gateio":
		return v.runGate(out)
	case "bybit":
		return v.runGate(out)
	case "bitget":
		return v.runGate(out)
	case "htx":
		return v.runHTX(out)
	case "okx":
		return v.runOKX(out)
	default:
		return v.runBinance(out)
	}
}

func (v SolvencyClaim) runGate(out Verification) Verification {
	const minRate = 100.0
	allPass := true
	checked := 0

	if rateStr, ok := v.Snapshot.Extra["totalReserveRate"]; ok && rateStr != "" {
		rate, err := strconv.ParseFloat(rateStr, 64)
		if err != nil {
			out.Findings = append(out.Findings, Finding{
				Subject: "total",
				Field:   "totalReserveRate",
				Claim:   rateStr,
				Status:  VerdictFail,
				Note:    fmt.Sprintf("parse totalReserveRate: %v", err),
			})
			allPass = false
		} else {
			checked++
			st := VerdictPass
			if rate < minRate {
				st = VerdictFail
				allPass = false
			}
			out.Findings = append(out.Findings, Finding{
				Subject: "total",
				Field:   "totalReserveRate",
				Claim:   rateStr,
				Actual:  fmt.Sprintf("%.2f", rate),
				Status:  st,
				Note:    solvencySelfReportedNote,
			})
		}
	} else {
		out.Verdict = VerdictFail
		out.Reason = "gate summary missing totalReserveRate"
		return out
	}

	for _, c := range v.Snapshot.CoinSummaries {
		rateStr := c.Extra["reserveRate"]
		if rateStr == "" {
			continue
		}
		rate, err := strconv.ParseFloat(rateStr, 64)
		if err != nil {
			out.Findings = append(out.Findings, Finding{
				Subject: c.Coin,
				Field:   "reserveRate",
				Claim:   rateStr,
				Status:  VerdictFail,
				Note:    fmt.Sprintf("parse reserveRate: %v", err),
			})
			allPass = false
			continue
		}
		checked++
		st := VerdictPass
		if rate < minRate {
			st = VerdictFail
			allPass = false
		}
		out.Findings = append(out.Findings, Finding{
			Subject: c.Coin,
			Field:   "reserveRate",
			Claim:   rateStr,
			Actual:  fmt.Sprintf("%.2f", rate),
			Status:  st,
			Note:    solvencySelfReportedNote,
		})
	}

	if checked == 0 {
		out.Verdict = VerdictFail
		out.Reason = "no solvency fields to check"
		return out
	}
	if len(v.Snapshot.CoinSummaries) > 0 {
		out.Coverage = float64(checked) / float64(1+len(v.Snapshot.CoinSummaries))
	}
	if allPass {
		out.Verdict = VerdictPass
	} else {
		out.Verdict = VerdictFail
	}
	return out
}

func (v SolvencyClaim) runHTX(out Verification) Verification {
	if v.Snapshot.Extra["summarySource"] == "zk-bundle-derived" {
		out.Verdict = VerdictUnverifiable
		out.Reason = "zk-derived summary has no public reserve ratios; import browser-captured ratio JSON via -summary-path"
		out.Coverage = 0
		return out
	}
	return v.runGate(out)
}

func (v SolvencyClaim) runBinance(out Verification) Verification {
	if len(v.Snapshot.CoinSummaries) == 0 {
		out.Verdict = VerdictFail
		out.Reason = "summary has no coin rows"
		return out
	}

	allPass := true
	for _, c := range v.Snapshot.CoinSummaries {
		claim := c.ExchangeLiability
		actual := c.CustomerLiability
		ok := claim.GreaterThan(actual) || claim.Equal(actual)
		st := VerdictPass
		if !ok {
			st = VerdictFail
			allPass = false
		}
		out.Findings = append(out.Findings, Finding{
			Subject: c.Coin,
			Field:   "binanceLiability>=customerLiability",
			Claim:   claim.String(),
			Actual:  actual.String(),
			Status:  st,
			Note:    solvencySelfReportedNote,
		})
	}
	if allPass {
		out.Verdict = VerdictPass
	} else {
		out.Verdict = VerdictFail
	}
	return out
}
