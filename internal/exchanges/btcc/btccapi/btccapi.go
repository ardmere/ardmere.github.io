// Package btccapi fetches and parses BTCC Proof-of-Reserves summary artifacts.
// See docs/btcc-por-data-guide.md for data availability.
package btccapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	userAgent = "ardmere/0.1 (+https://ardmere.org)"

	URLPoRPage     = "https://www.btcc.com/en-US/proof-of-reserves"
	URLReserveInfo = "https://www.btcc.com/social-api/contract/reserve/reserveInfo"
	// Login-gated (token required):
	URLReserveHistory = "https://www.btcc.com/social-api/contract/reserve/getReserveHistoryPage"
	URLReserveValidate = "https://www.btcc.com/social-api/contract/reserve/validate"
)

// Preferred coin order for display (matches BTCC PoR page).
var preferredCoinOrder = []string{"BTC", "ETH", "USDT", "USDC", "XRP", "ADA"}

// SummaryBundle is the archived summarySnapshot for one BTCC PoR period.
type SummaryBundle struct {
	FetchedAt  time.Time       `json:"fetchedAt"`
	Source     string          `json:"source,omitempty"`
	SourceURLs []string        `json:"sourceUrls,omitempty"`
	ReportDate string          `json:"reportDate,omitempty"`
	SnapshotID string          `json:"snapshotId,omitempty"`
	Info       json.RawMessage `json:"info"`
}

// ParsedSummary is normalized fields extracted from reserveInfo.
// Reserve rates are percentages (100 = 100%) for solvency-claim.
type ParsedSummary struct {
	AuditID      string
	ReportDate   string
	AuditTime    time.Time
	AuditTimeRaw string
	MerkleRoot   string // truncated ~16 hex as published
	TotalReserveRate float64
	CoinRows     []CoinRow
}

type CoinRow struct {
	Coin        string
	ReserveRate float64
}

type apiEnvelope struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	OK   bool            `json:"ok"`
	Data json.RawMessage `json:"data"`
}

type reserveInfoJSON struct {
	CreateDate string            `json:"createDate"`
	RootHash   string            `json:"rootHash"`
	Rate       map[string]string `json:"rate"`
}

var compactDateRe = regexp.MustCompile(`^\d{8}$`)
var percentRe = regexp.MustCompile(`^\s*([0-9]+(?:\.[0-9]+)?)\s*%?\s*$`)

// SummaryURL returns the public summary API used for artifact provenance.
func SummaryURL() string { return URLReserveInfo }

// FetchPublicSummary downloads reserveInfo for the latest known day-15 report.
func FetchPublicSummary(ctx context.Context) ([]byte, ParsedSummary, error) {
	return FetchPublicSummaryByReportDate(ctx, "")
}

// FetchPublicSummaryByReportDate selects a reportDate (e.g. "2026-10-15" or "20261015").
// Empty reportDate walks backward from today's preferred day-15 until a non-null payload is found.
func FetchPublicSummaryByReportDate(ctx context.Context, reportDate string) ([]byte, ParsedSummary, error) {
	reportDate = strings.TrimSpace(reportDate)
	if reportDate != "" {
		normalized, err := normalizeReportDate(reportDate)
		if err != nil {
			return nil, ParsedSummary{}, err
		}
		return fetchOne(ctx, normalized)
	}

	// Prefer the UI rule: on/after day 15 of month → that month's 15th; else prior month.
	candidates := day15Candidates(time.Now().UTC(), 24)
	var lastErr error
	for _, d := range candidates {
		raw, parsed, err := fetchOne(ctx, d)
		if err != nil {
			lastErr = err
			continue
		}
		if parsed.ReportDate != "" {
			return raw, parsed, nil
		}
	}
	if lastErr != nil {
		return nil, ParsedSummary{}, lastErr
	}
	return nil, ParsedSummary{}, fmt.Errorf("btcc reserveInfo: no published report found in recent day-15 window")
}

// ParseSummaryBytes parses an archived SummaryBundle or a raw reserveInfo payload.
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
		if data, ok := top["data"]; ok {
			return parseInfo(data)
		}
	}
	return parseInfo(json.RawMessage(raw))
}

// BuildSummaryBundle wraps a browser-captured reserveInfo data object (or full API envelope) as an archive bundle.
func BuildSummaryBundle(infoFile string) ([]byte, error) {
	raw, err := os.ReadFile(infoFile)
	if err != nil {
		return nil, err
	}
	var env apiEnvelope
	if err := json.Unmarshal(raw, &env); err == nil && len(env.Data) > 0 && string(env.Data) != "null" {
		raw = env.Data
	}
	out, _, err := buildBundle(json.RawMessage(raw), URLReserveInfo)
	return out, err
}

func fetchOne(ctx context.Context, reportDate string) ([]byte, ParsedSummary, error) {
	u, err := url.Parse(URLReserveInfo)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	q := u.Query()
	q.Set("reportDate", reportDate)
	u.RawQuery = q.Encode()

	raw, err := getJSON(ctx, u.String())
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	var env apiEnvelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, ParsedSummary{}, fmt.Errorf("decode reserveInfo envelope: %w", err)
	}
	if env.Code != 200 && !env.OK {
		return nil, ParsedSummary{}, fmt.Errorf("reserveInfo %s: code=%d msg=%s", reportDate, env.Code, env.Msg)
	}
	if len(env.Data) == 0 || string(env.Data) == "null" {
		return nil, ParsedSummary{}, fmt.Errorf("reserveInfo %s: empty data (not published)", reportDate)
	}
	return buildBundle(env.Data, u.String())
}

func buildBundle(infoRaw json.RawMessage, source string) ([]byte, ParsedSummary, error) {
	parsed, err := parseInfo(infoRaw)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	bundle := SummaryBundle{
		FetchedAt:  time.Now().UTC(),
		Source:     source,
		SourceURLs: []string{URLReserveInfo, URLPoRPage},
		ReportDate: parsed.ReportDate,
		SnapshotID: parsed.AuditID,
		Info:       infoRaw,
	}
	out, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	return out, parsed, nil
}

func parseInfo(infoRaw json.RawMessage) (ParsedSummary, error) {
	var info reserveInfoJSON
	if err := json.Unmarshal(infoRaw, &info); err != nil {
		return ParsedSummary{}, fmt.Errorf("decode btcc reserveInfo: %w", err)
	}
	createDate := strings.TrimSpace(info.CreateDate)
	if createDate == "" {
		return ParsedSummary{}, fmt.Errorf("btcc summary missing createDate")
	}
	ts, err := time.Parse("2006-01-02", createDate)
	if err != nil {
		return ParsedSummary{}, fmt.Errorf("parse createDate %q: %w", createDate, err)
	}
	ts = ts.UTC()
	auditID := strings.ReplaceAll(createDate, "-", "")

	root := strings.TrimSpace(strings.ToLower(info.RootHash))
	if root == "" {
		return ParsedSummary{}, fmt.Errorf("btcc summary missing rootHash")
	}

	rows := make([]CoinRow, 0, len(info.Rate))
	minRate := math.MaxFloat64
	seen := map[string]CoinRow{}
	for coin, rawRate := range info.Rate {
		coin = strings.TrimSpace(strings.ToUpper(coin))
		if coin == "" {
			continue
		}
		rate, err := parsePercent(rawRate)
		if err != nil {
			return ParsedSummary{}, fmt.Errorf("parse rate for %s (%q): %w", coin, rawRate, err)
		}
		row := CoinRow{Coin: coin, ReserveRate: rate}
		seen[coin] = row
		if rate < minRate {
			minRate = rate
		}
	}
	if len(seen) == 0 {
		return ParsedSummary{}, fmt.Errorf("btcc summary missing rate map")
	}
	for _, coin := range preferredCoinOrder {
		if row, ok := seen[coin]; ok {
			rows = append(rows, row)
			delete(seen, coin)
		}
	}
	extras := make([]string, 0, len(seen))
	for coin := range seen {
		extras = append(extras, coin)
	}
	sort.Strings(extras)
	for _, coin := range extras {
		rows = append(rows, seen[coin])
	}
	if minRate == math.MaxFloat64 {
		minRate = 0
	}

	return ParsedSummary{
		AuditID:          auditID,
		ReportDate:       createDate,
		AuditTime:        ts,
		AuditTimeRaw:     ts.Format(time.RFC3339),
		MerkleRoot:       root,
		TotalReserveRate: minRate,
		CoinRows:         rows,
	}, nil
}

func parsePercent(s string) (float64, error) {
	m := percentRe.FindStringSubmatch(s)
	if m == nil {
		return 0, fmt.Errorf("not a percentage")
	}
	return strconv.ParseFloat(m[1], 64)
}

func normalizeReportDate(s string) (string, error) {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "/", "-")
	if compactDateRe.MatchString(s) {
		return s[:4] + "-" + s[4:6] + "-" + s[6:8], nil
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t.Format("2006-01-02"), nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.UTC().Format("2006-01-02"), nil
	}
	return "", fmt.Errorf("unsupported reportDate %q (want YYYY-MM-DD or YYYYMMDD)", s)
}

func day15Candidates(now time.Time, monthsBack int) []string {
	y, m, d := now.Date()
	var start time.Time
	if d >= 15 {
		start = time.Date(y, m, 15, 0, 0, 0, 0, time.UTC)
	} else {
		start = time.Date(y, m-1, 15, 0, 0, 0, 0, time.UTC)
	}
	out := make([]string, 0, monthsBack)
	cur := start
	for i := 0; i < monthsBack; i++ {
		out = append(out, cur.Format("2006-01-02"))
		cur = cur.AddDate(0, -1, 0)
		cur = time.Date(cur.Year(), cur.Month(), 15, 0, 0, 0, 0, time.UTC)
	}
	return out
}

func getJSON(ctx context.Context, rawURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", URLPoRPage)
	req.Header.Set("Origin", "https://www.btcc.com")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
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
