package rpc

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Provider describes one JSON-RPC endpoint and how we should use it.
type Provider struct {
	URL         string `json:"url"`
	Archive     bool   `json:"archive"`
	RateLimitMs int    `json:"rateLimitMs"`
	Weight      int    `json:"weight"`
	Note        string `json:"note,omitempty"`
}

type providerConfigFile map[string][]Provider

// DefaultProviderConfig returns the built-in fallback when no config file exists.
func DefaultProviderConfig() map[Network][]Provider {
	return map[Network][]Provider{
		NetEthereum: {
			{URL: "https://eth.drpc.org", Archive: true, Weight: 30, RateLimitMs: 100},
			{URL: "https://eth.llamarpc.com", Archive: true, Weight: 20},
			{URL: "https://ethereum-rpc.publicnode.com", Archive: false, Weight: 10, RateLimitMs: 200},
			{URL: "https://eth.merkle.io", Archive: true, Weight: 10},
		},
		NetBSC: {
			{URL: "https://bsc-mainnet.public.blastapi.io", Archive: true, Weight: 100, RateLimitMs: 500},
			{URL: "https://bsc.rpc.sentio.xyz", Archive: true, Weight: 90, RateLimitMs: 500},
			{URL: "https://bsc.drpc.org", Archive: true, Weight: 20, RateLimitMs: 1000},
		},
		NetArbitrum: {
			{URL: "https://arbitrum-one.public.blastapi.io", Archive: true, Weight: 100, RateLimitMs: 400},
			{URL: "https://arb1.arbitrum.io/rpc", Archive: true, Weight: 90, RateLimitMs: 300},
			{URL: "https://arbitrum.llamarpc.com", Archive: true, Weight: 80, RateLimitMs: 200},
			{URL: "https://arbitrum-one-rpc.publicnode.com", Archive: true, Weight: 70, RateLimitMs: 400},
			{URL: "https://arbitrum.drpc.org", Archive: true, Weight: 50, RateLimitMs: 600},
		},
		NetOptimism: {
			{URL: "https://optimism-mainnet.public.blastapi.io", Archive: true, Weight: 100, RateLimitMs: 400},
			{URL: "https://mainnet.optimism.io", Archive: true, Weight: 90, RateLimitMs: 300},
			{URL: "https://optimism.llamarpc.com", Archive: true, Weight: 80, RateLimitMs: 200},
			{URL: "https://optimism-rpc.publicnode.com", Archive: true, Weight: 70, RateLimitMs: 400},
			{URL: "https://optimism.drpc.org", Archive: true, Weight: 50, RateLimitMs: 600},
		},
		NetBase: {
			{URL: "https://base-mainnet.public.blastapi.io", Archive: true, Weight: 100, RateLimitMs: 400},
			{URL: "https://mainnet.base.org", Archive: true, Weight: 90, RateLimitMs: 300},
			{URL: "https://base.llamarpc.com", Archive: true, Weight: 80, RateLimitMs: 200},
			{URL: "https://base-rpc.publicnode.com", Archive: true, Weight: 70, RateLimitMs: 400},
			{URL: "https://base.drpc.org", Archive: true, Weight: 50, RateLimitMs: 600},
		},
		NetPolygon: {
			{URL: "https://polygon-mainnet.public.blastapi.io", Archive: true, Weight: 100, RateLimitMs: 400},
			{URL: "https://polygon-rpc.com", Archive: true, Weight: 90, RateLimitMs: 300},
			{URL: "https://polygon.llamarpc.com", Archive: true, Weight: 80, RateLimitMs: 200},
			{URL: "https://polygon-bor-rpc.publicnode.com", Archive: true, Weight: 70, RateLimitMs: 400},
			{URL: "https://polygon.drpc.org", Archive: true, Weight: 50, RateLimitMs: 600},
		},
		NetAvalanche: {
			{URL: "https://api.avax.network/ext/bc/C/rpc", Archive: true, Weight: 100, RateLimitMs: 300},
			{URL: "https://avalanche-c-chain-rpc.publicnode.com", Archive: true, Weight: 80, RateLimitMs: 300},
			{URL: "https://avalanche.drpc.org", Archive: true, Weight: 50, RateLimitMs: 500},
		},
		NetOpBNB: {
			{URL: "https://opbnb-mainnet-rpc.bnbchain.org", Archive: true, Weight: 100, RateLimitMs: 300},
			{URL: "https://opbnb.public.blastapi.io", Archive: true, Weight: 80, RateLimitMs: 300},
			{URL: "https://opbnb.drpc.org", Archive: true, Weight: 50, RateLimitMs: 500},
		},
		NetSonic: {
			{URL: "https://rpc.soniclabs.com", Archive: true, Weight: 100, RateLimitMs: 200},
			{URL: "https://sonic.drpc.org", Archive: true, Weight: 60, RateLimitMs: 400},
		},
		NetWorld: {
			{URL: "https://worldchain-mainnet.g.alchemy.com/public", Archive: true, Weight: 100, RateLimitMs: 200},
			{URL: "https://worldchain-mainnet.gateway.tenderly.co", Archive: true, Weight: 80, RateLimitMs: 300},
		},
		NetCelo: {
			{URL: "https://forno.celo.org", Archive: true, Weight: 100, RateLimitMs: 300},
			{URL: "https://celo.drpc.org", Archive: true, Weight: 60, RateLimitMs: 500},
		},
		NetPlasma: {
			{URL: "https://rpc.plasma.to", Archive: true, Weight: 100, RateLimitMs: 300},
			{URL: "https://plasma.drpc.org", Archive: true, Weight: 60, RateLimitMs: 500},
		},
		NetKaia: {
			{URL: "https://public-en.node.kaia.io", Archive: true, Weight: 100, RateLimitMs: 300},
			{URL: "https://kaia.drpc.org", Archive: true, Weight: 60, RateLimitMs: 500},
		},
		NetKavaEVM: {
			{URL: "https://evm.kava.io", Archive: true, Weight: 100, RateLimitMs: 300},
			{URL: "https://kava-evm.drpc.org", Archive: true, Weight: 60, RateLimitMs: 500},
		},
		NetZkSync: {
			{URL: "https://mainnet.era.zksync.io", Archive: true, Weight: 100, RateLimitMs: 300},
			{URL: "https://zksync.drpc.org", Archive: true, Weight: 60, RateLimitMs: 500},
		},
		NetRon: {
			{URL: "https://api.roninchain.com/rpc", Archive: true, Weight: 100, RateLimitMs: 300},
			{URL: "https://ronin.drpc.org", Archive: true, Weight: 60, RateLimitMs: 500},
		},
		NetSeiEVM: {
			{URL: "https://evm-rpc.sei-apis.com", Archive: true, Weight: 100, RateLimitMs: 300},
			{URL: "https://sei.drpc.org", Archive: true, Weight: 60, RateLimitMs: 500},
		},
		NetManta: {
			{URL: "https://pacific-rpc.manta.network/http", Archive: true, Weight: 100, RateLimitMs: 300},
			{URL: "https://manta-pacific.drpc.org", Archive: true, Weight: 60, RateLimitMs: 500},
		},
		NetXLayer: {
			{URL: "https://rpc.xlayer.tech", Archive: true, Weight: 100, RateLimitMs: 300},
			{URL: "https://xlayerrpc.okx.com", Archive: true, Weight: 90, RateLimitMs: 300},
		},
		NetFEVM: {
			{URL: "https://api.node.glif.io/rpc/v1", Archive: false, Weight: 100, RateLimitMs: 400},
			{URL: "https://rpc.ankr.com/filecoin", Archive: false, Weight: 80, RateLimitMs: 500},
		},
		NetTron: {
			{URL: "https://api.trongrid.io", Archive: true, Weight: 100, RateLimitMs: 250},
			{URL: "https://tron-rpc.publicnode.com", Archive: true, Weight: 80, RateLimitMs: 300},
		},
		NetSolana: {
			{URL: "https://solana.drpc.org", Archive: false, Weight: 100, RateLimitMs: 150},
			{URL: "https://solana-rpc.publicnode.com", Archive: false, Weight: 90, RateLimitMs: 200},
			{URL: "https://solana-mainnet.g.alchemy.com/v2/demo", Archive: false, Weight: 70, RateLimitMs: 250},
			{URL: "https://api.mainnet-beta.solana.com", Archive: false, Weight: 30, RateLimitMs: 400},
		},
	}
}

// LoadProviderConfig reads config/rpc-providers.json (or RPC_PROVIDERS_CONFIG).
func LoadProviderConfig() (map[Network][]Provider, error) {
	path := os.Getenv("RPC_PROVIDERS_CONFIG")
	if path == "" {
		for _, candidate := range []string{
			"config/rpc-providers.json",
			filepath.Join("..", "config", "rpc-providers.json"),
		} {
			if _, err := os.Stat(candidate); err == nil {
				path = candidate
				break
			}
		}
	}
	if path == "" {
		return DefaultProviderConfig(), nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read rpc providers %s: %w", path, err)
	}
	var file providerConfigFile
	if err := json.Unmarshal(raw, &file); err != nil {
		return nil, fmt.Errorf("parse rpc providers %s: %w", path, err)
	}
	out := map[Network][]Provider{}
	for netKey, providers := range file {
		net := Network(strings.ToUpper(netKey))
		out[net] = append(out[net], providers...)
		sortProviders(out[net])
	}
	return out, nil
}

func sortProviders(providers []Provider) {
	sort.Slice(providers, func(i, j int) bool {
		if providers[i].Weight != providers[j].Weight {
			return providers[i].Weight > providers[j].Weight
		}
		return providers[i].URL < providers[j].URL
	})
}

// ProvidersFor returns providers for a network, optionally preferring archive nodes
// when querying a historical block height.
func ProvidersFor(all map[Network][]Provider, net Network, historical bool) []Provider {
	src, ok := all[net]
	if !ok || len(src) == 0 {
		return nil
	}
	if !historical {
		return append([]Provider(nil), src...)
	}
	var archive, other []Provider
	for _, p := range src {
		if p.Archive {
			archive = append(archive, p)
		} else {
			other = append(other, p)
		}
	}
	if len(archive) > 0 {
		return append(archive, other...)
	}
	return append([]Provider(nil), src...)
}
