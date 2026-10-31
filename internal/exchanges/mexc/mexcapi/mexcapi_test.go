package mexcapi_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ardmere/ardmere/internal/exchanges/mexc/mexcapi"
)

func fixturePath(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "fixtures", "mexc", name)
}

func TestParseSummaryBytesBundle(t *testing.T) {
	raw := []byte(`{
		"fetchedAt":"2026-07-15T10:00:00Z",
		"source":"press-release-scrape",
		"info":{"audit_id":"202607","audit_time":"2026-07-15 00:00:00","totalReserveRate":114,"merkleRoot":""},
		"coinList":{"list":[
			{"coin":"BTC","reserveRate":281,"userAssets":4439.51,"platformAssets":12474.0231},
			{"coin":"USDT","reserveRate":119,"userAssets":1811076206.47,"platformAssets":2155180685.6993}
		]}
	}`)
	got, err := mexcapi.ParseSummaryBytes(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.AuditID != "202607" || got.TotalReserveRate != 114 {
		t.Fatalf("unexpected summary: %+v", got)
	}
	if len(got.CoinRows) != 2 || got.CoinRows[0].Coin != "BTC" || got.CoinRows[0].ReserveRate != 281 {
		t.Fatalf("unexpected rows: %+v", got.CoinRows)
	}
}

func TestParseSummaryFixture(t *testing.T) {
	raw, err := os.ReadFile(fixturePath("202607-summary.json"))
	if err != nil {
		t.Fatal(err)
	}
	got, err := mexcapi.ParseSummaryBytes(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.AuditID != "202607" {
		t.Fatalf("audit id: %s", got.AuditID)
	}
	if len(got.CoinRows) < 4 {
		t.Fatalf("expected 4 coin rows, got %d", len(got.CoinRows))
	}
}
