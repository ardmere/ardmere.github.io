package coingecko

import "strings"

// SymbolToID maps common exchange reserve symbols to CoinGecko coin ids.
// Unmapped symbols are excluded from both coverage numerator and denominator.
var SymbolToID = map[string]string{
	"BTC":      "bitcoin",
	"WBTC":     "wrapped-bitcoin",
	"ETH":      "ethereum",
	"WETH":     "weth",
	"USDT":     "tether",
	"USDC":     "usd-coin",
	"DAI":      "dai",
	"BNB":      "binancecoin",
	"SOL":      "solana",
	"XRP":      "ripple",
	"DOGE":     "dogecoin",
	"ADA":      "cardano",
	"TRX":      "tron",
	"DOT":      "polkadot",
	"MATIC":    "matic-network",
	"POL":      "polygon-ecosystem-token",
	"AVAX":     "avalanche-2",
	"LINK":     "chainlink",
	"UNI":      "uniswap",
	"ATOM":     "cosmos",
	"LTC":      "litecoin",
	"BCH":      "bitcoin-cash",
	"APT":      "aptos",
	"SUI":      "sui",
	"TON":      "the-open-network",
	"SHIB":     "shiba-inu",
	"PEPE":     "pepe",
	"ARB":      "arbitrum",
	"OP":       "optimism",
	"FIL":      "filecoin",
	"NEAR":     "near",
	"ICP":      "internet-computer",
	"ETC":      "ethereum-classic",
	"XLM":      "stellar",
	"HBAR":     "hedera-hashgraph",
	"VET":      "vechain",
	"AAVE":     "aave",
	"MKR":      "maker",
	"CRV":      "curve-dao-token",
	"LDO":      "lido-dao",
	"STETH":    "staked-ether",
	"WSTETH":   "wrapped-steth",
	"CBETH":    "coinbase-wrapped-staked-eth",
	"RETH":     "rocket-pool-eth",
	"TUSD":     "true-usd",
	"FDUSD":    "first-digital-usd",
	"USDE":     "ethena-usde",
	"PYUSD":    "paypal-usd",
	"BUSD":     "binance-usd",
	"OKB":      "okb",
	"OKT":      "oec-token",
}

// ResolveID returns the CoinGecko id for an exchange asset symbol, or "" if unmapped.
func ResolveID(symbol string) string {
	s := strings.ToUpper(strings.TrimSpace(symbol))
	if id, ok := SymbolToID[s]; ok {
		return id
	}
	return ""
}
