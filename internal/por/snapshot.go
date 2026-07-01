// Package por defines exchange-neutral Proof-of-Reserves domain types.
// Exchange-specific fetchers normalize into these structs before verifiers run.
package por

import (
	"time"

	"github.com/shopspring/decimal"
)

// Snapshot is a normalized view of one exchange PoR audit period.
type Snapshot struct {
	Exchange        string
	ID              string
	SnapshotTime    time.Time
	SnapshotTimeRaw string
	MerkleRoot      string
	PeriodSeq       uint32
	TimeAnchor      *TimeAnchor
	CoinSummaries   []CoinSummary
	Extra           map[string]string // exchange-specific summary fields
	Auditor         string
	AuditorLink     string
	AuditDate       string
}

// TimeAnchor binds snapshot time to an external clock (e.g. BTC block height).
type TimeAnchor struct {
	Kind   string // "btc_block"
	Height uint32
}

// CoinSummary is per-coin reserve claims from the exchange summary artifact.
type CoinSummary struct {
	Coin              string
	CustomerLiability decimal.Decimal
	ExchangeReserve   decimal.Decimal // self-custody reserve (Binance: exchangeBalance)
	ThirdPartyReserve decimal.Decimal // third-party custody (Binance: thirdPartyCustody)
	ExchangeLiability decimal.Decimal // total exchange-side liability (Binance: binanceLiability)
	Extra             map[string]string
}
