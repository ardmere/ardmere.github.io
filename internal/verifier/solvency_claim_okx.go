package verifier

import (
	"fmt"
	"strconv"
	"strings"
)

func (v SolvencyClaim) runOKX(out Verification) Verification {
	const minRate = 100.0
	if len(v.Snapshot.CoinSummaries) == 0 {
		out.Verdict = VerdictFail
		out.Reason = "okx summary has no reserve currencies"
		return out
	}

	allPass := true
	checked := 0
	for _, c := range v.Snapshot.CoinSummaries {
		ratioStr := c.Extra["capitalRatio"]
		if ratioStr == "" {
			continue
		}
		ratioStr = strings.TrimSuffix(strings.TrimSpace(ratioStr), "%")
		rate, err := strconv.ParseFloat(ratioStr, 64)
		if err != nil {
			out.Findings = append(out.Findings, Finding{
				Subject: c.Coin,
				Field:   "capitalRatio",
				Claim:   c.Extra["capitalRatio"],
				Status:  VerdictFail,
				Note:    fmt.Sprintf("parse capitalRatio: %v", err),
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
			Field:   "capitalRatio>=100%",
			Claim:   c.Extra["capitalRatio"],
			Actual:  fmt.Sprintf("%.0f%%", rate),
			Status:  st,
			Note:    solvencySelfReportedNote,
		})
	}
	if checked == 0 {
		out.Verdict = VerdictFail
		out.Reason = "no capitalRatio fields to check"
		return out
	}
	out.Coverage = float64(checked) / float64(len(v.Snapshot.CoinSummaries))
	if allPass {
		out.Verdict = VerdictPass
	} else {
		out.Verdict = VerdictFail
	}
	return out
}
