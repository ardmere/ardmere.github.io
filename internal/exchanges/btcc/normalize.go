package btcc

import (
	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/exchanges/btcc/btccapi"
	"github.com/ardmere/ardmere/internal/por"
)

const ID = "btcc"

// Meta enriches parsed BTCC summary with bundle metadata.
type Meta struct {
	Exchange  string
	PeriodSeq uint32
}

// Normalize maps BTCC public summary into por.Snapshot.
// BTCC publishes ratios + truncated merkle root only — no wallet/liability amounts.
func Normalize(raw btccapi.ParsedSummary, meta Meta) por.Snapshot {
	if meta.Exchange == "" {
		meta.Exchange = ID
	}
	out := por.Snapshot{
		Exchange:        meta.Exchange,
		ID:              raw.AuditID,
		SnapshotTime:    raw.AuditTime,
		SnapshotTimeRaw: raw.AuditTimeRaw,
		MerkleRoot:      raw.MerkleRoot,
		PeriodSeq:       meta.PeriodSeq,
		Extra: map[string]string{
			"totalReserveRate":      decimal.NewFromFloat(raw.TotalReserveRate).String(),
			"totalReserveRateBasis": "min_per_coin_reserve_rate",
			"algorithm":             "merkle_tree (login leaf self-check; public root truncated to ~16 hex)",
			"auditor":               "self-attested (no named third-party in public summary)",
			"reportDate":            raw.ReportDate,
			"merkleRootTruncated":   "true",
			"merkleRootHexLen":      decimal.NewFromInt(int64(len(raw.MerkleRoot))).String(),
			"amountsPublished":      "false",
			"walletListPublished":   "false",
			"cadence":               "monthly ~day-15 (UTC createDate)",
			"uiMerkle":              "paste leaf hash after login (View My Report); validate API requires token",
		},
	}
	for _, c := range raw.CoinRows {
		out.CoinSummaries = append(out.CoinSummaries, por.CoinSummary{
			Coin: c.Coin,
			Extra: map[string]string{
				"reserveRate": decimal.NewFromFloat(c.ReserveRate).String(),
				// Explicit: BTCC does not publish reserve/liability amounts.
				"reserveAmount":   "unpublished",
				"liabilityAmount": "unpublished",
			},
		})
	}
	return out
}
