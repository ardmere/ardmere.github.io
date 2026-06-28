package htx

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/exchanges/htx/htxapi"
	"github.com/ardmere/ardmere/internal/por"
)

const ID = "htx"

// Meta enriches parsed HTX summary with bundle metadata.
type Meta struct {
	Exchange  string
	PeriodSeq uint32
	ZkDerived bool
}

// Normalize maps HTX summary into por.Snapshot.
func Normalize(raw htxapi.ParsedSummary, meta Meta) por.Snapshot {
	if meta.Exchange == "" {
		meta.Exchange = ID
	}
	snapTime := raw.AuditTime
	if snapTime.IsZero() && raw.AuditTimeRaw != "" {
		if t, err := time.Parse(time.RFC3339, raw.AuditTimeRaw); err == nil {
			snapTime = t.UTC()
		}
	}

	extra := map[string]string{
		"algorithm": "groth16+zkpor500",
	}
	if meta.ZkDerived {
		extra["summarySource"] = "zk-bundle-derived"
	}
	if raw.TotalReserveRate > 0 {
		extra["totalReserveRate"] = fmtDecimal(raw.TotalReserveRate)
	}

	out := por.Snapshot{
		Exchange:        meta.Exchange,
		ID:              raw.AuditID,
		SnapshotTime:    snapTime,
		SnapshotTimeRaw: raw.AuditTimeRaw,
		MerkleRoot:      raw.MerkleRoot,
		PeriodSeq:       meta.PeriodSeq,
		Extra:           extra,
		Auditor:         "HTX zkPoR (Groth16)",
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
