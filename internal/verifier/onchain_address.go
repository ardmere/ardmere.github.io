package verifier

import (
	"fmt"
	"strings"

	"github.com/ardmere/ardmere/internal/rpc"
)

func nonQueryableAddressNote(coin, network, addr string) string {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return "empty wallet address in CSV"
	}
	lower := strings.ToLower(addr)
	if strings.Contains(lower, "bridging") {
		return "placeholder wallet address (bridging in progress) — not a queryable on-chain address"
	}
	switch network {
	case "STARKNET":
		if _, err := normalizeStarknetAddress(addr); err != nil {
			return fmt.Sprintf("invalid Starknet wallet address for %s@%s", coin, network)
		}
	default:
		if isEVMOnchainNetwork(network) {
			if _, err := rpc.EncodeHexAddress(addr); err != nil {
				return fmt.Sprintf("invalid EVM wallet address for %s@%s", coin, network)
			}
		}
	}
	return ""
}

func isEVMOnchainNetwork(network string) bool {
	switch network {
	case "ETH", "BSC", "ARBITRUM", "OPTIMISM", "BASE", "MATIC", "AVAXC", "OPBNB", "SONIC",
		"WLD", "CELO", "PLASMA", "KAIA", "KAVAEVM", "ZKSYNCERA", "RON", "SEIEVM", "MANTA",
		"LINEA", "SCROLL", "CHZ2", "MTL", "XLAYER", "FEVM":
		return true
	default:
		return false
	}
}

func normalizeStarknetAddress(addr string) (string, error) {
	addr = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(addr)), "0x")
	if addr == "" {
		return "", fmt.Errorf("empty starknet address")
	}
	if len(addr) > 64 {
		return "", fmt.Errorf("starknet address too long")
	}
	return "0x" + strings.Repeat("0", 64-len(addr)) + addr, nil
}
