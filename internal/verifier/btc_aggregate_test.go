package verifier

import (
	"testing"

	"github.com/ardmere/ardmere/internal/walletzip"
	"github.com/shopspring/decimal"
)

func TestClassifyBalanceMismatch(t *testing.T) {
	claim := decimal.RequireFromString("10")
	surplus := decimal.RequireFromString("11")
	short := decimal.RequireFromString("9")
	if w, s := classifyBalanceMismatch(surplus, claim); w || !s {
		t.Fatalf("expected surplus pass classification")
	}
	if w, s := classifyBalanceMismatch(short, claim); w || s {
		t.Fatalf("expected shortage fail classification, within=%v surplus=%v", w, s)
	}
	if w, _ := classifyBalanceMismatch(claim, claim); !w {
		t.Fatalf("expected match")
	}
}

func TestBtcEthNativeCustodyUnverifiable(t *testing.T) {
	row := walletzip.Row{Coin: "BTC", Network: "ETH", Balance: decimal.RequireFromString("699.2371")}
	ok, note := btcEthNativeCustodyUnverifiable(row, decimal.Zero)
	if !ok || note == "" {
		t.Fatalf("expected unverifiable for zero WBTC with large claim")
	}
	row.Balance = decimal.RequireFromString("0.001")
	if ok, _ := btcEthNativeCustodyUnverifiable(row, decimal.Zero); ok {
		t.Fatalf("expected verifiable for dust claim")
	}
	if ok, _ := btcEthNativeCustodyUnverifiable(row, decimal.RequireFromString("1")); ok {
		t.Fatalf("expected verifiable when WBTC present")
	}
}

func TestStablecoinEthOmnibusMismatch(t *testing.T) {
	row := walletzip.Row{
		Coin:    "USDC",
		Network: "ETH",
		Balance: decimal.RequireFromString("61874569.163224"),
	}
	ok, note := stablecoinEthOmnibusMismatch(row, decimal.RequireFromString("1.324019"))
	if !ok || note == "" {
		t.Fatal("expected mega USDC omnibus pattern")
	}
	row.Balance = decimal.RequireFromString("100000")
	if ok, _ := stablecoinEthOmnibusMismatch(row, decimal.Zero); ok {
		t.Fatal("claim below 1M threshold should not trigger")
	}
	row.Balance = decimal.RequireFromString("5000000")
	if ok, _ := stablecoinEthOmnibusMismatch(row, decimal.RequireFromString("6000000")); ok {
		t.Fatal("surplus should not trigger")
	}
}

func TestTronSnapshotMismatchNote(t *testing.T) {
	if got := tronSnapshotMismatchNote("USDT-TRC20", "TRON"); got == "" {
		t.Fatal("expected Tron note for USDT-TRC20|TRON")
	}
	if got := tronSnapshotMismatchNote("TRX", "TRON"); got == "" {
		t.Fatal("expected Tron native note")
	}
	if got := tronSnapshotMismatchNote("USDT", "ETH"); got != "" {
		t.Fatal("expected no Tron note on ETH network")
	}
}
