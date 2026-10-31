package kucoin

import (
	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/exchanges/kucoin/kucoinapi"
	"github.com/ardmere/ardmere/internal/por"
)

const ID = "kucoin"

type Meta struct {
	Exchange  string
	PeriodSeq uint32
}

func Normalize(raw kucoinapi.ParsedSummary, meta Meta) por.Snapshot {
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
			"auditor":          "third-party (page unnamed for latest; Hacken historically)",
		},
	}
	for _, c := range raw.CoinRows {
		out.CoinSummaries = append(out.CoinSummaries, por.CoinSummary{
			Coin:              c.Coin,
			ExchangeReserve:   decimal.NewFromFloat(c.ReserveAmount),
			CustomerLiability: decimal.NewFromFloat(c.LiabilityAmount),
			Extra: map[string]string{
				"reserveRate": decimal.NewFromFloat(c.ReserveRate).String(),
			},
		})
	}
	return out
}
