package bitgetapi_test

import (
	"testing"

	"github.com/ardmere/ardmere/internal/exchanges/bitget/bitgetapi"
)

func TestParseSummaryBytesBundle(t *testing.T) {
	raw := []byte(`{
		"fetchedAt":"2026-06-14T13:40:00Z",
		"info":{"audit_id":"202605","audit_time":"2026-05-27 10:00:00","merkelRootHash":"ca89456bb711c913","totalReserveRate":127},
		"coinList":{"list":[{"coin":"BTC","platformAssets":323.4,"userAssets":300,"reserveRate":107.8}]}
	}`)
	got, err := bitgetapi.ParseSummaryBytes(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.AuditID != "202605" || got.MerkleRoot != "ca89456bb711c913" || got.TotalReserveRate != 127 {
		t.Fatalf("unexpected summary: %+v", got)
	}
	if len(got.CoinRows) != 1 || got.CoinRows[0].Coin != "BTC" || got.CoinRows[0].ReserveRate != 107.8 {
		t.Fatalf("unexpected rows: %+v", got.CoinRows)
	}
}
