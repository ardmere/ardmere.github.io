package gateio

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/exchanges/gateio/gateapi"
	"github.com/ardmere/ardmere/internal/por"
)

const ID = "gateio"

// Meta enriches parsed Gate summary with bundle metadata.
type Meta struct {
	Exchange  string
	PeriodSeq uint32
}

// Normalize maps Gate public summary into por.Snapshot.
func Normalize(raw gateapi.ParsedSummary, meta Meta) por.Snapshot {
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
			"totalReserveUSD":  fmtDecimal(raw.TotalReserveUSD),
			"customerNetUSD":   fmtDecimal(raw.CustomerNetUSD),
			"excessReserveUSD": fmtDecimal(raw.ExcessReserveUSD),
			"algorithm":        "merkle_tree+zk_snark",
		},
	}

	anyAmounts := false
	for _, c := range raw.CoinRows {
		extra := map[string]string{
			por.ExtraReserveRate: fmtDecimal(c.ReserveRate),
		}
		cs := por.CoinSummary{Coin: c.Coin}
		if c.ReserveAmount == 0 && c.LiabilityAmount == 0 {
			extra = por.MarkAmountsUnpublished(extra)
		} else {
			anyAmounts = true
			cs.CustomerLiability = decimal.NewFromFloat(c.LiabilityAmount)
			cs.ExchangeReserve = decimal.NewFromFloat(c.ReserveAmount)
		}
		cs.Extra = extra
		out.CoinSummaries = append(out.CoinSummaries, cs)
	}
	if !anyAmounts {
		out.Extra = por.SnapshotAmountsNotPublished(out.Extra)
	}
	return out
}

func fmtDecimal(v float64) string {
	return decimal.NewFromFloat(v).String()
}
