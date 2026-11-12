package poloniex

import (
	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/exchanges/poloniex/poloniexapi"
	"github.com/ardmere/ardmere/internal/por"
)

const ID = "poloniex"

type Meta struct {
	Exchange  string
	PeriodSeq uint32
}

// Normalize maps Poloniex public summary into por.Snapshot.
// Missing reserve/liability amounts are unpublished, not decimal 0.
func Normalize(raw poloniexapi.ParsedSummary, meta Meta) por.Snapshot {
	if meta.Exchange == "" {
		meta.Exchange = ID
	}
	out := por.Snapshot{
		Exchange:        meta.Exchange,
		ID:              raw.AuditID,
		SnapshotTime:    raw.AuditTime,
		SnapshotTimeRaw: raw.AuditTimeRaw,
		MerkleRoot:      raw.MerkleRoot,
		PeriodSeq:       meta.PeriodSeq,
		Extra: map[string]string{
			"totalReserveRate": decimal.NewFromFloat(raw.TotalReserveRate).String(),
			"algorithm":        "merkle_tree+undocumented",
			"auditor":          "self-attested (no named third-party in public summary)",
		},
	}
	anyAmounts := false
	for _, c := range raw.CoinRows {
		extra := map[string]string{
			por.ExtraReserveRate: decimal.NewFromFloat(c.ReserveRate).String(),
		}
		if c.CurrencyGroup != "" {
			extra["currencyGroup"] = c.CurrencyGroup
		}
		cs := por.CoinSummary{Coin: c.Coin}
		if c.ReserveAmount == 0 && c.LiabilityAmount == 0 {
			extra = por.MarkAmountsUnpublished(extra)
		} else {
			anyAmounts = true
			cs.ExchangeReserve = decimal.NewFromFloat(c.ReserveAmount)
			cs.CustomerLiability = decimal.NewFromFloat(c.LiabilityAmount)
		}
		cs.Extra = extra
		out.CoinSummaries = append(out.CoinSummaries, cs)
	}
	if !anyAmounts {
		out.Extra = por.SnapshotAmountsNotPublished(out.Extra)
	}
	return out
}
