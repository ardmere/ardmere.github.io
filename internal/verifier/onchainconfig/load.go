// Package onchainconfig loads per-exchange (coin|network) → RPC mapping tables.
package onchainconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/ardmere/ardmere/internal/rpc"
)

type TokenSpec struct {
	Net      rpc.Network
	Contract string
	Decimals int
	Native   bool
}

type tokenSpec = TokenSpec // alias for internal maps

// Config holds native and ERC20/BEP20 mappings for one exchange wallet CSV schema.
type Config struct {
	Native map[string]rpc.Network
	Tokens map[string]TokenSpec
}

type fileConfig struct {
	Native         map[string]string `json:"native"`
	Tokens         []tokenEntry      `json:"tokens"`
	NetworkAliases map[string]string `json:"networkAliases,omitempty"`
	CoinAliases    map[string]string `json:"coinAliases,omitempty"`
}

type tokenEntry struct {
	Coin       string `json:"coin"`
	Network    string `json:"network"`
	RPCNetwork string `json:"rpcNetwork"`
	Contract   string `json:"contract,omitempty"`
	Decimals   int    `json:"decimals"`
	Native     bool   `json:"native,omitempty"`
}

var (
	mu      sync.RWMutex
	byEx    = map[string]*Config{}
	configRoot = "config/exchanges"
)

// ResetForTest clears cached on-chain configs (tests only).
func ResetForTest() {
	mu.Lock()
	byEx = map[string]*Config{}
	mu.Unlock()
	ledgerMu.Lock()
	ledgerByEx = map[string]map[string]LedgerSpec{}
	ledgerMu.Unlock()
}

// Root overrides the default config/exchanges directory (for tests).
func Root(dir string) { configRoot = dir }

// ForExchange returns on-chain mappings for an exchange id (e.g. binance, gateio).
func ForExchange(exchangeID string) (*Config, error) {
	mu.RLock()
	if c, ok := byEx[exchangeID]; ok {
		mu.RUnlock()
		return c, nil
	}
	mu.RUnlock()

	mu.Lock()
	defer mu.Unlock()
	if c, ok := byEx[exchangeID]; ok {
		return c, nil
	}
	path, err := resolveConfigPath(exchangeID)
	if err != nil {
		return nil, err
	}
	c, err := loadFile(path)
	if err != nil {
		return nil, err
	}
	byEx[exchangeID] = c
	return c, nil
}

func resolveConfigPath(exchangeID string) (string, error) {
	if p := os.Getenv("ONCHAIN_CONFIG"); p != "" {
		return p, nil
	}
	rel := filepath.Join(configRoot, exchangeID, "onchain.json")
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
		filepath.Join("..", "..", "..", rel),
	} {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("onchain config not found for exchange %q (tried %s)", exchangeID, rel)
}

func findModuleRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for i := 0; i < 8; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

func loadFile(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read onchain config %s: %w", path, err)
	}
	var fc fileConfig
	if err := json.Unmarshal(raw, &fc); err != nil {
		return nil, fmt.Errorf("decode onchain config %s: %w", path, err)
	}
	out := &Config{
		Native: map[string]rpc.Network{},
		Tokens: map[string]tokenSpec{},
	}
	for k, v := range fc.Native {
		net, err := parseNetwork(v)
		if err != nil {
			return nil, fmt.Errorf("%s native %q: %w", path, k, err)
		}
		out.Native[k] = net
	}
	for _, t := range fc.Tokens {
		net, err := parseNetwork(t.RPCNetwork)
		if err != nil {
			return nil, fmt.Errorf("%s token %s|%s: %w", path, t.Coin, t.Network, err)
		}
		key := t.Coin + "|" + t.Network
		out.Tokens[key] = TokenSpec{
			Net:      net,
			Contract: t.Contract,
			Decimals: t.Decimals,
			Native:   t.Native,
		}
	}
	expandCoinNetworkAliases(out.Native, aliasTables{
		NetworkAliases: fc.NetworkAliases,
		CoinAliases:    fc.CoinAliases,
	})
	expandCoinNetworkAliases(out.Tokens, aliasTables{
		NetworkAliases: fc.NetworkAliases,
		CoinAliases:    fc.CoinAliases,
	})
	return out, nil
}

func parseNetwork(s string) (rpc.Network, error) {
	switch s {
	case "ETH":
		return rpc.NetEthereum, nil
	case "BSC":
		return rpc.NetBSC, nil
	case "ARBITRUM":
		return rpc.NetArbitrum, nil
	case "OPTIMISM":
		return rpc.NetOptimism, nil
	case "BASE":
		return rpc.NetBase, nil
	case "MATIC":
		return rpc.NetPolygon, nil
	case "AVAXC":
		return rpc.NetAvalanche, nil
	case "OPBNB":
		return rpc.NetOpBNB, nil
	case "SONIC":
		return rpc.NetSonic, nil
	case "WLD":
		return rpc.NetWorld, nil
	case "TRX":
		return rpc.NetTron, nil
	case "CELO":
		return rpc.NetCelo, nil
	case "PLASMA":
		return rpc.NetPlasma, nil
	case "KAIA":
		return rpc.NetKaia, nil
	case "KAVAEVM":
		return rpc.NetKavaEVM, nil
	case "ZKSYNCERA":
		return rpc.NetZkSync, nil
	case "RON":
		return rpc.NetRon, nil
	case "SEIEVM":
		return rpc.NetSeiEVM, nil
	case "MANTA":
		return rpc.NetManta, nil
	case "LINEA":
		return rpc.Network("LINEA"), nil
	case "SCROLL":
		return rpc.Network("SCROLL"), nil
	case "CHZ2":
		return rpc.Network("CHZ2"), nil
	case "MTL":
		return rpc.Network("MTL"), nil
	case "STARKNET":
		return rpc.Network("STARKNET"), nil
	case "AB":
		return rpc.Network("AB"), nil
	case "XLAYER":
		return rpc.NetXLayer, nil
	case "FEVM":
		return rpc.NetFEVM, nil
	default:
		return "", fmt.Errorf("unknown rpc network %q", s)
	}
}
