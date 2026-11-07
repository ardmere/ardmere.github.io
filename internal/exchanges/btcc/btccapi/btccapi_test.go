package btccapi

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseSummaryBytesFixture(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "fixtures", "btcc", "20261015-summary.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	parsed, err := ParseSummaryBytes(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if parsed.AuditID != "20261015" {
		t.Fatalf("audit id: got %q", parsed.AuditID)
	}
	if parsed.ReportDate != "2026-10-15" {
		t.Fatalf("report date: got %q", parsed.ReportDate)
	}
	if parsed.MerkleRoot != "182f650c44a74f6d" {
		t.Fatalf("merkle root: got %q", parsed.MerkleRoot)
	}
	if len(parsed.MerkleRoot) != 16 {
		t.Fatalf("expected truncated 16-hex root, got len=%d", len(parsed.MerkleRoot))
	}
	if len(parsed.CoinRows) != 6 {
		t.Fatalf("expected 6 coin rows, got %d", len(parsed.CoinRows))
	}
	if parsed.TotalReserveRate != 111 {
		t.Fatalf("min reserve rate (USDC): want 111, got %v", parsed.TotalReserveRate)
	}
	wantOrder := []string{"BTC", "ETH", "USDT", "USDC", "XRP", "ADA"}
	for i, coin := range wantOrder {
		if parsed.CoinRows[i].Coin != coin {
			t.Fatalf("coin order[%d]: want %s got %s", i, coin, parsed.CoinRows[i].Coin)
		}
		if parsed.CoinRows[i].ReserveRate < 100 {
			t.Fatalf("%s reserve rate <100: %v", coin, parsed.CoinRows[i].ReserveRate)
		}
	}
}

func TestNormalizeReportDate(t *testing.T) {
	got, err := normalizeReportDate("20261015")
	if err != nil || got != "2026-10-15" {
		t.Fatalf("compact: got %q err=%v", got, err)
	}
	got, err = normalizeReportDate("2026-10-15")
	if err != nil || got != "2026-10-15" {
		t.Fatalf("dashed: got %q err=%v", got, err)
	}
}

func TestParsePercent(t *testing.T) {
	v, err := parsePercent("147%")
	if err != nil || v != 147 {
		t.Fatalf("147%%: got %v err=%v", v, err)
	}
	v, err = parsePercent("111")
	if err != nil || v != 111 {
		t.Fatalf("111: got %v err=%v", v, err)
	}
}
