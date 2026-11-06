package kraken

import (
	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/exchanges/kraken/krakenapi"
	"github.com/ardmere/ardmere/internal/por"
)

const ID = "kraken"

// Meta enriches parsed Kraken summary with bundle metadata.
type Meta struct {
	Exchange  string
	PeriodSeq uint32
}

// Normalize maps Kraken public summary into por.Snapshot.
// Public disclosure: CMS reserve ratios (+ capped "105%+" stables) and LedgerLens accountant Merkle root.
// No permissionless wallet amounts/addresses.
func Normalize(raw krakenapi.ParsedSummary, meta Meta) por.Snapshot {
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
			"totalReserveRateBasis": "min_per_coin_effective_reserve_rate",
			"algorithm":             "accountant Merkle tree (The Network Firm LedgerLens; user leaf login-gated)",
			"auditor":               raw.Accountant,
			"reportDate":            raw.ReportDate,
			"auditPeriodId":         raw.AuditPeriodID,
			"ledgerlensDatasetId":   raw.DatasetID,
			"merkleRootHexLen":      decimal.NewFromInt(int64(len(raw.MerkleRoot))).String(),
			"amountsPublished":      "false",
			"walletListPublished":   "false",
			"cadence":               "quarterly accountant attestation",
			"uiMerkle":              "Kraken Pro login → Proof of Reserves → Verify yourself / Merkle Leaf ID in accountant tool",
			"attestationModel":      "true",
		},
	}
	for _, c := range raw.CoinRows {
		extra := map[string]string{
			"reserveRate":       decimal.NewFromFloat(c.ReserveRate).String(),
			"reserveRateDisplay": c.DisplayRatio,
			"reserveAmount":     "unpublished",
			"liabilityAmount":   "unpublished",
		}
		if c.CapPercentage {
			extra["reserveRateCapped"] = "true"
			extra["reserveRateFloor"] = "105"
		}
		if c.HideBalances {
			extra["balancesHidden"] = "true"
		}
		out.CoinSummaries = append(out.CoinSummaries, por.CoinSummary{
			Coin:  c.Coin,
			Extra: extra,
		})
	}
	return out
}
