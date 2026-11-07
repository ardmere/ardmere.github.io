package wooxapi

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseSummaryBytesFixture(t *testing.T) {
	path := filepath.Join("..", "..", "..", "..", "fixtures", "woox", "20261015-summary.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	parsed, err := ParseSummaryBytes(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if parsed.AuditID != "20261015" {
		t.Fatalf("audit id: got %q want 20261015", parsed.AuditID)
	}
	if parsed.OfficialAuditID != "1b80d170-770a-4ed3-b36d-b514e1831271" {
		t.Fatalf("official audit id: %q", parsed.OfficialAuditID)
	}
	if parsed.MerkleTreeURL == "" {
		t.Fatal("missing merkle tree URL")
	}
	if parsed.TotalReserveRate < 90 || parsed.TotalReserveRate > 110 {
		t.Fatalf("unexpected reserve ratio %v", parsed.TotalReserveRate)
	}
	if parsed.TotalAssetsUSDT <= 0 || parsed.TotalLiabilitiesUSDT <= 0 {
		t.Fatalf("missing totals assets=%v liab=%v", parsed.TotalAssetsUSDT, parsed.TotalLiabilitiesUSDT)
	}
}

func TestReserveRatioRounding(t *testing.T) {
	got := reserveRatio(48046392.98871899, 49455631.19857185)
	// 97.151… → UI toFixed(1) ≈ 97.2
	if got != 97.2 {
		t.Fatalf("ratio: got %v want 97.2", got)
	}
}
