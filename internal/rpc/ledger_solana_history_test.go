package rpc

import "testing"

func TestParseSolanaIndexBalanceRaw(t *testing.T) {
	raw, err := parseSolanaIndexBalanceRaw([]byte(`{
		"balanceRaw":"27790751441",
		"decimals":6,
		"slot":423478906
	}`))
	if err != nil || raw != 27790751441 {
		t.Fatalf("got %d err=%v", raw, err)
	}
	if _, err := parseSolanaIndexBalanceRaw([]byte(`{"error":"Unauthorized"}`)); err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadSolanaHistoryConfig(t *testing.T) {
	t.Setenv("SOLANA_INDEX_API_KEY", "")
	t.Setenv("SOLANA_HISTORY_API_KEY", "")
	if _, ok := loadSolanaHistoryConfig(); ok {
		t.Fatal("expected disabled without key")
	}
	t.Setenv("SOLANA_INDEX_API_KEY", "test-key")
	cfg, ok := loadSolanaHistoryConfig()
	if !ok || cfg.APIKey != "test-key" || cfg.BaseURL != "https://solanaindex.top" {
		t.Fatalf("cfg=%+v ok=%v", cfg, ok)
	}
	t.Setenv("SOLANA_HISTORY_API_URL", "https://example.test")
	cfg, ok = loadSolanaHistoryConfig()
	if !ok || cfg.BaseURL != "https://example.test" {
		t.Fatalf("custom base: %+v", cfg)
	}
}
