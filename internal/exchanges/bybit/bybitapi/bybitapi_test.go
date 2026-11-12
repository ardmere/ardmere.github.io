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
	if got.AuditID != "2025061709" || got.TotalReserveRate == nil || *got.TotalReserveRate != 102.5 {
		t.Fatalf("audit/rate: %+v", got)
	}
	if got.MerkleRoot != "c91d5de0554f97244e4d9f8056fad70fa0cb2cdb23c290ec597a042645dcbc03" {
		t.Fatalf("merkle: %s", got.MerkleRoot)
	}
	if len(got.CoinRows) != 2 || got.CoinRows[0].Coin != "BTC" {
		t.Fatalf("coins: %+v", got.CoinRows)
	}
	if got.CoinRows[0].ReserveAmount == nil || *got.CoinRows[0].ReserveAmount != 1000 {
		t.Fatalf("btc reserve: %+v", got.CoinRows[0])
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
	if got.TotalReserveRate == nil || *got.TotalReserveRate != 100.1 {
		t.Fatalf("rate: %+v", got.TotalReserveRate)
	}
}

func TestParseSummaryBundle_nullReservesStayNil(t *testing.T) {
	info := []byte(`{"code":0,"data":{"audit_id":"2026102706","audit_time":"2026-10-27 06:00:00","period":41}}`)
	coins := []byte(`{"code":0,"data":{"list":[{"coin":"BTC","customer_liability":57373.99,"reserve_rate":null,"reserve_amount":null},{"coin":"ETH","customer_liability":1,"reserve_amount":0,"reserve_rate":0}]}}`)
	got, err := bybitapi.ParseSummaryBundle(info, coins)
	if err != nil {
		t.Fatal(err)
	}
	if got.TotalReserveRate != nil {
		t.Fatalf("total rate should be nil, got %v", *got.TotalReserveRate)
	}
	if len(got.CoinRows) != 2 {
		t.Fatalf("coins: %+v", got.CoinRows)
	}
	btc := got.CoinRows[0]
	if btc.ReserveAmount != nil || btc.ReserveRate != nil {
		t.Fatalf("btc missing reserves must be nil: %+v", btc)
	}
	if btc.LiabilityAmount == nil || *btc.LiabilityAmount != 57373.99 {
		t.Fatalf("btc liability: %+v", btc.LiabilityAmount)
	}
	eth := got.CoinRows[1]
	if eth.ReserveAmount == nil || *eth.ReserveAmount != 0 {
		t.Fatalf("explicit zero reserve must stay 0: %+v", eth.ReserveAmount)
	}
	if eth.ReserveRate == nil || *eth.ReserveRate != 0 {
		t.Fatalf("explicit zero rate must stay 0: %+v", eth.ReserveRate)
	}
}

func TestParseSummaryBundle_omittedReserveFieldsStayNil(t *testing.T) {
	info := []byte(`{"code":0,"data":{"audit_id":"2026082604","audit_time":"2026-08-26 04:00:00"}}`)
	coins := []byte(`{"code":0,"data":{"list":[{"coin":"USDT","customer_liability":100}]}}`)
	got, err := bybitapi.ParseSummaryBundle(info, coins)
	if err != nil {
		t.Fatal(err)
	}
	if got.CoinRows[0].ReserveAmount != nil || got.CoinRows[0].ReserveRate != nil {
		t.Fatalf("omitted reserves must be nil: %+v", got.CoinRows[0])
	}
}
