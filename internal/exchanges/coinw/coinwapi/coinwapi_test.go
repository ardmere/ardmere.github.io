package coinwapi

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseAnnouncementFixture(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "fixtures", "coinw", "202609-summary.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	parsed, err := ParseSummaryBytes(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if parsed.AuditID != "202609" {
		t.Fatalf("audit id: got %q", parsed.AuditID)
	}
	if parsed.TotalReserveRate != 109.65 {
		t.Fatalf("total rate: got %v", parsed.TotalReserveRate)
	}
	if parsed.MerkleRoot != "" {
		t.Fatalf("announcement fixture should not invent merkle root, got %q", parsed.MerkleRoot)
	}
	if parsed.AuditTime.IsZero() {
		t.Fatal("expected audit time")
	}
}

func TestParseAnnouncementRatio(t *testing.T) {
	html := `In the latest update, CoinW’s total reserve ratio stands at 109.65%.`
	rate, ok := parseAnnouncementRatio(html)
	if !ok || rate != 109.65 {
		t.Fatalf("got %v ok=%v", rate, ok)
	}
}

func TestCloudflareDetect(t *testing.T) {
	body := []byte(`<!DOCTYPE html><html><title>Just a moment...</title>`)
	if !isCloudflareChallenge(body, 403) {
		t.Fatal("expected waf detect")
	}
}
