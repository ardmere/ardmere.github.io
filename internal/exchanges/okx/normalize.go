package okx

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/exchanges/okx/okxapi"
	"github.com/ardmere/ardmere/internal/por"
)

const ID = "okx"

// Meta enriches parsed OKX summary with bundle metadata.
type Meta struct {
	Exchange  string
	PeriodSeq uint32
}

// Normalize maps OKX auditRootInfo into por.Snapshot.
func Normalize(raw okxapi.AuditRootInfo, meta Meta) por.Snapshot {
	if meta.Exchange == "" {
		meta.Exchange = ID
	}
	snapTime := time.UnixMilli(raw.CreateTime).UTC()

	out := por.Snapshot{
		Exchange:        meta.Exchange,
		ID:              raw.AuditID,
		SnapshotTime:    snapTime,
		SnapshotTimeRaw: snapTime.Format(time.RFC3339),
		MerkleRoot:      raw.MerkleHash,
		PeriodSeq:       meta.PeriodSeq,
		Extra: map[string]string{
			"typeNum":     itoa(raw.TypeNum),
			"algorithm":   "merkle_tree+zk_stark",
			"reserveCoin": joinCoins(raw.ReserveCurrencies),
		},
	}

	for _, coin := range raw.ReserveCurrencies {
		ratio := raw.CapitalRatio[coin]
		out.CoinSummaries = append(out.CoinSummaries, por.CoinSummary{
			Coin:              coin,
			CustomerLiability: parseDec(raw.LiabilityBalances[coin]),
			ExchangeReserve:   parseDec(raw.ExchangeReserveBalances[coin]),
			ThirdPartyReserve: parseDec(raw.CustodyReserveBalances[coin]),
			ExchangeLiability: parseDec(raw.ReserveBalances[coin]),
			Extra: map[string]string{
				"capitalRatio": ratio,
			},
		})
	}
	return out
}

func parseDec(s string) decimal.Decimal {
	if s == "" {
		return decimal.Zero
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero
	}
	return d
}

func itoa(v int) string {
	return decimal.NewFromInt(int64(v)).String()
}

func joinCoins(coins []string) string {
	if len(coins) == 0 {
		return ""
	}
	out := coins[0]
	for i := 1; i < len(coins); i++ {
		out += "," + coins[i]
	}
	return out
}
