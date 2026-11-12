package bitunix

import (
	"strings"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/exchanges/bitunix/bitunixapi"
	"github.com/ardmere/ardmere/internal/por"
)

const ID = "bitunix"

// Meta enriches parsed Bitunix summary with bundle metadata.
type Meta struct {
	Exchange  string
	PeriodSeq uint32
}

// Normalize maps Bitunix public summary into por.Snapshot.
// Public summary publishes ratios only — amounts are unpublished.
func Normalize(raw bitunixapi.ParsedSummary, meta Meta) por.Snapshot {
	if meta.Exchange == "" {
		meta.Exchange = ID
	}
	out := por.Snapshot{
		Exchange:        meta.Exchange,
		ID:              raw.AuditID,
		SnapshotTime:    raw.AuditTime,
		SnapshotTimeRaw: raw.AuditTimeRaw,
		MerkleRoot:      normalizeRoot(raw.MerkleRoot),
		PeriodSeq:       meta.PeriodSeq,
		Extra: por.SnapshotAmountsNotPublished(map[string]string{
			"totalReserveRate":      decimal.NewFromFloat(raw.TotalReserveRate).String(),
			"totalReserveRateBasis": "min_per_coin_reserve_rate",
			"algorithm":             "merkle_tree+sha256",
			"hashSecurity":          "full-sha256",
			"auditor":               "self-attested (no named third-party in public summary)",
			"walletAddress":         raw.WalletAddress,
		}),
	}
	for _, c := range raw.CoinRows {
		out.CoinSummaries = append(out.CoinSummaries, por.CoinSummary{
			Coin: c.Coin,
			Extra: por.MarkAmountsUnpublished(map[string]string{
				por.ExtraReserveRate: decimal.NewFromFloat(c.ReserveRate).String(),
			}),
		})
	}
	return out
}

func normalizeRoot(root string) string {
	return strings.TrimPrefix(strings.ToLower(strings.TrimSpace(root)), "0x")
}
