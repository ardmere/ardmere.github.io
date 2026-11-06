package bitunixapi_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ardmere/ardmere/internal/exchanges/bitunix/bitunixapi"
)

func fixture(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "fixtures", "bitunix", name)
}

func TestParseSummaryFixture(t *testing.T) {
	raw, err := os.ReadFile(fixture("20261103-summary.json"))
	if err != nil {
		t.Fatal(err)
	}
	got, err := bitunixapi.ParseSummaryBytes(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.AuditID != "20261103" {
		t.Fatalf("AuditID=%q", got.AuditID)
	}
	if got.MerkleRoot != "d21a6f9c4e47d35f3083086eef4344e3d91cc627c866f9d5ff9228c12d9bd16b" {
		t.Fatalf("MerkleRoot=%q", got.MerkleRoot)
	}
	if got.TotalReserveRate != 105 {
		t.Fatalf("TotalReserveRate=%v (want min=105)", got.TotalReserveRate)
	}
	if len(got.CoinRows) != 3 {
		t.Fatalf("CoinRows=%d", len(got.CoinRows))
	}
}
