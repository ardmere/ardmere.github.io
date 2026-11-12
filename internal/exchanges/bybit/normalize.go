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
// Missing reserve or liability amounts stay unpublished (Extra marker + zero
// decimal) so pipeline totals can distinguish "not published" from a true zero.
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

	extra := map[string]string{
		"algorithm": "merkle_tree+sha256",
		"auditor":   "Hacken",
	}
	if raw.TotalReserveRate != nil {
		extra["totalReserveRate"] = fmtDecimal(*raw.TotalReserveRate)
	}

	out := por.Snapshot{
		Exchange:        meta.Exchange,
		ID:              raw.AuditID,
		SnapshotTime:    snapTime,
		SnapshotTimeRaw: raw.AuditTimeRaw,
		MerkleRoot:      raw.MerkleRoot,
		PeriodSeq:       meta.PeriodSeq,
		Extra:           extra,
	}

	anyAmounts := false
	for _, c := range raw.CoinRows {
		coinExtra := map[string]string{}
		cs := por.CoinSummary{Coin: c.Coin}
		if c.LiabilityAmount != nil {
			anyAmounts = true
			cs.CustomerLiability = decimal.NewFromFloat(*c.LiabilityAmount)
		} else {
			coinExtra[por.ExtraLiabilityAmount] = por.ExtraUnpublished
		}
		if c.ReserveAmount != nil {
			anyAmounts = true
			cs.ExchangeReserve = decimal.NewFromFloat(*c.ReserveAmount)
		} else {
			coinExtra[por.ExtraReserveAmount] = por.ExtraUnpublished
		}
		if c.ReserveRate != nil {
			coinExtra[por.ExtraReserveRate] = fmtDecimal(*c.ReserveRate)
		} else {
			coinExtra[por.ExtraReserveRate] = por.ExtraUnpublished
		}
		cs.Extra = coinExtra
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
