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
	for _, c := range raw.CoinRows {
		extra := map[string]string{
			"reserveRate": decimal.NewFromFloat(c.ReserveRate).String(),
		}
		if c.CurrencyGroup != "" {
			extra["currencyGroup"] = c.CurrencyGroup
		}
		out.CoinSummaries = append(out.CoinSummaries, por.CoinSummary{
			Coin:              c.Coin,
			ExchangeReserve:   decimal.NewFromFloat(c.ReserveAmount),
			CustomerLiability: decimal.NewFromFloat(c.LiabilityAmount),
			Extra:             extra,
		})
	}
	return out
}
