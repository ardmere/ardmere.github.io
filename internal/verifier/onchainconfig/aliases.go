package onchainconfig

import "strings"

type aliasTables struct {
	NetworkAliases map[string]string `json:"networkAliases,omitempty"`
	CoinAliases    map[string]string `json:"coinAliases,omitempty"`
}

func splitCoinNetwork(key string) (coin, network string) {
	i := strings.Index(key, "|")
	if i < 0 {
		return key, ""
	}
	return key[:i], key[i+1:]
}

// expandCoinNetworkAliases registers wallet CSV labels that map to canonical config keys.
// networkAliases: wallet network → canonical network (e.g. TRON → TRX).
// coinAliases: wallet coin → canonical coin (e.g. APTOS → APT).
func expandCoinNetworkAliases[V any](base map[string]V, aliases aliasTables) {
	if len(base) == 0 {
		return
	}
	walletCoinsByCanon := map[string][]string{}
	walletNetsByCanon := map[string][]string{}
	for wc, cc := range aliases.CoinAliases {
		walletCoinsByCanon[cc] = append(walletCoinsByCanon[cc], wc)
	}
	for wn, cn := range aliases.NetworkAliases {
		walletNetsByCanon[cn] = append(walletNetsByCanon[cn], wn)
	}
	for key, spec := range base {
		coin, net := splitCoinNetwork(key)
		coins := []string{coin}
		coins = append(coins, walletCoinsByCanon[coin]...)
		nets := []string{net}
		nets = append(nets, walletNetsByCanon[net]...)
		for _, wc := range uniqueStrings(coins) {
			for _, wn := range uniqueStrings(nets) {
				aliasKey := wc + "|" + wn
				if aliasKey == key {
					continue
				}
				if _, exists := base[aliasKey]; !exists {
					base[aliasKey] = spec
				}
			}
		}
	}
}

func uniqueStrings(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range in {
		if s == "" || seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return out
}
