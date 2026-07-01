package bitget

import (
	"strings"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/exchanges/bitget/bitgetapi"
	"github.com/ardmere/ardmere/internal/por"
)

const ID = "bitget"

type Meta struct {
	Exchange  string
	PeriodSeq uint32
}

func Normalize(raw bitgetapi.ParsedSummary, meta Meta) por.Snapshot {
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
		Extra: map[string]string{
			"totalReserveRate": decimal.NewFromFloat(raw.TotalReserveRate).String(),
			"algorithm":        "merkle_tree+sha256_truncated_64",
			"hashSecurity":     "weak-64-bit-truncated-sha256",
			"auditor":          "unidentified",
		},
	}
	for _, c := range raw.CoinRows {
		out.CoinSummaries = append(out.CoinSummaries, por.CoinSummary{
			Coin:              c.Coin,
			ExchangeReserve:   decimal.NewFromFloat(c.ReserveAmount),
			CustomerLiability: decimal.NewFromFloat(c.LiabilityAmount),
			Extra: map[string]string{
				"reserveRate": decimal.NewFromFloat(c.ReserveRate).String(),
			},
		})
	}
	return out
}

func normalizeRoot(root string) string {
	root = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(root)), "0x")
	if len(root) >= 64 {
		return root
	}
	if root == "" {
		return ""
	}
	return strings.Repeat("0", 64-len(root)) + root
}
