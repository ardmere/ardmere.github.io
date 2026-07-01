package verifier

import (
	"sync"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/verifier/onchainconfig"
	"github.com/ardmere/ardmere/internal/walletzip"
)

var (
	ledgerMu    sync.RWMutex
	ledgerCache = map[string]map[string]onchainconfig.LedgerSpec{}
)

// SetLedgerExchange clears cached ledger maps (tests).
func SetLedgerExchange(id string) {
	ledgerMu.Lock()
	delete(ledgerCache, id)
	ledgerMu.Unlock()
}

func loadLedgerSupported() map[string]onchainconfig.LedgerSpec {
	return loadLedgerSupportedFor("binance")
}

func loadLedgerSupportedFor(exchangeID string) map[string]onchainconfig.LedgerSpec {
	if exchangeID == "" {
		exchangeID = "binance"
	}
	ledgerMu.RLock()
	if m, ok := ledgerCache[exchangeID]; ok {
		ledgerMu.RUnlock()
		return m
	}
	ledgerMu.RUnlock()

	ledgerMu.Lock()
	defer ledgerMu.Unlock()
	if m, ok := ledgerCache[exchangeID]; ok {
		return m
	}
	cfg, err := onchainconfig.LedgerForExchange(exchangeID)
	if err != nil {
		ledgerCache[exchangeID] = map[string]onchainconfig.LedgerSpec{}
		return ledgerCache[exchangeID]
	}
	ledgerCache[exchangeID] = cfg
	return cfg
}

// btcNativeCustodyUnverifiable marks BTC|BTC rows where visible UTXOs are far below CSV.
func btcNativeCustodyUnverifiable(r walletzip.Row, actual decimal.Decimal) (bool, string) {
	if r.Coin != "BTC" || r.Network != "BTC" {
		return false, ""
	}
	if actual.IsZero() || r.Balance.IsZero() {
		return false, ""
	}
	if actual.Div(r.Balance).GreaterThan(decimal.NewFromFloat(0.01)) {
		return false, ""
	}
	return true, "on-chain UTXO sum << CSV; likely omnibus/internal BTC custody not visible at this address"
}

func ledgerSnapshotNote(spec onchainconfig.LedgerSpec) string {
	if !spec.LiveSnapshot {
		return ""
	}
	switch spec.Kind {
	case onchainconfig.LedgerSolanaNative, onchainconfig.LedgerSolanaSPL:
		return "live Solana balance vs snapshot CSV slot; public RPC has no historical slot state"
	case onchainconfig.LedgerNear, onchainconfig.LedgerNearFT:
		return "live NEAR balance vs snapshot block; public archival RPC unavailable for this height"
	case onchainconfig.LedgerSui, onchainconfig.LedgerSuiCoin:
		return "Sui public RPC may return latest checkpoint balance rather than snapshot checkpoint"
	case onchainconfig.LedgerTonJetton:
		return "TON jetton balance via tonapi (live); no public historical seqno API"
	case onchainconfig.LedgerTonNative:
		return "live TON balance via tonapi vs snapshot CSV; no public historical seqno API"
	case onchainconfig.LedgerFilecoin:
		return "live Filecoin StateGetActor balance vs snapshot CSV; Glif lookback limited on public nodes"
	case onchainconfig.LedgerChromia:
		return "Chromia postchain query; set CHROMIA_NODE if bootstrap node is unreachable"
	default:
		return "live chain balance vs snapshot CSV height"
	}
}

func ledgerUsedHistoricalSlot(components map[string]string) bool {
	return components != nil && components["historical_slot"] != ""
}

func ledgerLiveSnapshotNote(spec onchainconfig.LedgerSpec, components map[string]string) string {
	if ledgerUsedHistoricalSlot(components) {
		return ""
	}
	return ledgerSnapshotNote(spec)
}

func isLedgerCandidate(coin string) bool {
	switch coin {
	case "BTC", "DOGE", "LTC", "BCH", "BCHN", "ZEC", "XRP", "RIPPLE", "SOL",
		"USDC", "USDT", "USD1", "BOME", "WIF", "TRUMP", "WLFI", "FDUSD", "HFT", "RLUSD",
		"NEAR", "APT", "SUI", "DOT", "HBAR", "CHZ", "ENJ", "CHR", "CAKE", "FIL", "TRX",
		"TONCOIN-NEW", "TONCOIN":
		return true
	default:
		return false
	}
}
