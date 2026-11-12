package bybit_test

import (
	"testing"

	"github.com/ardmere/ardmere/internal/exchanges/bybit"
	"github.com/ardmere/ardmere/internal/exchanges/bybit/bybitapi"
)

func TestNormalize_missingReservesUnpublished(t *testing.T) {
	liab := 100.0
	zero := 0.0
	snap := bybit.Normalize(bybitapi.ParsedSummary{
		AuditID: "2026102706",
		CoinRows: []bybitapi.CoinRow{
			{Coin: "BTC", LiabilityAmount: &liab},
			{Coin: "ETH", LiabilityAmount: &liab, ReserveAmount: &zero, ReserveRate: &zero},
		},
	}, bybit.Meta{})
	if len(snap.CoinSummaries) != 2 {
		t.Fatalf("coins: %+v", snap.CoinSummaries)
	}
	btc := snap.CoinSummaries[0]
	if btc.Extra["reserveAmount"] != "unpublished" || btc.Extra["reserveRate"] != "unpublished" {
		t.Fatalf("btc extra: %+v", btc.Extra)
	}
	if !btc.ExchangeReserve.IsZero() {
		t.Fatalf("btc reserve decimal should be zero placeholder: %s", btc.ExchangeReserve)
	}
	eth := snap.CoinSummaries[1]
	if eth.Extra["reserveAmount"] == "unpublished" {
		t.Fatalf("explicit zero must not be marked unpublished: %+v", eth.Extra)
	}
	if !eth.ExchangeReserve.IsZero() {
		t.Fatalf("eth published zero: %s", eth.ExchangeReserve)
	}
	if eth.Extra["reserveRate"] != "0" {
		t.Fatalf("eth rate: %+v", eth.Extra)
	}
	if _, ok := snap.Extra["totalReserveRate"]; ok {
		t.Fatalf("missing total rate must be omitted: %+v", snap.Extra)
	}
}
