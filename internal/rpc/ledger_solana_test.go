package rpc

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestLoadSolanaProvidersFromConfig(t *testing.T) {
	t.Setenv("SOLANA_RPC", "")
	providers := loadSolanaProviders()
	if len(providers) < 2 {
		t.Fatalf("expected multiple solana providers, got %d", len(providers))
	}
	if providers[0].URL == "https://api.mainnet-beta.solana.com" {
		t.Fatalf("official solana RPC should not be first provider: %v", providers[0].URL)
	}
}

func TestLoadSolanaProvidersEnvOverride(t *testing.T) {
	t.Setenv("SOLANA_RPC", "https://example.invalid/solana")
	providers := loadSolanaProviders()
	if len(providers) != 1 || providers[0].URL != "https://example.invalid/solana" {
		t.Fatalf("SOLANA_RPC override: %+v", providers)
	}
}

func TestSolanaLiveProbe(t *testing.T) {
	if os.Getenv("SOLANA_LIVE_PROBE") != "1" {
		t.Skip("set SOLANA_LIVE_PROBE=1 to run")
	}
	c := NewLedger()
	// Binance hot wallet (public, high activity)
	addr := "9WzDXwBbmkg8ZTbNMqUxvQRAyrZzDsGYdLVL9zYtAWWM"
	lamports, used, err := c.SolanaNativeBalance(context.Background(), addr)
	if err != nil {
		t.Fatalf("SolanaNativeBalance: %v (provider=%s)", err, used)
	}
	if lamports == 0 {
		t.Fatalf("expected non-zero lamports")
	}
	if used == "" || strings.Contains(used, "mainnet-beta.solana.com") {
		t.Logf("provider=%s lamports=%d", used, lamports)
	}
}
