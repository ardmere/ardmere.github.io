// Package wooxapi fetches and parses WOO X Proof-of-Reserves summary artifacts.
// See docs/woox-por-data-guide.md for data availability.
package wooxapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	// Prefer a browser-like UA: bare library UAs intermittently see HTTP 500 from api.woox.io.
	userAgent = "Mozilla/5.0 (compatible; ardmere/0.1; +https://ardmere.org)"

	URLPoRPage   = "https://woox.io/proof-of-reserves"
	URLAPIBase   = "https://api.woox.io"
	URLAudits    = URLAPIBase + "/asset/porl/v2/audits"
	URLOverview  = URLAPIBase + "/asset/porl/v2/snapshot_overview"
	URLAssets    = URLAPIBase + "/asset/porl/assets"
	URLLiab      = URLAPIBase + "/asset/porl/liabilities"
	URLOSSVerify = "https://github.com/woo-x/merkle-tree"
)

// SummaryBundle is the archived summarySnapshot for one WOO X PoR period.
type SummaryBundle struct {
	FetchedAt  time.Time       `json:"fetchedAt"`
	Source     string          `json:"source,omitempty"`
	SourceURLs []string        `json:"sourceUrls,omitempty"`
	Audit      json.RawMessage `json:"audit"`
	Overview   json.RawMessage `json:"overview"`
	Live       json.RawMessage `json:"live,omitempty"`
	Notes      string          `json:"notes,omitempty"`
}

// ParsedSummary is normalized fields extracted from a WOO X audit overview.
type ParsedSummary struct {
	AuditID              string
	OfficialAuditID      string
	AuditTime            time.Time
	AuditTimeRaw         string
	MerkleTreeURL        string
	MerkleTreeBytes      int64
	MerkleTreeETag       string
	TotalReserveRate     float64
	TotalAssetsUSDT      float64
	TotalLiabilitiesUSDT float64
	AssetTokenCount      int
	LiabilityTokenCount  int
	LiveReserveRate      float64
	VenuesNote           string
}

type apiEnvelope struct {
	Success   bool            `json:"success"`
	Timestamp int64           `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
	Code      any             `json:"code"`
	Message   string          `json:"message"`
}

type auditRow struct {
	AuditID    string `json:"audit_id"`
	AuditDate  int64  `json:"audit_date"`
	MerkleTree string `json:"merkle_tree"`
}

type auditsData struct {
	Rows []auditRow `json:"rows"`
}

type tokenNotional struct {
	Token        string   `json:"token"`
	USDTNotional *float64 `json:"usdt_notional"`
	Quantity     *float64 `json:"quantity"`
}

type overviewData struct {
	Asset       []tokenNotional `json:"asset"`
	Liabilities []tokenNotional `json:"liabilities"`
}

type assetsData struct {
	Tokens            []tokenNotional `json:"tokens"`
	Venues            []venueRow      `json:"venues"`
	TotalUSDTNotional float64         `json:"total_usdt_notional"`
}

type venueRow struct {
	Name         string  `json:"name"`
	Type         string  `json:"type"`
	USDTNotional float64 `json:"usdt_notional"`
}

type liveAssetsEnvelope struct {
	Success     bool       `json:"success"`
	Data        assetsData `json:"data"`
	UpdatedTime int64      `json:"updated_time"`
}

type liveLiabEnvelope struct {
	Success     bool            `json:"success"`
	Data        []tokenNotional `json:"data"`
	UpdatedTime int64           `json:"updated_time"`
}

// SummaryURL returns the public summary API used for artifact provenance.
func SummaryURL() string { return URLOverview }

// FetchPublicSummary downloads the latest public audit overview (+ live dashboard totals).
func FetchPublicSummary(ctx context.Context) ([]byte, ParsedSummary, error) {
	auditsRaw, err := getJSON(ctx, URLAudits)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	auditsBody, err := unwrapData(auditsRaw)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	var audits auditsData
	if err := json.Unmarshal(auditsBody, &audits); err != nil {
		return nil, ParsedSummary{}, fmt.Errorf("decode woox audits: %w", err)
	}
	if len(audits.Rows) == 0 {
		return nil, ParsedSummary{}, fmt.Errorf("woox audit list empty")
	}
	latest := audits.Rows[0]
	for _, row := range audits.Rows[1:] {
		if row.AuditDate > latest.AuditDate {
			latest = row
		}
	}

	overviewURL := URLOverview + "?auditId=" + url.QueryEscape(latest.AuditID)
	overviewRaw, err := getJSON(ctx, overviewURL)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	overviewBody, err := unwrapData(overviewRaw)
	if err != nil {
		return nil, ParsedSummary{}, err
	}

	var ov overviewData
	if err := json.Unmarshal(overviewBody, &ov); err != nil {
		return nil, ParsedSummary{}, fmt.Errorf("decode woox overview: %w", err)
	}

	assetsTot, liabTot := sumNotional(ov.Asset), sumNotional(ov.Liabilities)
	ratio := reserveRatio(assetsTot, liabTot)

	merkleBytes, merkleETag, merkleMD5, _ := headMerkle(ctx, latest.MerkleTree)

	liveRaw, liveParsed, liveErr := fetchLive(ctx)
	_ = liveErr // live is supplementary; audit overview is authoritative for the snapshot

	auditObj := map[string]any{
		"audit_id":            latest.AuditID,
		"audit_date_ms":       latest.AuditDate,
		"audit_date":          msToRFC3339(latest.AuditDate),
		"merkle_tree":         latest.MerkleTree,
		"merkle_tree_bytes":   merkleBytes,
		"merkle_tree_etag":    merkleETag,
		"merkle_tree_md5_b64": merkleMD5,
	}
	overviewObj := map[string]any{
		"totalAssetsUSDT":      assetsTot,
		"totalLiabilitiesUSDT": liabTot,
		"reserveRatioPercent":  ratio,
		"assetTokenCount":      len(ov.Asset),
		"liabilityTokenCount":  len(ov.Liabilities),
		"topAssets":            topN(ov.Asset, 10),
		"topLiabilities":       topN(ov.Liabilities, 10),
	}

	auditJSON, _ := json.Marshal(auditObj)
	overviewJSON, _ := json.Marshal(overviewObj)

	bundle := SummaryBundle{
		FetchedAt: time.Now().UTC(),
		Source:    overviewURL,
		SourceURLs: []string{
			URLAudits, URLOverview, URLAssets, URLLiab, URLPoRPage, latest.MerkleTree, URLOSSVerify,
		},
		Audit:    auditJSON,
		Overview: overviewJSON,
		Notes: "Reserve ratio = sum(overview.asset.usdt_notional)/sum(overview.liabilities.usdt_notional)*100 " +
			"(same formula as woox.io PoR UI). Full Merkle tree is permissionlessly downloadable (~200MB protobuf); " +
			"user Merkle hash is login-gated. Official validator: woo-x/merkle-tree. Live dashboard fields are ~15-min " +
			"rolling and not the audit merkle period.",
	}
	if liveRaw != nil {
		bundle.Live = liveRaw
	}

	out, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	parsed := buildParsed(latest, assetsTot, liabTot, ratio, merkleBytes, merkleETag, liveParsed)
	return out, parsed, nil
}

// ParseSummaryBytes parses an archived SummaryBundle.
func ParseSummaryBytes(raw []byte) (ParsedSummary, error) {
	var bundle SummaryBundle
	if err := json.Unmarshal(raw, &bundle); err != nil {
		return ParsedSummary{}, fmt.Errorf("decode woox summary bundle: %w", err)
	}
	if len(bundle.Audit) == 0 || len(bundle.Overview) == 0 {
		return ParsedSummary{}, fmt.Errorf("woox summary missing audit/overview")
	}

	var audit struct {
		AuditID         string `json:"audit_id"`
		AuditDateMS     int64  `json:"audit_date_ms"`
		AuditDate       string `json:"audit_date"`
		MerkleTree      string `json:"merkle_tree"`
		MerkleTreeBytes int64  `json:"merkle_tree_bytes"`
		MerkleTreeETag  string `json:"merkle_tree_etag"`
	}
	if err := json.Unmarshal(bundle.Audit, &audit); err != nil {
		return ParsedSummary{}, fmt.Errorf("decode woox audit: %w", err)
	}
	var overview struct {
		TotalAssetsUSDT      float64 `json:"totalAssetsUSDT"`
		TotalLiabilitiesUSDT float64 `json:"totalLiabilitiesUSDT"`
		ReserveRatioPercent  float64 `json:"reserveRatioPercent"`
		AssetTokenCount      int     `json:"assetTokenCount"`
		LiabilityTokenCount  int     `json:"liabilityTokenCount"`
	}
	if err := json.Unmarshal(bundle.Overview, &overview); err != nil {
		return ParsedSummary{}, fmt.Errorf("decode woox overview: %w", err)
	}

	ts, tsRaw := parseAuditTime(audit.AuditDate, audit.AuditDateMS)
	auditID := ts.Format("20060102")
	if ts.IsZero() {
		auditID = strings.ReplaceAll(audit.AuditID, "-", "")[:8]
	}

	liveRate := 0.0
	venues := ""
	if len(bundle.Live) > 0 {
		var live struct {
			ReserveRatioPercent float64    `json:"reserveRatioPercent"`
			Venues              []venueRow `json:"venues"`
		}
		if err := json.Unmarshal(bundle.Live, &live); err == nil {
			liveRate = live.ReserveRatioPercent
			parts := make([]string, 0, len(live.Venues))
			for _, v := range live.Venues {
				parts = append(parts, fmt.Sprintf("%s(%s)", v.Name, v.Type))
			}
			venues = strings.Join(parts, "; ")
		}
	}

	ratio := overview.ReserveRatioPercent
	if ratio == 0 && overview.TotalLiabilitiesUSDT > 0 {
		ratio = reserveRatio(overview.TotalAssetsUSDT, overview.TotalLiabilitiesUSDT)
	}

	return ParsedSummary{
		AuditID:              auditID,
		OfficialAuditID:      audit.AuditID,
		AuditTime:            ts,
		AuditTimeRaw:         tsRaw,
		MerkleTreeURL:        audit.MerkleTree,
		MerkleTreeBytes:      audit.MerkleTreeBytes,
		MerkleTreeETag:       audit.MerkleTreeETag,
		TotalReserveRate:     ratio,
		TotalAssetsUSDT:      overview.TotalAssetsUSDT,
		TotalLiabilitiesUSDT: overview.TotalLiabilitiesUSDT,
		AssetTokenCount:      overview.AssetTokenCount,
		LiabilityTokenCount:  overview.LiabilityTokenCount,
		LiveReserveRate:      liveRate,
		VenuesNote:           venues,
	}, nil
}

func buildParsed(latest auditRow, assetsTot, liabTot, ratio float64, merkleBytes int64, merkleETag string, live ParsedSummary) ParsedSummary {
	ts, tsRaw := parseAuditTime("", latest.AuditDate)
	out := ParsedSummary{
		AuditID:              ts.Format("20060102"),
		OfficialAuditID:      latest.AuditID,
		AuditTime:            ts,
		AuditTimeRaw:         tsRaw,
		MerkleTreeURL:        latest.MerkleTree,
		MerkleTreeBytes:      merkleBytes,
		MerkleTreeETag:       merkleETag,
		TotalReserveRate:     ratio,
		TotalAssetsUSDT:      assetsTot,
		TotalLiabilitiesUSDT: liabTot,
		LiveReserveRate:      live.LiveReserveRate,
		VenuesNote:           live.VenuesNote,
	}
	return out
}

func fetchLive(ctx context.Context) (json.RawMessage, ParsedSummary, error) {
	assetsURL := URLAssets + "?excludeWooToken=false"
	liabURL := URLLiab + "?excludeWooToken=false"
	assetsRaw, err := getJSON(ctx, assetsURL)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	liabRaw, err := getJSON(ctx, liabURL)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	var assets liveAssetsEnvelope
	if err := json.Unmarshal(assetsRaw, &assets); err != nil {
		return nil, ParsedSummary{}, err
	}
	var liab liveLiabEnvelope
	if err := json.Unmarshal(liabRaw, &liab); err != nil {
		return nil, ParsedSummary{}, err
	}
	liveAssets := assets.Data.TotalUSDTNotional
	liveLiab := sumNotional(liab.Data)
	ratio := reserveRatio(liveAssets, liveLiab)
	parts := make([]string, 0, len(assets.Data.Venues))
	for _, v := range assets.Data.Venues {
		parts = append(parts, fmt.Sprintf("%s(%s)", v.Name, v.Type))
	}
	obj := map[string]any{
		"excludeWooToken":        false,
		"assetsUpdatedTime":      assets.UpdatedTime,
		"liabilitiesUpdatedTime": liab.UpdatedTime,
		"totalAssetsUSDT":        liveAssets,
		"totalLiabilitiesUSDT":   liveLiab,
		"reserveRatioPercent":    ratio,
		"venues":                 assets.Data.Venues,
	}
	raw, err := json.Marshal(obj)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	return raw, ParsedSummary{LiveReserveRate: ratio, VenuesNote: strings.Join(parts, "; ")}, nil
}

func headMerkle(ctx context.Context, rawURL string) (bytes int64, etag, md5b64 string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, rawURL, nil)
	if err != nil {
		return 0, "", "", err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, "", "", fmt.Errorf("HEAD %s: HTTP %d", rawURL, resp.StatusCode)
	}
	if cl := resp.Header.Get("Content-Length"); cl != "" {
		if n, e := strconv.ParseInt(cl, 10, 64); e == nil {
			bytes = n
		}
	}
	etag = strings.Trim(resp.Header.Get("ETag"), `"`)
	for _, h := range resp.Header.Values("X-Goog-Hash") {
		if strings.HasPrefix(h, "md5=") {
			md5b64 = strings.TrimPrefix(h, "md5=")
		}
	}
	return bytes, etag, md5b64, nil
}

func sumNotional(rows []tokenNotional) float64 {
	var sum float64
	for _, r := range rows {
		if r.USDTNotional != nil {
			sum += *r.USDTNotional
		}
	}
	return sum
}

func reserveRatio(assets, liabilities float64) float64 {
	if liabilities == 0 {
		return 0
	}
	// Match UI: (Number(assets)/Number(liabilities)*100).toFixed(1)
	return math.Round(assets/liabilities*1000) / 10
}

func topN(rows []tokenNotional, n int) []map[string]any {
	cp := append([]tokenNotional(nil), rows...)
	for i := 0; i < len(cp); i++ {
		for j := i + 1; j < len(cp); j++ {
			ai, aj := 0.0, 0.0
			if cp[i].USDTNotional != nil {
				ai = *cp[i].USDTNotional
			}
			if cp[j].USDTNotional != nil {
				aj = *cp[j].USDTNotional
			}
			if aj > ai {
				cp[i], cp[j] = cp[j], cp[i]
			}
		}
	}
	if len(cp) > n {
		cp = cp[:n]
	}
	out := make([]map[string]any, 0, len(cp))
	for _, r := range cp {
		m := map[string]any{"token": r.Token}
		if r.USDTNotional != nil {
			m["usdt_notional"] = *r.USDTNotional
		}
		if r.Quantity != nil {
			m["quantity"] = *r.Quantity
		}
		out = append(out, m)
	}
	return out
}

func parseAuditTime(s string, ms int64) (time.Time, string) {
	s = strings.TrimSpace(s)
	if s != "" {
		layouts := []string{
			time.RFC3339Nano,
			"2006-01-02T15:04:05.000Z",
			"2006-01-02T15:04:05Z",
			time.RFC3339,
		}
		for _, layout := range layouts {
			if t, err := time.Parse(layout, s); err == nil {
				t = t.UTC()
				return t, t.Format(time.RFC3339)
			}
		}
	}
	if ms > 0 {
		t := time.UnixMilli(ms).UTC()
		return t, t.Format(time.RFC3339)
	}
	return time.Time{}, s
}

func msToRFC3339(ms int64) string {
	return time.UnixMilli(ms).UTC().Format("2006-01-02T15:04:05.000Z")
}

func unwrapData(raw []byte) (json.RawMessage, error) {
	var env apiEnvelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, fmt.Errorf("decode woox envelope: %w", err)
	}
	if len(env.Data) == 0 {
		return json.RawMessage(raw), nil
	}
	if env.Message != "" && !env.Success {
		return nil, fmt.Errorf("woox API error: %s", env.Message)
	}
	if !env.Success && env.Code != nil {
		return nil, fmt.Errorf("woox API code=%v message=%s", env.Code, env.Message)
	}
	return env.Data, nil
}

func getJSON(ctx context.Context, rawURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Origin", "https://woox.io")
	req.Header.Set("Referer", URLPoRPage)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 16<<20))
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

// ReadFile is a thin helper for tests / capture import.
func ReadFile(path string) ([]byte, error) { return os.ReadFile(path) }
