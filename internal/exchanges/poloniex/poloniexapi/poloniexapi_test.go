package poloniexapi_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ardmere/ardmere/internal/exchanges/poloniex/poloniexapi"
)

func fixturePath(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "fixtures", "poloniex", name)
}

func TestParseSummaryBytesAPIEnvelope(t *testing.T) {
	raw := []byte(`{
		"fetchedAt":"2026-11-06T12:00:00Z",
		"source":"proxy-sapi-public",
		"version":"20261001",
		"info":{"version":"20261001"},
		"coinList":{
			"code":200,
			"message":"Success",
			"data":{
				"totalReserveRate":"1.0047",
				"merkleHash":"6980318bc8260fbc0d4ba0e69ff480614d61c8e0617c34e0e63c0305a2e9fb11",
				"balanceList":[
					{"currency":"BTC","userBalance":"12521.91","walletBalance":"12534.98","reserveRate":"1.0010"},
					{"currency":"ETH","currencyGroup":"ETH+stETH","userBalance":"1","walletBalance":"1.01","reserveRate":"1.01"}
				]
			}
		}
	}`)
	got, err := poloniexapi.ParseSummaryBytes(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.AuditID != "20261001" {
		t.Fatalf("audit id: %s", got.AuditID)
	}
	if got.MerkleRoot == "" {
		t.Fatal("missing merkle root")
	}
	if got.TotalReserveRate < 100 || got.TotalReserveRate > 101 {
		t.Fatalf("total rate should be percent-scale, got %v", got.TotalReserveRate)
	}
	if len(got.CoinRows) != 2 || got.CoinRows[0].Coin != "BTC" {
		t.Fatalf("rows: %+v", got.CoinRows)
	}
	if got.CoinRows[0].ReserveRate < 100 {
		t.Fatalf("coin rate should be percent-scale, got %v", got.CoinRows[0].ReserveRate)
	}
}

func TestParseSummaryFixture(t *testing.T) {
	raw, err := os.ReadFile(fixturePath("20261001-summary.json"))
	if err != nil {
		t.Fatal(err)
	}
	got, err := poloniexapi.ParseSummaryBytes(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.AuditID != "20261001" || got.MerkleRoot == "" || len(got.CoinRows) < 5 {
		t.Fatalf("unexpected: %+v", got)
	}
	if got.TotalReserveRate < 100 {
		t.Fatalf("expected percent-scale total rate, got %v", got.TotalReserveRate)
	}
}
