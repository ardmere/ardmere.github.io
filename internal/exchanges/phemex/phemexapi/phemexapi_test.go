package phemexapi

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseSummaryBytesFixture(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "fixtures", "phemex", "20261101-summary.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	parsed, err := ParseSummaryBytes(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if parsed.AuditID != "20261101" {
		t.Fatalf("audit id: got %q", parsed.AuditID)
	}
	if parsed.OfficialVersionID != "2026-11-01" {
		t.Fatalf("official version: got %q", parsed.OfficialVersionID)
	}
	if len(parsed.CoinRows) < 8 {
		t.Fatalf("expected >=8 coin rows, got %d", len(parsed.CoinRows))
	}
	if parsed.TotalReserveRate < 100 {
		t.Fatalf("min reserve rate should be >=100, got %v", parsed.TotalReserveRate)
	}
	foundBTC := false
	for _, c := range parsed.CoinRows {
		if c.Coin == "BTC" {
			foundBTC = true
			if c.ReserveRate < 100 {
				t.Fatalf("BTC reserve rate <100: %v", c.ReserveRate)
			}
			if c.ReserveAmount <= 0 || c.LiabilityAmount <= 0 {
				t.Fatalf("BTC balances missing: reserve=%v liability=%v", c.ReserveAmount, c.LiabilityAmount)
			}
		}
	}
	if !foundBTC {
		t.Fatal("BTC row missing")
	}
}
