package woox

import (
	"fmt"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/exchanges/woox/wooxapi"
	"github.com/ardmere/ardmere/internal/por"
)

const ID = "woox"

// Meta enriches parsed WOO X summary with bundle metadata.
type Meta struct {
	Exchange  string
	PeriodSeq uint32
}

// Normalize maps WOO X public summary into por.Snapshot.
func Normalize(raw wooxapi.ParsedSummary, meta Meta) por.Snapshot {
	if meta.Exchange == "" {
		meta.Exchange = ID
	}
	out := por.Snapshot{
		Exchange:        meta.Exchange,
		ID:              raw.AuditID,
		SnapshotTime:    raw.AuditTime,
		SnapshotTimeRaw: raw.AuditTimeRaw,
		// Period merkle root hex is not published separately; full protobuf tree is public.
		MerkleRoot: "",
		PeriodSeq:  meta.PeriodSeq,
		Extra: map[string]string{
			"totalReserveRate":      decimal.NewFromFloat(raw.TotalReserveRate).String(),
			"totalReserveRateBasis": "sum_asset_usdt_notional / sum_liability_usdt_notional * 100",
			"totalAssetsUSDT":       decimal.NewFromFloat(raw.TotalAssetsUSDT).String(),
			"totalLiabilitiesUSDT":  decimal.NewFromFloat(raw.TotalLiabilitiesUSDT).String(),
			"officialAuditId":       raw.OfficialAuditID,
			"merkleTreeURL":         raw.MerkleTreeURL,
			"merkleTreeBytes":       fmt.Sprintf("%d", raw.MerkleTreeBytes),
			"merkleTreeETag":        raw.MerkleTreeETag,
			"assetTokenCount":       fmt.Sprintf("%d", raw.AssetTokenCount),
			"liabilityTokenCount":   fmt.Sprintf("%d", raw.LiabilityTokenCount),
			"algorithm":             "merkle_tree+sha256+protobuf (woo-x/merkle-tree)",
			"hashSecurity":          "full-sha256 (leaf + parent; official Java OSS)",
			"auditor":               "self-attested live dashboard + monthly merkle audit (no named third-party in public summary)",
			"liveReserveRate":       decimal.NewFromFloat(raw.LiveReserveRate).String(),
			"assetVenues":           raw.VenuesNote,
		},
	}
	return out
}
