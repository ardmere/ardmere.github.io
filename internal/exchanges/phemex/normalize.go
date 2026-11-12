package phemex

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/exchanges/phemex/phemexapi"
	"github.com/ardmere/ardmere/internal/por"
)

const ID = "phemex"

// Meta enriches parsed Phemex summary with bundle metadata.
type Meta struct {
	Exchange  string
	PeriodSeq uint32
}

// Normalize maps Phemex public summary into por.Snapshot.
func Normalize(raw phemexapi.ParsedSummary, meta Meta) por.Snapshot {
	if meta.Exchange == "" {
		meta.Exchange = ID
	}
	out := por.Snapshot{
		Exchange:        meta.Exchange,
		ID:              raw.AuditID,
		SnapshotTime:    raw.AuditTime,
		SnapshotTimeRaw: raw.AuditTimeRaw,
		PeriodSeq:       meta.PeriodSeq,
		Extra: map[string]string{
			"totalReserveRate":      decimal.NewFromFloat(raw.TotalReserveRate).String(),
			"totalReserveRateBasis": "min_per_coin_reserve_rate",
			"algorithm":             "merkle_tree (UI / login user inclusion; root not in public porVersions)",
			"auditor":               "self-attested (no named third-party in public summary)",
			"officialVersionId":     raw.OfficialVersionID,
			"dataVersion":           raw.DataVersion,
			"uiWalletList":          "public PoR page section Part of Our Cold Wallets (no download; api10 gated)",
			"uiMerkle":              "monthly Merkle inclusion via logged-in hashed client ID tool",
		},
	}
	for _, h := range raw.BlockHeights {
		if strings.EqualFold(h.Chain, "Bitcoin") {
			out.Extra["btcBlockHeight"] = h.Height
			break
		}
	}
	anyAmounts := false
	for _, c := range raw.CoinRows {
		extra := map[string]string{
			por.ExtraReserveRate: decimal.NewFromFloat(c.ReserveRate).String(),
		}
		cs := por.CoinSummary{Coin: c.Coin}
		if c.ReserveAmount == 0 && c.LiabilityAmount == 0 {
			extra = por.MarkAmountsUnpublished(extra)
		} else {
			anyAmounts = true
			cs.ExchangeReserve = decimal.NewFromFloat(c.ReserveAmount)
			cs.CustomerLiability = decimal.NewFromFloat(c.LiabilityAmount)
			extra["chainBalance"] = formatFloat(c.ReserveAmount)
			extra["cexBalance"] = formatFloat(c.LiabilityAmount)
		}
		cs.Extra = extra
		out.CoinSummaries = append(out.CoinSummaries, cs)
	}
	if !anyAmounts {
		out.Extra = por.SnapshotAmountsNotPublished(out.Extra)
	}
	return out
}

func formatFloat(v float64) string {
	return fmt.Sprintf("%g", v)
}
