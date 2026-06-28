package verifier

import (
	"sync"

	"github.com/ardmere/ardmere/internal/verifier/onchainconfig"
)

type tokenSpec = onchainconfig.TokenSpec

var (
	tokenMu    sync.RWMutex
	tokenCache = map[string]map[string]tokenSpec{}
)

// SetTokenExchange resets cached token maps (tests).
func SetTokenExchange(id string) {
	tokenMu.Lock()
	delete(tokenCache, id)
	tokenMu.Unlock()
}

func loadTokenSupported() map[string]tokenSpec {
	return loadTokenSupportedFor("binance")
}

func loadTokenSupportedFor(exchangeID string) map[string]tokenSpec {
	if exchangeID == "" {
		exchangeID = "binance"
	}
	tokenMu.RLock()
	if m, ok := tokenCache[exchangeID]; ok {
		tokenMu.RUnlock()
		return m
	}
	tokenMu.RUnlock()

	tokenMu.Lock()
	defer tokenMu.Unlock()
	if m, ok := tokenCache[exchangeID]; ok {
		return m
	}
	cfg, err := onchainconfig.ForExchange(exchangeID)
	if err != nil {
		panic("load onchain token config for " + exchangeID + ": " + err.Error())
	}
	tokenCache[exchangeID] = cfg.Tokens
	return cfg.Tokens
}
