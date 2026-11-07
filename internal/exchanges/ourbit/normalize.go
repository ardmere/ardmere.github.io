package ourbit

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/exchanges/ourbit/ourbitapi"
	"github.com/ardmere/ardmere/internal/por"
)

const ID = "ourbit"

// Meta enriches parsed Ourbit summary with bundle metadata.
type Meta struct {
	Exchange  string
	PeriodSeq uint32
}

// Normalize maps Ourbit public summary into por.Snapshot.
func Normalize(raw ourbitapi.ParsedSummary, meta Meta) por.Snapshot {
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
			"totalReserveRateBasis": "min_per_coin_stockRate",
			"algorithm":        "merkle_tree+sha256_truncated_128",
			"hashSecurity":     "128-bit-truncated-sha256",
			"auditor":          "self-attested (no named third-party in public summary)",
			"walletRows":       fmt.Sprintf("%d", len(raw.WalletRows)),
			"walletChannel":    "live_wallet_detail_api",
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
