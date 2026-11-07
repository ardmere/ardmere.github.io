package krakenapi

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseSummaryBytesFixture(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "fixtures", "kraken", "20260630-summary.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	parsed, err := ParseSummaryBytes(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if parsed.AuditID != "20260630" {
		t.Fatalf("audit id: got %q", parsed.AuditID)
	}
	if parsed.ReportDate != "2026-06-30" {
		t.Fatalf("report date: got %q", parsed.ReportDate)
	}
	if parsed.MerkleRoot != "56272172594d525355c22220eda22e1f1f842c37455d49f0588b91ed97362ad1" {
		t.Fatalf("merkle root: got %q", parsed.MerkleRoot)
	}
	if len(parsed.MerkleRoot) != 64 {
		t.Fatalf("expected full 64-hex root, got len=%d", len(parsed.MerkleRoot))
	}
	if parsed.Accountant == "" || parsed.Accountant == "unpublished" {
		t.Fatalf("accountant: got %q", parsed.Accountant)
	}
	if len(parsed.CoinRows) < 6 {
		t.Fatalf("expected >=6 coin rows, got %d", len(parsed.CoinRows))
	}
	// Min effective rate should be ADA 100.3 (capped stables floor at 105).
	if parsed.TotalReserveRate != 100.3 {
		t.Fatalf("min reserve rate: want 100.3, got %v", parsed.TotalReserveRate)
	}
	byCoin := map[string]CoinRow{}
	for _, c := range parsed.CoinRows {
		byCoin[c.Coin] = c
	}
	if byCoin["BTC"].ReserveRate != 102.9 {
		t.Fatalf("BTC rate: got %v", byCoin["BTC"].ReserveRate)
	}
	if !byCoin["USDC"].CapPercentage || byCoin["USDC"].ReserveRate != 105 {
		t.Fatalf("USDC capped floor: %+v", byCoin["USDC"])
	}
	if byCoin["USDC"].DisplayRatio != "105%+" {
		t.Fatalf("USDC display: got %q", byCoin["USDC"].DisplayRatio)
	}
}

func TestNormalizeAuditDate(t *testing.T) {
	got, err := normalizeAuditDate("20260630")
	if err != nil || got.Format("2006-01-02") != "2026-06-30" {
		t.Fatalf("compact: got %v err=%v", got, err)
	}
	got, err = normalizeAuditDate("2026-06-30")
	if err != nil || got.Format("2006-01-02") != "2026-06-30" {
		t.Fatalf("dashed: got %v err=%v", got, err)
	}
}

func TestParseAuditLabel(t *testing.T) {
	got, err := parseAuditLabel("June 30, 2026")
	if err != nil || got.Format("2006-01-02") != "2026-06-30" {
		t.Fatalf("label: got %v err=%v", got, err)
	}
}
