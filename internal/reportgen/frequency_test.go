package reportgen

import (
	"path/filepath"
	"testing"
)

func TestCollectFrequency_BinanceJUN(t *testing.T) {
	reportsDir := filepath.Join("..", "..", "docs", "reports")
	freq := collectFrequency("binance", "PR01JUN26", "2026-06-01T00:00:00Z", reportsDir)

	if freq.LatestSnapshot != "2026-06-01T00:00:00Z" {
		t.Fatalf("latest: %s", freq.LatestSnapshot)
	}
	if freq.PreviousSnapshot != "2026-05-01T00:00:00Z" {
		t.Fatalf("previous: %s", freq.PreviousSnapshot)
	}
	if freq.ObservedCadence != "~monthly" {
		t.Fatalf("cadence: %s", freq.ObservedCadence)
	}
	if freq.HistoryAvailable != "3 snapshot(s) in public evaluation set" {
		t.Fatalf("history: %s", freq.HistoryAvailable)
	}
}

func TestCollectFrequency_BinanceAPR(t *testing.T) {
	reportsDir := filepath.Join("..", "..", "docs", "reports")
	freq := collectFrequency("binance", "PR01APR26", "2026-04-01T00:00:00Z", reportsDir)

	if freq.PreviousSnapshot != "UNVERIFIABLE" {
		t.Fatalf("previous: %s", freq.PreviousSnapshot)
	}
}
