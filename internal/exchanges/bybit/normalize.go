package bybit

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/exchanges/bybit/bybitapi"
	"github.com/ardmere/ardmere/internal/por"
)

const ID = "bybit"

// Meta enriches parsed Bybit summary with bundle metadata.
type Meta struct {
	Exchange  string
	PeriodSeq uint32
}

// Normalize maps Bybit public summary into por.Snapshot.
func Normalize(raw bybitapi.ParsedSummary, meta Meta) por.Snapshot {
	if meta.Exchange == "" {
		meta.Exchange = ID
	}
	snapTime := raw.AuditTime
	if snapTime.IsZero() && raw.AuditTimeRaw != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", raw.AuditTimeRaw); err == nil {
			snapTime = t.UTC()
		}
	}

	out := por.Snapshot{
		Exchange:        meta.Exchange,
		ID:              raw.AuditID,
		SnapshotTime:    snapTime,
		SnapshotTimeRaw: raw.AuditTimeRaw,
		MerkleRoot:      raw.MerkleRoot,
		PeriodSeq:       meta.PeriodSeq,
		Extra: map[string]string{
			"totalReserveRate": fmtDecimal(raw.TotalReserveRate),
			"algorithm":        "merkle_tree+sha256",
			"auditor":          "Hacken",
		},
	}

	for _, c := range raw.CoinRows {
		out.CoinSummaries = append(out.CoinSummaries, por.CoinSummary{
			Coin:              c.Coin,
			CustomerLiability: decimal.NewFromFloat(c.LiabilityAmount),
			ExchangeReserve:   decimal.NewFromFloat(c.ReserveAmount),
			Extra: map[string]string{
				"reserveRate": fmtDecimal(c.ReserveRate),
			},
		})
	}
	return out
}

func fmtDecimal(v float64) string {
	return decimal.NewFromFloat(v).String()
}
