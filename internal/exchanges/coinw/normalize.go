package coinw

import (
	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/exchanges/coinw/coinwapi"
	"github.com/ardmere/ardmere/internal/por"
)

const ID = "coinw"

// Meta enriches parsed CoinW summary with bundle metadata.
type Meta struct {
	Exchange  string
	PeriodSeq uint32
}

// Normalize maps CoinW public/imported summary into por.Snapshot.
// CoinW publishes monthly Merkle PoR + USDT total reserve ratio on the marketing page,
// but datacenter fetches commonly hit Cloudflare WAF; announcement text is the stable
// permissionless ratio channel assessed here. No permissionless wallet ZIP.
func Normalize(raw coinwapi.ParsedSummary, meta Meta) por.Snapshot {
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
			"totalReserveRateBasis": "announcement_or_page_total_usdt_ratio",
			"algorithm":             "merkle_tree+sha256 (user leaf via merkel_tree_cw.json + closed-source por-verify binary)",
			"auditor":               "self-attested (no named third-party in public announcement)",
			"reportDate":            raw.ReportDate,
			"amountsPublished":      "false",
			"walletListPublished":   "false",
			"cadence":               "monthly",
			"uiMerkle":              "Assets → Overview → Proof of Reserves (login) or offline por-verify + merkel_tree_cw.json",
			"ossVerifier":           "https://github.com/connectCoinw/proof-of-reserves (binaries only; no source to port)",
		},
	}
	if raw.SourceNote != "" {
		out.Extra["sourceNote"] = raw.SourceNote
	}
	if raw.WAFBlocked {
		out.Extra["wafBlocked"] = "true"
	}
	if raw.MerkleUnstable {
		out.Extra["merkleRootUnstable"] = "true"
	}
	if raw.MerkleRoot == "" {
		out.Extra["merkleRootPublished"] = "false"
	} else {
		out.Extra["merkleRootHexLen"] = decimal.NewFromInt(int64(len(raw.MerkleRoot))).String()
	}
	// Announcement-only fixtures expose a synthetic TOTAL row so solvency-claim has a subject.
	if raw.TotalReserveRate > 0 {
		out.CoinSummaries = append(out.CoinSummaries, por.CoinSummary{
			Coin: "TOTAL",
			Extra: map[string]string{
				"reserveRate":     decimal.NewFromFloat(raw.TotalReserveRate).String(),
				"reserveAmount":   "unpublished",
				"liabilityAmount": "unpublished",
				"basis":           "announcement_total_usdt_ratio",
			},
		})
	}
	return out
}
