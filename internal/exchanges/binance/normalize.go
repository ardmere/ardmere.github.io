package binance

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/exchanges/binance/bapi"
	"github.com/ardmere/ardmere/internal/por"
)

// Meta enriches a raw BAPI snapshot with bundle / list metadata.
type Meta struct {
	Exchange       string
	PeriodSeq        uint32
	BTCBlockHeight   uint32
	SnapshotTime     time.Time
}

// Normalize maps Binance BAPI JSON into the exchange-neutral por.Snapshot.
func Normalize(raw bapi.Snapshot, meta Meta) por.Snapshot {
	if meta.Exchange == "" {
		meta.Exchange = ID
	}
	snapTime := meta.SnapshotTime
	if snapTime.IsZero() {
		if t, err := bapi.ParseSnapshotTime(raw.SnapshotTime); err == nil {
			snapTime = t
		}
	}

	out := por.Snapshot{
		Exchange:        meta.Exchange,
		ID:              raw.AuditID,
		SnapshotTime:    snapTime,
		SnapshotTimeRaw: raw.SnapshotTime,
		MerkleRoot:      raw.MerkleRootHash,
		PeriodSeq:       meta.PeriodSeq,
		Auditor:         raw.Auditor,
		AuditorLink:     raw.AuditorLink,
		AuditDate:       raw.AuditDate,
	}
	if meta.BTCBlockHeight > 0 {
		out.TimeAnchor = &por.TimeAnchor{
			Kind:   "btc_block",
			Height: meta.BTCBlockHeight,
		}
	}

	for _, c := range raw.SnapshotDataList {
		out.CoinSummaries = append(out.CoinSummaries, por.CoinSummary{
			Coin:              c.Coin,
			CustomerLiability: decimal.NewFromFloat(c.CustomerLiability),
			ExchangeReserve:   decimal.NewFromFloat(c.ExchangeBalance),
			ThirdPartyReserve: decimal.NewFromFloat(c.ThirdPartyCustody),
			ExchangeLiability: decimal.NewFromFloat(c.BinanceLiability),
		})
	}
	return out
}
