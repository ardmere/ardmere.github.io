package bybitapi_test

import (
	"testing"

	"github.com/ardmere/ardmere/internal/exchanges/bybit/bybitapi"
)

func TestParseSummaryBundleGateStyle(t *testing.T) {
	info := []byte(`{"code":0,"data":{"audit_id":"2025061709","audit_time":"2025-06-17 09:00:00","merkle_root_hash":"c91d5de0554f97244e4d9f8056fad70fa0cb2cdb23c290ec597a042645dcbc03","total_reserve_rate":102.5}}`)
	coins := []byte(`{"code":0,"data":{"list":[{"coin":"BTC","reserve_rate":105.2,"reserve_amount":1000,"customer_liability":950},{"coin":"ETH","reserve_rate":101.0,"reserve_amount":5000,"customer_liability":4950}]}}`)
	got, err := bybitapi.ParseSummaryBundle(info, coins)
	if err != nil {
		t.Fatal(err)
	}
	if got.AuditID != "2025061709" || got.TotalReserveRate != 102.5 {
		t.Fatalf("audit/rate: %+v", got)
	}
	if got.MerkleRoot != "c91d5de0554f97244e4d9f8056fad70fa0cb2cdb23c290ec597a042645dcbc03" {
		t.Fatalf("merkle: %s", got.MerkleRoot)
	}
	if len(got.CoinRows) != 2 || got.CoinRows[0].Coin != "BTC" {
		t.Fatalf("coins: %+v", got.CoinRows)
	}
}

func TestParseSummaryBytesBundle(t *testing.T) {
	raw := []byte(`{"fetchedAt":"2026-06-14T00:00:00Z","info":{"code":0,"data":{"auditId":"2025061709","merkleRootHash":"abc","totalReserveRate":100.1}}}`)
	got, err := bybitapi.ParseSummaryBytes(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.AuditID != "2025061709" || got.MerkleRoot != "abc" {
		t.Fatalf("got %+v", got)
	}
}
