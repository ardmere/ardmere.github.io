package rpc

import (
	"context"
	"os"
	"testing"
)

func TestParseSolanaHistoryBalanceRaw(t *testing.T) {
	raw, err := parseSolanaHistoryBalanceRaw([]byte(`{
		"balanceRaw":"27790751441",
		"decimals":6,
		"slot":423478906
	}`))
	if err != nil || raw != 27790751441 {
		t.Fatalf("got %d err=%v", raw, err)
	}
	if _, err := parseSolanaHistoryBalanceRaw([]byte(`{"error":"Unauthorized"}`)); err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadSolanaHistoryConfig(t *testing.T) {
	t.Setenv("HELIUS_API_KEY", "")
	t.Setenv("SOLANA_INDEX_API_KEY", "")
	t.Setenv("SOLANA_HISTORY_API_KEY", "")
	if _, ok := loadSolanaHistoryConfig(); ok {
		t.Fatal("expected disabled without key")
	}

	t.Setenv("HELIUS_API_KEY", "helius-test")
	cfg, ok := loadSolanaHistoryConfig()
	if !ok || cfg.Provider != solanaHistoryProviderHelius || cfg.APIKey != "helius-test" || cfg.BaseURL != "https://api.helius.xyz" {
		t.Fatalf("helius cfg=%+v ok=%v", cfg, ok)
	}

	t.Setenv("HELIUS_API_KEY", "")
	t.Setenv("SOLANA_INDEX_API_KEY", "test-key")
	cfg, ok = loadSolanaHistoryConfig()
	if !ok || cfg.Provider != solanaHistoryProviderSolanaIndex || cfg.APIKey != "test-key" || cfg.BaseURL != "https://solanaindex.top" {
		t.Fatalf("solanaindex cfg=%+v ok=%v", cfg, ok)
	}

	t.Setenv("SOLANA_HISTORY_API_URL", "https://example.test")
	cfg, ok = loadSolanaHistoryConfig()
	if !ok || cfg.BaseURL != "https://example.test" {
		t.Fatalf("custom base: %+v", cfg)
	}
}

func TestSolanaHistoryProviderFromUsed(t *testing.T) {
	if got := SolanaHistoryProviderFromUsed("https://api.helius.xyz/v1/wallet"); got != "helius" {
		t.Fatalf("got %q", got)
	}
	if got := SolanaHistoryProviderFromUsed("https://solanaindex.top/api/v1"); got != "solanaindex" {
		t.Fatalf("got %q", got)
	}
}

func TestHeliusBalanceAtLive(t *testing.T) {
	if os.Getenv("HELIUS_API_KEY") == "" {
		t.Skip("HELIUS_API_KEY not set")
	}
	t.Setenv("SOLANA_INDEX_API_KEY", "")
	t.Setenv("SOLANA_HISTORY_API_KEY", "")

	const (
		addr = "Ff1wYqLwVQou6T2VbSNNUbF7YimdJ34qFvAHKJwLc9nL"
		mint = "ukHH6c7mMyiWCf1b9pnWe25TSpkDDt3H5pQZgZ74J82"
		slot = int64(423478906)
	)

	c := NewLedger()
	raw, used, historical, err := c.SolanaSPLBalanceAtSlot(context.Background(), addr, mint, slot)
	if err != nil {
		t.Fatalf("SolanaSPLBalanceAtSlot: %v (provider=%s)", err, used)
	}
	if !historical {
		t.Fatalf("expected historical=true, provider=%s", used)
	}
	if used == "" || SolanaHistoryProviderFromUsed(used) != "helius" {
		t.Fatalf("unexpected provider label: %q", used)
	}
	t.Logf("OK helius BOME balanceRaw=%d at slot=%d via %s", raw, slot, used)
}
