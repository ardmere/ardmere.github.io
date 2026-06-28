package binance

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/exchanges/binance/bapi"
)

func TestNormalize(t *testing.T) {
	raw := bapi.Snapshot{
		AuditID:        "PR01JUN26",
		SnapshotTime:   "01/06/26 00:00:00 UTC",
		MerkleRootHash: "0xabc",
		SnapshotDataList: []bapi.CoinRow{{
			Coin:              "BTC",
			CustomerLiability: 100,
			BinanceLiability:  110,
			ExchangeBalance:   80,
			ThirdPartyCustody: 30,
		}},
	}
	snapTime := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	got := Normalize(raw, Meta{
		Exchange:       "binance",
		PeriodSeq:      43,
		BTCBlockHeight: 900000,
		SnapshotTime:   snapTime,
	})

	if got.Exchange != "binance" || got.ID != "PR01JUN26" {
		t.Fatalf("identity: %+v", got)
	}
	if got.PeriodSeq != 43 || got.TimeAnchor.Height != 900000 {
		t.Fatalf("meta: %+v", got)
	}
	if len(got.CoinSummaries) != 1 {
		t.Fatalf("coins: %d", len(got.CoinSummaries))
	}
	c := got.CoinSummaries[0]
	if !c.ExchangeReserve.Equal(decimal.NewFromInt(80)) {
		t.Fatalf("exchange reserve: %s", c.ExchangeReserve)
	}
	if !c.ThirdPartyReserve.Equal(decimal.NewFromInt(30)) {
		t.Fatalf("third party: %s", c.ThirdPartyReserve)
	}
}
