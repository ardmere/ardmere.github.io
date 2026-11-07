package ourbitapi_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ardmere/ardmere/internal/exchanges/ourbit/ourbitapi"
)

func fixturePath(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "fixtures", "ourbit", name)
}

func TestParseSummaryFixture(t *testing.T) {
	raw, err := os.ReadFile(fixturePath("20261101-summary.json"))
	if err != nil {
		t.Fatal(err)
	}
	got, err := ourbitapi.ParseSummaryBytes(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.AuditID != "20261101" {
		t.Fatalf("audit id: %s", got.AuditID)
	}
	if len(got.CoinRows) < 3 {
		t.Fatalf("expected >=3 coin rows, got %d", len(got.CoinRows))
	}
	if got.TotalReserveRate < 100 {
		t.Fatalf("expected totalReserveRate >= 100, got %v", got.TotalReserveRate)
	}
	foundBTC := false
	for _, c := range got.CoinRows {
		if c.Coin == "BTC" && c.ReserveRate > 100 {
			foundBTC = true
		}
	}
	if !foundBTC {
		t.Fatalf("expected BTC row with rate>100: %+v", got.CoinRows)
	}
	if len(got.WalletRows) < 1 {
		t.Fatalf("expected archived wallet rows, got %d", len(got.WalletRows))
	}
}
