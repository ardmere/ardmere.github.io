package onchainconfig_test

import (
	"testing"

	"github.com/ardmere/ardmere/internal/verifier/onchainconfig"
)

func TestBinanceOnchainConfigLoads(t *testing.T) {
	cfg, err := onchainconfig.ForExchange("binance")
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Native) < 2 {
		t.Fatalf("native pairs: %d", len(cfg.Native))
	}
	if len(cfg.Tokens) < 50 {
		t.Fatalf("token pairs: %d", len(cfg.Tokens))
	}
	usdt, ok := cfg.Tokens["USDT|ETH"]
	if !ok || usdt.Decimals != 6 {
		t.Fatalf("USDT|ETH: %+v ok=%v", usdt, ok)
	}
}
