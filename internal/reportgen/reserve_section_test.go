package reportgen

import (
	"strings"
	"testing"

	"github.com/ardmere/ardmere/internal/pipeline"
)

func TestWriteReserveSection_unavailableNotZero(t *testing.T) {
	var b strings.Builder
	writeReserveSection(&b, pipeline.ReserveFigures{
		TotalLiabilityUsd:     floatPtr(100),
		TotalReserveUsd:       nil,
		UnpricedCoins:         []string{},
		MajorCoinReserveRates: []pipeline.MajorCoinReserveRate{},
		AmountsPublished:      true,
	})
	got := b.String()
	if strings.Contains(got, "`0.00`") || strings.Contains(got, "| `0` |") {
		t.Fatalf("must not render invented zeros:\n%s", got)
	}
	if !strings.Contains(got, "| Total reserves (USD) | Unavailable |") {
		t.Fatalf("reserves cell:\n%s", got)
	}
	if !strings.Contains(got, "| — | Unavailable | — | — |") {
		t.Fatalf("major rates:\n%s", got)
	}
	if !strings.Contains(got, "not $0 / 0%") {
		t.Fatalf("missing unavailable explanation:\n%s", got)
	}
}

func TestWriteReserveSection_nullTotalsRatesOnly(t *testing.T) {
	var b strings.Builder
	writeReserveSection(&b, pipeline.ReserveFigures{
		UnpricedCoins: []string{},
		MajorCoinReserveRates: []pipeline.MajorCoinReserveRate{{
			Coin: "BTC", ReserveRate: "1.1", Basis: "reported", SelfReported: true,
		}},
		AmountsPublished: false,
	})
	got := b.String()
	if strings.Contains(got, "`null`") {
		t.Fatalf("human report must say Unavailable, not null:\n%s", got)
	}
	if !strings.Contains(got, "Unavailable (snapshot did not publish amounts)") {
		t.Fatalf("liability/reserve cells:\n%s", got)
	}
}

func floatPtr(v float64) *float64 { return &v }
