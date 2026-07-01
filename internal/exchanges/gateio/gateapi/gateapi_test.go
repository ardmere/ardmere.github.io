package gateapi

import (
	"encoding/json"
	"testing"
	"time"
)

func TestParseSummaryBundle(t *testing.T) {
	info := []byte(`{"code":0,"message":"success","data":{
		"audit_id":"20231110",
		"audit_time":"2023-11-10 00:00:00",
		"merkle_root_hash":"0xabc123",
		"total_reserve_rate":115.69,
		"total_reserve_amount":3206000000,
		"customer_net_balance":2770000000,
		"excess_reserve_value":436000000
	}}`)
	coins := []byte(`{"code":0,"data":{"list":[
		{"coin":"BTC","reserve_rate":123.52,"reserve_amount":17100,"customer_liability":13800},
		{"coin":"ETH","reserve_rate":110.0,"reserve_amount":500000,"customer_liability":450000}
	]}}`)

	got, err := ParseSummaryBundle(info, coins)
	if err != nil {
		t.Fatal(err)
	}
	if got.MerkleRoot != "0xabc123" || got.TotalReserveRate != 115.69 {
		t.Fatalf("summary: %+v", got)
	}
	if len(got.CoinRows) != 2 || got.CoinRows[0].Coin != "BTC" {
		t.Fatalf("coins: %+v", got.CoinRows)
	}
	if got.AuditTime.IsZero() {
		t.Fatalf("expected parsed audit time")
	}
}

func TestParseSummaryBundleWrapped(t *testing.T) {
	wrapped, _ := json.Marshal(SummaryBundle{
		FetchedAt: time.Now().UTC(),
		Info: json.RawMessage(`{"code":0,"data":{
			"audit_time":"2024-01-03T00:00:00Z",
			"merkle_root_hash":"0xdead",
			"total_reserve_rate":101.0,
			"total_reserve_amount":1,
			"customer_net_balance":1,
			"excess_reserve_value":0
		}}`),
	})
	got, err := ParseSummaryBytes(wrapped)
	if err != nil {
		t.Fatal(err)
	}
	if got.MerkleRoot != "0xdead" {
		t.Fatalf("got %+v", got)
	}
}
