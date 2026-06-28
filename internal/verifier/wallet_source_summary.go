package verifier

import (
	"fmt"
	"sort"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/walletzip"
)

func walletSourceSummary(agg *walletzip.Aggregate) []Finding {
	var out []Finding
	hc := agg.BySource[walletzip.HotCold]
	dep := agg.BySource[walletzip.Deposit]
	var hcRows, depRows int64
	if hc != nil {
		hcRows = hc.TotalRows
	}
	if dep != nil {
		depRows = dep.TotalRows
	}
	out = append(out, Finding{
		Subject: "wallet-sources",
		Field:   "row_counts",
		Actual:  fmt.Sprintf("HotCold=%d Deposit=%d total=%d", hcRows, depRows, agg.TotalRows),
		Status:  VerdictPass,
		Note:    "HotCold+Deposit aggregated for BAPI reconciliation",
	})

	for _, coin := range summaryCoinOrder(agg) {
		hcEx := decimal.Zero
		depEx := decimal.Zero
		if hc != nil {
			hcEx = hc.Exchange[coin]
		}
		if dep != nil {
			depEx = dep.Exchange[coin]
		}
		if hcEx.IsZero() && depEx.IsZero() {
			continue
		}
		total := hcEx.Add(depEx)
		if total.IsZero() {
			continue
		}
		hcPct := decimal.Zero
		if !total.IsZero() {
			hcPct = hcEx.Div(total).Mul(decimal.NewFromInt(100))
		}
		out = append(out, Finding{
			Subject: coin,
			Field:   "exchangeBalance_split",
			Claim:   total.String(),
			Actual:  fmt.Sprintf("HotCold=%s Deposit=%s (HotCold %.1f%%)", hcEx, depEx, hcPct.InexactFloat64()),
			Status:  VerdictPass,
			Note:    "exchange-owned rows only (custodian==\"\")",
		})
	}
	return out
}

func summaryCoinOrder(agg *walletzip.Aggregate) []string {
	set := map[string]bool{}
	for c := range agg.Exchange {
		set[c] = true
	}
	var coins []string
	for c := range set {
		coins = append(coins, c)
	}
	sort.Slice(coins, func(i, j int) bool {
		ti := agg.Exchange[coins[i]]
		tj := agg.Exchange[coins[j]]
		return ti.GreaterThan(tj)
	})
	if len(coins) > 14 {
		coins = coins[:14]
	}
	return coins
}
