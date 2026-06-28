package gateio

import (
	"testing"

	"github.com/ardmere/ardmere/internal/exchanges/gateio/gateapi"
)

func TestNormalize(t *testing.T) {
	raw := gateapi.ParsedSummary{
		AuditID:          "20231110",
		AuditTimeRaw:     "2023-11-10 00:00:00",
		MerkleRoot:       "0xabc",
		TotalReserveRate: 115.69,
		CoinRows: []gateapi.CoinRow{{
			Coin: "BTC", ReserveRate: 123.52, ReserveAmount: 17100, LiabilityAmount: 13800,
		}},
	}
	got := Normalize(raw, Meta{Exchange: "gateio", PeriodSeq: 1})
	if got.Exchange != "gateio" || got.MerkleRoot != "0xabc" {
		t.Fatalf("identity: %+v", got)
	}
	if got.Extra["totalReserveRate"] != "115.69" {
		t.Fatalf("extra: %+v", got.Extra)
	}
	if len(got.CoinSummaries) != 1 || got.CoinSummaries[0].Coin != "BTC" {
		t.Fatalf("coins: %+v", got.CoinSummaries)
	}
}
