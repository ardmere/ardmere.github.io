package onchainconfig_test

import (
	"testing"

	"github.com/ardmere/ardmere/internal/rpc"
	"github.com/ardmere/ardmere/internal/verifier/onchainconfig"
)

func TestOKXNetworkCoinAliases(t *testing.T) {
	onchainconfig.ResetForTest()
	cfg, err := onchainconfig.ForExchange("okx")
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{
		"TRX|TRON",
		"USDT-TRC20|TRON",
		"USDT|TRON",
		"USDT|POLYGON",
		"ETH|ZKSYNC2",
		"OKB-X1-USDC|XLAYER",
		"FIL-EVM|FEVM",
	} {
		if _, ok := cfg.Tokens[key]; !ok {
			if _, ok := cfg.Native[key]; !ok {
				t.Fatalf("missing aliased key %q", key)
			}
		}
	}
	if spec, ok := cfg.Tokens["TRX|TRON"]; !ok || !spec.Native || spec.Net != rpc.NetTron {
		t.Fatalf("TRX|TRON: %+v ok=%v", spec, ok)
	}
	if spec, ok := cfg.Tokens["USDT-TRC20|TRON"]; !ok || spec.Contract == "" {
		t.Fatalf("USDT-TRC20|TRON: %+v ok=%v", spec, ok)
	}

	onchainconfig.ResetForTest()
	ledger, err := onchainconfig.LedgerForExchange("okx")
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{
		"DOGE|DOGE",
		"APTOS|APTOS",
		"USDT-APTOS|APTOS",
		"USDC-SUI|SUI",
		"FIL|FIL",
		"TONCOIN-NEW|TON",
	} {
		if _, ok := ledger[key]; !ok {
			t.Fatalf("missing ledger key %q", key)
		}
	}
}
