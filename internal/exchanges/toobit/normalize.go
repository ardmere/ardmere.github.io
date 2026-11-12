package toobit

import (
	"strconv"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/exchanges/toobit/toobitapi"
	"github.com/ardmere/ardmere/internal/por"
)

const ID = "toobit"

// Meta enriches parsed Toobit summary with bundle metadata.
type Meta struct {
	Exchange  string
	PeriodSeq uint32
}

// Normalize maps Toobit public summary into por.Snapshot.
// Public summary publishes ratios only — amounts are unpublished.
func Normalize(raw toobitapi.ParsedSummary, meta Meta) por.Snapshot {
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
			"algorithm":             "SHA256-SMT",
			"hashSecurity":          "full-sha256+domain-sep",
			"auditor":               "self-attested (no named third-party in public summary)",
			"officialAuditId":       raw.OfficialAuditID,
			"treeHeight":            strconv.Itoa(raw.TreeHeight),
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
