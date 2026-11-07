package toobitapi

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseSummaryFixture(t *testing.T) {
	candidates := []string{
		filepath.Join("..", "..", "..", "..", "fixtures", "toobit", "20261001-summary.json"),
		"fixtures/toobit/20261001-summary.json",
	}
	var raw []byte
	var err error
	for _, p := range candidates {
		raw, err = os.ReadFile(p)
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	parsed, err := ParseSummaryBytes(raw)
	if err != nil {
		t.Fatalf("ParseSummaryBytes: %v", err)
	}
	if parsed.AuditID != "20261001" {
		t.Fatalf("AuditID = %q, want 20261001", parsed.AuditID)
	}
	if parsed.OfficialAuditID != "Toobit_POR_202610" {
		t.Fatalf("OfficialAuditID = %q", parsed.OfficialAuditID)
	}
	wantRoot := "2a406020be83997ddae3359fea169eb97d51f5524c88f3bf05851a6bcf8ec599"
	if parsed.MerkleRoot != wantRoot {
		t.Fatalf("MerkleRoot = %s", parsed.MerkleRoot)
	}
	if parsed.TotalReserveRate != 104 {
		t.Fatalf("TotalReserveRate = %v, want 104 (min)", parsed.TotalReserveRate)
	}
	if len(parsed.CoinRows) != 4 {
		t.Fatalf("CoinRows len = %d", len(parsed.CoinRows))
	}
}

func TestCodeOKAccepts200(t *testing.T) {
	if !codeOK(200) || !codeOK(float64(200)) || !codeOK("200") {
		t.Fatal("expected code 200 to be OK for Toobit API")
	}
}
