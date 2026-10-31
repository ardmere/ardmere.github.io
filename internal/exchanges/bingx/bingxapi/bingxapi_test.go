package bingxapi_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ardmere/ardmere/internal/exchanges/bingx/bingxapi"
)

func fixturePath(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "fixtures", "bingx", name)
}

func TestParseSummaryFixture(t *testing.T) {
	raw, err := os.ReadFile(fixturePath("20261015-summary.json"))
	if err != nil {
		t.Fatal(err)
	}
	got, err := bingxapi.ParseSummaryBytes(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.AuditID != "20261015" || got.TotalReserveRate != 123.93 {
		t.Fatalf("unexpected: %+v", got)
	}
	if len(got.CoinRows) != 4 || got.CoinRows[0].ReserveRate != 143.42 {
		t.Fatalf("rows: %+v", got.CoinRows)
	}
}
