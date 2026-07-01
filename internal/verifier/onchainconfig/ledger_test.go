package onchainconfig_test

import (
	"testing"

	"github.com/ardmere/ardmere/internal/verifier/onchainconfig"
)

func TestLedgerForExchangeBinance(t *testing.T) {
	m, err := onchainconfig.LedgerForExchange("binance")
	if err != nil {
		t.Fatal(err)
	}
	if len(m) < 30 {
		t.Fatalf("expected >=30 ledger entries, got %d", len(m))
	}
	spec, ok := m["APT|APT"]
	if !ok || spec.Kind != onchainconfig.LedgerAptos {
		t.Fatalf("APT|APT: %+v", m["APT|APT"])
	}
	spec, ok = m["BTC|BTC"]
	if !ok || spec.Kind != onchainconfig.LedgerEsplora {
		t.Fatalf("BTC|BTC: %+v", spec)
	}
	spec, ok = m["XRP|XRP"]
	if !ok || spec.Kind != onchainconfig.LedgerXRPL || spec.Decimals != 6 {
		t.Fatalf("XRP|XRP: %+v", spec)
	}
}
