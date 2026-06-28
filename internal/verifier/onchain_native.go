package verifier

import (
	"sync"

	"github.com/ardmere/ardmere/internal/rpc"
	"github.com/ardmere/ardmere/internal/verifier/onchainconfig"
)

var (
	nativeMu    sync.Mutex
	nativeCache = map[string]map[string]rpc.Network{}
)

// SetNativeExchange is deprecated; pass Exchange on OnchainBalanceHot instead.
func SetNativeExchange(id string) {
	nativeMu.Lock()
	delete(nativeCache, id)
	nativeMu.Unlock()
	loadNativeSupported(id)
}

func loadNativeSupported(exchangeID string) map[string]rpc.Network {
	nativeMu.Lock()
	defer nativeMu.Unlock()
	if m, ok := nativeCache[exchangeID]; ok {
		return m
	}
	cfg, err := onchainconfig.ForExchange(exchangeID)
	if err != nil {
		panic("load onchain native config for " + exchangeID + ": " + err.Error())
	}
	nativeCache[exchangeID] = cfg.Native
	return cfg.Native
}
