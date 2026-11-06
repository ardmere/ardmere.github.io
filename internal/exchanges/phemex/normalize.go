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
	for _, c := range raw.CoinRows {
		extra := map[string]string{
			"reserveRate":  decimal.NewFromFloat(c.ReserveRate).String(),
			"chainBalance": formatFloat(c.ReserveAmount),
			"cexBalance":   formatFloat(c.LiabilityAmount),
		}
		out.CoinSummaries = append(out.CoinSummaries, por.CoinSummary{
			Coin:              c.Coin,
			ExchangeReserve:   decimal.NewFromFloat(c.ReserveAmount),
			CustomerLiability: decimal.NewFromFloat(c.LiabilityAmount),
			Extra:             extra,
		})
	}
	return out
}

func formatFloat(v float64) string {
	return fmt.Sprintf("%g", v)
}
