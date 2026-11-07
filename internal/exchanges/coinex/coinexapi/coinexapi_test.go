package coinexapi_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ardmere/ardmere/internal/exchanges/coinex/coinexapi"
)

func fixture(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "fixtures", "coinex", name)
}

func TestParseSummaryFixture(t *testing.T) {
	raw, err := os.ReadFile(fixture("20260914-summary.json"))
	if err != nil {
		t.Fatal(err)
	}
	got, err := coinexapi.ParseSummaryBytes(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.AuditID != "20260914" {
		t.Fatalf("AuditID=%q", got.AuditID)
	}
	if got.MerkleRoot == "" {
		t.Fatal("missing merkle root")
	}
	if len(got.CoinRows) == 0 {
		t.Fatal("expected coin rows")
	}
	if got.TotalReserveRate < 100 {
		t.Fatalf("unexpected min reserve rate %v", got.TotalReserveRate)
	}
}
