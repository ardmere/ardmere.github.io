// Package phemexapi fetches and parses Phemex Proof-of-Reserves summary artifacts.
// See docs/phemex-por-data-guide.md for data availability.
package phemexapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	userAgent = "ardmere/0.1 (+https://ardmere.org)"

	URLPoRPage     = "https://phemex.com/proof-of-reserves"
	URLPorVersions = "https://por.phemex.com/porVersions"
	// Browser-context APIs (CloudFront / auth gated from many datacenters):
	URLReserveFunds    = "https://api10.phemex.com/phemex-user/public/v2/por/reserve-funds"
	URLMerkleRoot      = "https://api10.phemex.com/phemex-user/public/v2/por/merkle-tree/root"
	URLWalletAddresses = "https://api10.phemex.com/phemex-user/public/v2/por/wallet-addresses"
)

// SummaryBundle is the archived summarySnapshot for one Phemex PoR period.
type SummaryBundle struct {
	FetchedAt         time.Time       `json:"fetchedAt"`
	Source            string          `json:"source,omitempty"`
	SourceURLs        []string        `json:"sourceUrls,omitempty"`
	OfficialVersionID string          `json:"officialVersionId,omitempty"`
	SnapshotID        string          `json:"snapshotId,omitempty"`
	Info              json.RawMessage `json:"info"`
}

// ParsedSummary is normalized fields extracted from porVersions entry.
// Reserve rates are percentages (100 = 100%) for solvency-claim.
type ParsedSummary struct {
	AuditID           string
	OfficialVersionID string
	AuditTime         time.Time
	AuditTimeRaw      string
	DataVersion       string
	TotalReserveRate  float64
	CoinRows          []CoinRow
	BlockHeights      []BlockHeight
}

type CoinRow struct {
	Coin            string
	ReserveRate     float64
	ReserveAmount   float64 // chainBalance
	LiabilityAmount float64 // cexBalance
}

type BlockHeight struct {
	Chain  string
	Height string
}

type versionJSON struct {
	VersionID          string        `json:"versionID"`
	VersionDescription string        `json:"versionDescription"`
	DataVersion        string        `json:"dataVersion"`
	BlockHeights       []blockJSON   `json:"blockHeights"`
	TotalBalances      []balanceJSON `json:"totalBalances"`
}

type blockJSON struct {
	Chain  string `json:"chain"`
	Height string `json:"height"`
}

type balanceJSON struct {
	Symbol       string `json:"symbol"`
	ChainBalance string `json:"chainBalance"`
	CexBalance   string `json:"cexBalance"`
}

// SummaryURL returns the public summary API used for artifact provenance.
func SummaryURL() string { return URLPorVersions }

// FetchPublicSummary downloads porVersions and archives the latest entry with balances.
func FetchPublicSummary(ctx context.Context) ([]byte, ParsedSummary, error) {
	return FetchPublicSummaryByVersionID(ctx, "")
}

// FetchPublicSummaryByVersionID selects a versionID (e.g. "2026-11-01") from porVersions.
func FetchPublicSummaryByVersionID(ctx context.Context, versionID string) ([]byte, ParsedSummary, error) {
	raw, err := getJSON(ctx, URLPorVersions)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	var versions []json.RawMessage
	if err := json.Unmarshal(raw, &versions); err != nil {
		return nil, ParsedSummary{}, fmt.Errorf("decode porVersions list: %w", err)
	}
	if len(versions) == 0 {
		return nil, ParsedSummary{}, fmt.Errorf("porVersions empty")
	}

	versionID = strings.TrimSpace(versionID)
	var chosen json.RawMessage
	if versionID == "" {
		chosen = versions[0]
	} else {
		want := normalizeVersionKey(versionID)
		for _, v := range versions {
			var meta versionJSON
			if err := json.Unmarshal(v, &meta); err != nil {
				continue
			}
			if normalizeVersionKey(meta.VersionID) == want || normalizeVersionKey(meta.VersionDescription) == want {
				chosen = v
				break
			}
			// Also accept compact snapshot ids like 20261101.
			if compactSnapshotID(meta.VersionID) == want {
				chosen = v
				break
			}
		}
		if len(chosen) == 0 {
			return nil, ParsedSummary{}, fmt.Errorf("porVersions: version %q not found", versionID)
		}
	}
	return buildBundle(chosen, URLPorVersions)
}

// ParseSummaryBytes parses an archived SummaryBundle or a raw porVersions entry.
func ParseSummaryBytes(raw []byte) (ParsedSummary, error) {
	var bundle SummaryBundle
	if err := json.Unmarshal(raw, &bundle); err == nil && len(bundle.Info) > 0 {
		return parseInfo(bundle.Info)
	}
	var top map[string]json.RawMessage
	if err := json.Unmarshal(raw, &top); err == nil {
		if info, ok := top["info"]; ok {
			return parseInfo(info)
		}
	}
	return parseInfo(json.RawMessage(raw))
}

// BuildSummaryBundle wraps a browser-captured porVersions entry (or full list) as an archive bundle.
func BuildSummaryBundle(infoFile string) ([]byte, error) {
	raw, err := os.ReadFile(infoFile)
	if err != nil {
		return nil, err
	}
	var list []json.RawMessage
	if err := json.Unmarshal(raw, &list); err == nil && len(list) > 0 {
		raw = list[0]
	}
	out, _, err := buildBundle(json.RawMessage(raw), URLPorVersions)
	return out, err
}

func buildBundle(infoRaw json.RawMessage, source string) ([]byte, ParsedSummary, error) {
	parsed, err := parseInfo(infoRaw)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	bundle := SummaryBundle{
		FetchedAt:         time.Now().UTC(),
		Source:            source,
		SourceURLs:        []string{source, URLPoRPage},
		OfficialVersionID: parsed.OfficialVersionID,
		SnapshotID:        parsed.AuditID,
		Info:              infoRaw,
	}
	out, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	return out, parsed, nil
}

func parseInfo(infoRaw json.RawMessage) (ParsedSummary, error) {
	var info versionJSON
	if err := json.Unmarshal(infoRaw, &info); err != nil {
		return ParsedSummary{}, fmt.Errorf("decode phemex por version: %w", err)
	}
	official := strings.TrimSpace(info.VersionID)
	if official == "" {
		return ParsedSummary{}, fmt.Errorf("phemex summary missing versionID")
	}
	if len(info.TotalBalances) == 0 {
		return ParsedSummary{}, fmt.Errorf("phemex version %s has no totalBalances (older dataVersion=%s)", official, info.DataVersion)
	}

	ts, tsRaw := parseVersionTime(info.VersionDescription, official)
	auditID := compactSnapshotID(official)
	if auditID == "" && !ts.IsZero() {
		auditID = ts.Format("20060102")
	}

	rows := make([]CoinRow, 0, len(info.TotalBalances))
	minRate := math.MaxFloat64
	for _, b := range info.TotalBalances {
		coin := strings.TrimSpace(b.Symbol)
		if coin == "" {
			continue
		}
		reserve, err := strconv.ParseFloat(strings.TrimSpace(b.ChainBalance), 64)
		if err != nil {
			return ParsedSummary{}, fmt.Errorf("parse chainBalance for %s: %w", coin, err)
		}
		liability, err := strconv.ParseFloat(strings.TrimSpace(b.CexBalance), 64)
		if err != nil {
			return ParsedSummary{}, fmt.Errorf("parse cexBalance for %s: %w", coin, err)
		}
		rate := 0.0
		if liability > 0 {
			rate = reserve / liability * 100.0
		}
		rows = append(rows, CoinRow{
			Coin:            coin,
			ReserveRate:     rate,
			ReserveAmount:   reserve,
			LiabilityAmount: liability,
		})
		if rate < minRate {
			minRate = rate
		}
	}
	if len(rows) == 0 {
		return ParsedSummary{}, fmt.Errorf("phemex summary missing usable totalBalances rows")
	}
	if minRate == math.MaxFloat64 {
		minRate = 0
	}

	heights := make([]BlockHeight, 0, len(info.BlockHeights))
	for _, h := range info.BlockHeights {
		heights = append(heights, BlockHeight{Chain: h.Chain, Height: h.Height})
	}

	return ParsedSummary{
		AuditID:           auditID,
		OfficialVersionID: official,
		AuditTime:         ts,
		AuditTimeRaw:      tsRaw,
		DataVersion:       strings.TrimSpace(info.DataVersion),
		TotalReserveRate:  minRate,
		CoinRows:          rows,
		BlockHeights:      heights,
	}, nil
}

func parseVersionTime(desc, versionID string) (time.Time, string) {
	desc = strings.TrimSpace(desc)
	for _, layout := range []string{
		"2006-01-02 15:04:05 UTC",
		"2006-01-02 15:04:05",
		time.RFC3339,
	} {
		if t, err := time.Parse(layout, desc); err == nil {
			t = t.UTC()
			return t, t.Format(time.RFC3339)
		}
	}
	if t, err := time.Parse("2006-01-02", strings.TrimSpace(versionID)); err == nil {
		t = t.UTC()
		return t, t.Format(time.RFC3339)
	}
	return time.Time{}, desc
}

func compactSnapshotID(versionID string) string {
	s := strings.TrimSpace(versionID)
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, ":", "")
	s = strings.ReplaceAll(s, " ", "")
	if len(s) >= 8 {
		return s[:8]
	}
	return s
}

func normalizeVersionKey(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " utc", "")
	return strings.ReplaceAll(s, " ", "")
}

func getJSON(ctx context.Context, rawURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", URLPoRPage)
	req.Header.Set("Origin", "https://phemex.com")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: HTTP %d: %s", rawURL, resp.StatusCode, truncate(string(body), 200))
	}
	return body, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
