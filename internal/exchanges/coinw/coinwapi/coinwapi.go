// Package coinwapi fetches and parses CoinW Proof-of-Reserves summary artifacts.
// See docs/coinw-por-data-guide.md for data availability.
package coinwapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	userAgent = "ardmere/0.1 (+https://ardmere.org)"

	URLPoRPage = "https://www.coinw.com/proof-of-reserves"
	URLAnnouncementSep2026 = "https://www.coinw.com/help-center/announcement-center/maintenance-updates/project-upgrades/2026-09-26/proof-of-reserves-update-september-2026/8128"
	URLVerifierRelease = "https://github.com/connectCoinw/proof-of-reserves"
)

var (
	compactDateRe = regexp.MustCompile(`^\d{6,8}$`)
	cloudflareRe  = regexp.MustCompile(`(?i)Just a moment\.\.\.|cf-browser-verification|cdn-cgi/challenge`)
	ratioRe       = regexp.MustCompile(`(?i)total reserve ratio stands at\s*([0-9]+(?:\.[0-9]+)?)\s*%`)
	nuxtDataRe    = regexp.MustCompile(`(?s)<script[^>]*id="__NUXT_DATA__"[^>]*>(.*?)</script>`)
	hex64Re       = regexp.MustCompile(`\b([a-fA-F0-9]{64})\b`)
	snapTimeRe    = regexp.MustCompile(`20\d{2}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}`)
)

// SummaryBundle is the archived summarySnapshot for one CoinW PoR period.
type SummaryBundle struct {
	FetchedAt  time.Time       `json:"fetchedAt"`
	Source     string          `json:"source,omitempty"`
	SourceURLs []string        `json:"sourceUrls,omitempty"`
	ReportDate string          `json:"reportDate,omitempty"`
	SnapshotID string          `json:"snapshotId,omitempty"`
	Info       json.RawMessage `json:"info"`
	Probe      json.RawMessage `json:"probe,omitempty"`
}

// ParsedSummary is normalized fields for por.Snapshot.
// Reserve rates are percentages (100 = 100%) for solvency-claim.
type ParsedSummary struct {
	AuditID          string
	ReportDate       string
	AuditTime        time.Time
	AuditTimeRaw     string
	MerkleRoot       string
	TotalReserveRate float64
	SourceNote       string
	WAFBlocked       bool
	MerkleUnstable   bool
}

type infoJSON struct {
	AuditID              string          `json:"audit_id"`
	AuditIDCamel         string          `json:"auditId"`
	AuditTime            string          `json:"audit_time"`
	AuditTimeCamel       string          `json:"auditTime"`
	SnapshotTime         string          `json:"snapshotTime"`
	Timezone             string          `json:"timezone"`
	MerkleRoot           string          `json:"merkle_root_hash"`
	MerkleRootAlt        string          `json:"merkleRoot"`
	MerkleRootHash       string          `json:"merkleRootHash"`
	MerkleRootStability  string          `json:"merkle_root_stability"`
	MerkleRootAltObserved []string       `json:"merkle_root_alt_observed"`
	TotalReserveRate     float64         `json:"total_reserve_rate"`
	TotalReserveAlt      float64         `json:"totalReserveRate"`
	ReserveRate          float64         `json:"reserveRate"`
	TotalReserveStatus   string          `json:"totalReserveRateStatus"`
	SourceURL            string          `json:"sourceUrl"`
	Note                 string          `json:"note"`
	WAF                  string          `json:"waf"`
	Raw                  json.RawMessage `json:"-"`
}

// SummaryURL returns the primary public PoR page used for provenance.
func SummaryURL() string { return URLPoRPage }

// FetchPublicSummary probes CoinW public pages from this host.
// Datacenter IPs commonly receive Cloudflare WAF 403 / challenge HTML.
func FetchPublicSummary(ctx context.Context) ([]byte, ParsedSummary, error) {
	return FetchPublicSummaryByAuditID(ctx, "")
}

// FetchPublicSummaryByAuditID attempts a live HTML scrape. Empty auditID = latest probe.
// On WAF / missing machine fields, returns an error directing callers to -summary-path.
func FetchPublicSummaryByAuditID(ctx context.Context, auditID string) ([]byte, ParsedSummary, error) {
	auditID = strings.TrimSpace(auditID)
	porBody, porCode, porErr := getBody(ctx, URLPoRPage)
	annBody, annCode, annErr := getBody(ctx, URLAnnouncementSep2026)

	var notes []string
	if porErr != nil {
		notes = append(notes, fmt.Sprintf("por page: %v", porErr))
	} else {
		notes = append(notes, fmt.Sprintf("por page HTTP %d (%d bytes)", porCode, len(porBody)))
	}
	if annErr != nil {
		notes = append(notes, fmt.Sprintf("announcement: %v", annErr))
	} else {
		notes = append(notes, fmt.Sprintf("announcement HTTP %d (%d bytes)", annCode, len(annBody)))
	}

	waf := (porErr == nil && isCloudflareChallenge(porBody, porCode)) ||
		(annErr == nil && isCloudflareChallenge(annBody, annCode))

	// Prefer announcement ratio when HTML is readable (permissionless month claim).
	if annErr == nil && !isCloudflareChallenge(annBody, annCode) {
		if rate, ok := parseAnnouncementRatio(string(annBody)); ok {
			info := map[string]any{
				"audit_id":          firstNonEmpty(auditID, "202609"),
				"audit_time":        "2026-09-01 00:00:00",
				"timezone":          "UTC+8 (September 2026 month label; exact snapshot clock not in announcement)",
				"totalReserveRate":  rate,
				"sourceUrl":         URLAnnouncementSep2026,
				"waf":               "live announcement HTML reachable this run",
				"note":              "Parsed from CoinW September 2026 PoR announcement. Merkle root not published in the announcement body.",
				"liveProbeNotes":    strings.Join(notes, "; "),
			}
			if porErr == nil && !isCloudflareChallenge(porBody, porCode) {
				if root, snap := extractPoRPageHints(string(porBody)); root != "" || snap != "" {
					info["porPageMerkleHint"] = root
					info["porPageSnapshotHint"] = snap
				}
			}
			rawInfo, _ := json.Marshal(info)
			bundle := SummaryBundle{
				FetchedAt:  time.Now().UTC(),
				Source:     "announcement-html",
				SourceURLs: []string{URLAnnouncementSep2026, URLPoRPage, URLVerifierRelease},
				ReportDate: "2026-09",
				SnapshotID: firstNonEmpty(auditID, "202609"),
				Info:       rawInfo,
			}
			out, err := json.MarshalIndent(bundle, "", "  ")
			if err != nil {
				return nil, ParsedSummary{}, err
			}
			parsed, err := ParseSummaryBytes(out)
			return out, parsed, err
		}
	}

	msg := "CoinW PoR JSON is not stably permissionless from this host"
	if waf {
		msg = "CoinW Cloudflare WAF challenge/403 on public PoR pages from datacenter"
	}
	return nil, ParsedSummary{WAFBlocked: waf}, fmt.Errorf("%s (%s); import announcement/browser fixture via -summary-path (see fixtures/coinw/)", msg, strings.Join(notes, "; "))
}

// ParseSummaryBytes parses an archived SummaryBundle or raw info object.
func ParseSummaryBytes(raw []byte) (ParsedSummary, error) {
	var bundle SummaryBundle
	if err := json.Unmarshal(raw, &bundle); err == nil && len(bundle.Info) > 0 {
		parsed, err := parseInfo(bundle.Info)
		if err != nil {
			return ParsedSummary{}, err
		}
		if parsed.AuditID == "" && bundle.SnapshotID != "" {
			parsed.AuditID = bundle.SnapshotID
		}
		if parsed.ReportDate == "" {
			parsed.ReportDate = bundle.ReportDate
		}
		return parsed, nil
	}
	return parseInfo(json.RawMessage(raw))
}

// BuildSummaryBundle wraps a browser-captured info JSON object as an archive bundle.
func BuildSummaryBundle(infoFile string) ([]byte, error) {
	raw, err := os.ReadFile(infoFile)
	if err != nil {
		return nil, err
	}
	var asMap map[string]json.RawMessage
	if err := json.Unmarshal(raw, &asMap); err == nil {
		if info, ok := asMap["info"]; ok && len(info) > 0 {
			raw = info
		}
	}
	parsed, err := parseInfo(json.RawMessage(raw))
	if err != nil {
		return nil, err
	}
	bundle := SummaryBundle{
		FetchedAt:  time.Now().UTC(),
		Source:     "browser-or-announcement-import",
		SourceURLs: []string{URLPoRPage, URLAnnouncementSep2026, URLVerifierRelease},
		ReportDate: parsed.ReportDate,
		SnapshotID: parsed.AuditID,
		Info:       json.RawMessage(raw),
	}
	return json.MarshalIndent(bundle, "", "  ")
}

func parseInfo(infoRaw json.RawMessage) (ParsedSummary, error) {
	var info infoJSON
	if err := json.Unmarshal(infoRaw, &info); err != nil {
		return ParsedSummary{}, fmt.Errorf("decode coinw info: %w", err)
	}
	out := ParsedSummary{
		AuditID:          firstNonEmpty(info.AuditID, info.AuditIDCamel),
		AuditTimeRaw:     firstNonEmpty(info.AuditTime, info.AuditTimeCamel, info.SnapshotTime),
		MerkleRoot:       strings.ToLower(strings.TrimPrefix(firstNonEmpty(info.MerkleRoot, info.MerkleRootAlt, info.MerkleRootHash), "0x")),
		TotalReserveRate: firstNonZero(info.TotalReserveRate, info.TotalReserveAlt, info.ReserveRate),
		SourceNote:       info.Note,
		WAFBlocked:       strings.Contains(strings.ToLower(info.WAF), "waf") || strings.Contains(strings.ToLower(info.WAF), "cloudflare"),
		MerkleUnstable:   strings.Contains(strings.ToLower(info.MerkleRootStability), "animation") || strings.Contains(strings.ToLower(info.MerkleRootStability), "conflict") || len(info.MerkleRootAltObserved) > 0,
	}
	if out.AuditID == "" && out.AuditTimeRaw != "" {
		if t, err := parseAuditTime(out.AuditTimeRaw); err == nil {
			out.AuditID = t.Format("200601")
		}
	}
	if out.AuditID != "" && compactDateRe.MatchString(out.AuditID) {
		// keep
	} else if out.AuditID == "" {
		out.AuditID = "unknown"
	}
	if t, err := parseAuditTime(out.AuditTimeRaw); err == nil {
		out.AuditTime = t
		if out.ReportDate == "" {
			out.ReportDate = t.Format("2006-01")
		}
	}
	if out.ReportDate == "" && len(out.AuditID) >= 6 {
		out.ReportDate = out.AuditID[:4] + "-" + out.AuditID[4:6]
	}
	// Animated / unsettled UI ratios must not drive solvency-claim PASS.
	if strings.EqualFold(info.TotalReserveStatus, "UNVERIFIABLE_UI_ANIMATION") {
		out.TotalReserveRate = 0
	}
	return out, nil
}

func parseAuditTime(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, fmt.Errorf("empty audit time")
	}
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05 -07:00",
		"2006-01-02",
		"20060102",
		"200601",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t.UTC(), nil
		}
	}
	// CoinW PoR page uses UTC+8 wall clock without offset in the label.
	if t, err := time.ParseInLocation("2006-01-02 15:04:05", raw, time.FixedZone("UTC+8", 8*3600)); err == nil {
		return t.UTC(), nil
	}
	return time.Time{}, fmt.Errorf("parse audit time %q", raw)
}

func getBody(ctx context.Context, url string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/json;q=0.9,*/*;q=0.8")
	resp, err := httpClient().Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return body, resp.StatusCode, nil
}

func httpClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}

func isCloudflareChallenge(body []byte, code int) bool {
	if code == 403 || code == 503 {
		return true
	}
	return cloudflareRe.Match(body)
}

func parseAnnouncementRatio(html string) (float64, bool) {
	m := ratioRe.FindStringSubmatch(html)
	if len(m) < 2 {
		return 0, false
	}
	rate, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return 0, false
	}
	return rate, true
}

func extractPoRPageHints(html string) (merkle, snapshot string) {
	if m := nuxtDataRe.FindStringSubmatch(html); len(m) > 1 {
		html = m[1] + "\n" + html
	}
	if m := snapTimeRe.FindString(html); m != "" {
		snapshot = m
	}
	if m := hex64Re.FindString(html); m != "" {
		merkle = strings.ToLower(m)
	}
	return merkle, snapshot
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func firstNonZero(vals ...float64) float64 {
	for _, v := range vals {
		if v != 0 {
			return v
		}
	}
	return 0
}
