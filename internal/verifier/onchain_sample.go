package verifier

import (
	"fmt"

	"github.com/ardmere/ardmere/internal/walletzip"
)

const (
	okxOnchainTopKPerPair = 40
	okxOnchainMaxTotal    = 800
)

func maybeSubsampleOnchainRows(exchange string, supported []walletzip.Row) ([]walletzip.Row, string) {
	if exchange != "okx" || len(supported) <= okxOnchainMaxTotal {
		return supported, ""
	}
	sampled := walletzip.SampleTopRowsByPair(supported, okxOnchainTopKPerPair, okxOnchainMaxTotal)
	cov := walletzip.ValueCoverageOf(supported, sampled)
	note := fmt.Sprintf("OKX wallet subsample: %d/%d supported rows (%.1f%% balance coverage)",
		len(sampled), len(supported), cov*100)
	return sampled, note
}
