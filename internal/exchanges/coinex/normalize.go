package coinex

import (
	"strings"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/exchanges/coinex/coinexapi"
	"github.com/ardmere/ardmere/internal/por"
)

const ID = "coinex"

type Meta struct {
	Exchange  string
	PeriodSeq uint32
}

// Normalize maps CoinEx public period summary into por.Snapshot.
func Normalize(raw coinexapi.ParsedSummary, meta Meta) por.Snapshot {
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
			"totalReserveRate":      decimal.NewFromFloat(raw.TotalReserveRate).String(),
			"totalReserveRateBasis": "min_per_coin_reserve_rate",
			"algorithm":             "merkle_tree+sha256",
			"hashSecurity":          "full-sha256",
			"auditor":               "self-attested (no named third-party in public summary)",
			"walletTotalUSD":        raw.WalletTotalUSD,
			"walletProofFile":       raw.WalletProofFile,
			"walletChannel":         "public_xlsx_address_list",
		},
	}
	anyAmounts := false
	for _, c := range raw.CoinRows {
		extra := map[string]string{
			por.ExtraReserveRate: decimal.NewFromFloat(c.ReserveRate).String(),
		}
		if c.WalletUSD != "" {
			extra["walletUSD"] = c.WalletUSD
		}
		if c.Price != "" {
			extra["price"] = c.Price
		}
		cs := por.CoinSummary{Coin: c.Coin}
		if c.ReserveAmount == 0 && c.LiabilityAmount == 0 {
			extra = por.MarkAmountsUnpublished(extra)
		} else {
			anyAmounts = true
			cs.ExchangeReserve = decimal.NewFromFloat(c.ReserveAmount)
			cs.CustomerLiability = decimal.NewFromFloat(c.LiabilityAmount)
		}
		cs.Extra = extra
		out.CoinSummaries = append(out.CoinSummaries, cs)
	}
	if !anyAmounts {
		out.Extra = por.SnapshotAmountsNotPublished(out.Extra)
	}
	return out
}

func normalizeRoot(root string) string {
	root = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(root)), "0x")
	return root
}
