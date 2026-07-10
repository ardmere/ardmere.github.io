package onchainconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// LedgerKind identifies how to query a non-EVM HotCold row.
type LedgerKind string

const (
	LedgerEsplora      LedgerKind = "esplora"
	LedgerAlchemy      LedgerKind = "alchemy"
	LedgerBlockchair   LedgerKind = "blockchair"
	LedgerBlockcypher  LedgerKind = "blockcypher"
	LedgerSolanaNative LedgerKind = "solana"
	LedgerSolanaSPL    LedgerKind = "solana-spl"
	LedgerXRPL         LedgerKind = "xrpl"
	LedgerNear         LedgerKind = "near"
	LedgerNearFT       LedgerKind = "near-ft"
	LedgerAptos        LedgerKind = "aptos"
	LedgerAptosFA      LedgerKind = "aptos-fa"
	LedgerSui          LedgerKind = "sui"
	LedgerSuiCoin      LedgerKind = "sui-coin"
	LedgerHbar         LedgerKind = "hbar"
	LedgerHbarHTS      LedgerKind = "hbar-hts"
	LedgerSubstrate      LedgerKind = "substrate"
	LedgerSubstrateAsset LedgerKind = "substrate-asset"
	LedgerTonJetton      LedgerKind = "ton-jetton"
	LedgerAlgoASA        LedgerKind = "algo-asa"
	LedgerAptosCoin      LedgerKind = "aptos-coin"
	LedgerStellarAsset   LedgerKind = "stellar-asset"
	LedgerTezosFA        LedgerKind = "tezos-fa"
	LedgerXRPLEmitted    LedgerKind = "xrpl-issued"
	LedgerChromia        LedgerKind = "chromia"
	LedgerFilecoin       LedgerKind = "filecoin"
	LedgerTonNative      LedgerKind = "ton-native"
)

// LedgerSpec maps one HotCold (coin|network) to a ledger backend.
type LedgerSpec struct {
	Kind         LedgerKind
	Decimals     int
	EsploraBase  string // esplora
	Alchemy      string // alchemy utxo chain slug (bitcoin, bitcoin-cash)
	Blockchair   string // blockchair chain slug
	Blockcypher  string // blockcypher chain slug
	MaxTxCount   int    // blockcypher guard
	Mint         string // solana-spl / aptos-fa metadata / sui coin type / near-ft contract / ton jetton master
	Issuer       string // stellar issuer / xrpl issuer
	AssetID      int    // algo ASA / substrate Assets pallet id
	RPCURL       string // substrate / chromia node / optional override
	LiveSnapshot bool   // solana: compare live vs csv with WARN
}

type ledgerFile struct {
	Entries        []ledgerEntry     `json:"entries"`
	NetworkAliases map[string]string `json:"networkAliases,omitempty"`
	CoinAliases    map[string]string `json:"coinAliases,omitempty"`
}

type ledgerEntry struct {
	Coin         string `json:"coin"`
	Network      string `json:"network"`
	Backend      string `json:"backend"`
	Decimals     int    `json:"decimals"`
	EsploraBase  string `json:"esploraBase,omitempty"`
	Alchemy      string `json:"alchemy,omitempty"`
	Blockchair   string `json:"blockchair,omitempty"`
	Blockcypher  string `json:"blockcypher,omitempty"`
	MaxTxCount   int    `json:"maxTxCount,omitempty"`
	Mint         string `json:"mint,omitempty"`
	Issuer       string `json:"issuer,omitempty"`
	AssetID      int    `json:"assetId,omitempty"`
	RPCURL       string `json:"rpcUrl,omitempty"`
	LiveSnapshot bool   `json:"liveSnapshot,omitempty"`
}

var (
	ledgerMu      sync.RWMutex
	ledgerByEx    = map[string]map[string]LedgerSpec{}
	ledgerRoot    = "config/exchanges"
)

// LedgerRoot overrides config/exchanges for tests.
func LedgerRoot(dir string) { ledgerRoot = dir }

// LedgerForExchange returns ledger mappings for an exchange.
func LedgerForExchange(exchangeID string) (map[string]LedgerSpec, error) {
	ledgerMu.RLock()
	if m, ok := ledgerByEx[exchangeID]; ok {
		ledgerMu.RUnlock()
		return m, nil
	}
	ledgerMu.RUnlock()

	ledgerMu.Lock()
	defer ledgerMu.Unlock()
	if m, ok := ledgerByEx[exchangeID]; ok {
		return m, nil
	}
	path, err := resolveLedgerPath(exchangeID)
	if err != nil {
		return nil, err
	}
	m, err := loadLedgerFile(path)
	if err != nil {
		return nil, err
	}
	ledgerByEx[exchangeID] = m
	return m, nil
}

func resolveLedgerPath(exchangeID string) (string, error) {
	if p := os.Getenv("LEDGER_CONFIG"); p != "" {
		return p, nil
	}
	rel := filepath.Join(ledgerRoot, exchangeID, "ledger.json")
	if root := findModuleRoot(); root != "" {
		candidate := filepath.Join(root, rel)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	for _, candidate := range []string{
		rel,
		filepath.Join("..", rel),
		filepath.Join("..", "..", rel),
	} {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("ledger config not found for exchange %q (tried %s)", exchangeID, rel)
}

func loadLedgerFile(path string) (map[string]LedgerSpec, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read ledger config %s: %w", path, err)
	}
	var lf ledgerFile
	if err := json.Unmarshal(raw, &lf); err != nil {
		return nil, fmt.Errorf("decode ledger config %s: %w", path, err)
	}
	out := map[string]LedgerSpec{}
	for _, e := range lf.Entries {
		key := e.Coin + "|" + e.Network
		spec := LedgerSpec{
			Kind:         LedgerKind(e.Backend),
			Decimals:     e.Decimals,
			EsploraBase:  e.EsploraBase,
			Alchemy:      e.Alchemy,
			Blockchair:   e.Blockchair,
			Blockcypher:  e.Blockcypher,
			MaxTxCount:   e.MaxTxCount,
			Mint:         e.Mint,
			Issuer:       e.Issuer,
			AssetID:      e.AssetID,
			RPCURL:       e.RPCURL,
			LiveSnapshot: e.LiveSnapshot,
		}
		if spec.RPCURL == "" {
			spec.RPCURL = e.EsploraBase
		}
		switch spec.Kind {
		case LedgerSolanaNative, LedgerSolanaSPL, LedgerSui, LedgerSuiCoin, LedgerNear, LedgerNearFT, LedgerTonJetton, LedgerTonNative, LedgerFilecoin:
			if !spec.LiveSnapshot {
				spec.LiveSnapshot = true
			}
		}
		out[key] = spec
	}
	expandCoinNetworkAliases(out, aliasTables{
		NetworkAliases: lf.NetworkAliases,
		CoinAliases:    lf.CoinAliases,
	})
	return out, nil
}
