// Package toobitapi fetches and parses Toobit Proof-of-Reserves summary artifacts.
// See docs/toobit-por-data-guide.md for data availability.
package toobitapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	userAgent = "ardmere/0.1 (+https://ardmere.org)"

	URLPoRPage     = "https://www.toobit.com/proof-of-reserves"
	URLHelpVerify  = "https://www.toobit.com/support/how-do-i-view-and-verify-my-toobit-assets"
	URLOSSVerifier = "https://github.com/toobit-docs/por-verifier"
	URLAPIBase     = "https://bapi.toobit.com/uapi/v1/por"
	URLLatestAudit = URLAPIBase + "/latest-audit"
)

// SummaryBundle is the archived summarySnapshot for one Toobit PoR period.
type SummaryBundle struct {
	FetchedAt        time.Time       `json:"fetchedAt"`
	Source           string          `json:"source,omitempty"`
	OfficialAuditID  string          `json:"officialAuditId,omitempty"`
	Info             json.RawMessage `json:"info"`
}

// ParsedSummary is normalized fields extracted from Toobit latest-audit data.
// Reserve rates are percentages (100 = 100%) for solvency-claim.
type ParsedSummary struct {
	AuditID          string
	OfficialAuditID  string
	AuditTime        time.Time
	AuditTimeRaw     string
	MerkleRoot       string
	TreeHeight       int
	TotalReserveRate float64
	CoinRows         []CoinRow
}

type CoinRow struct {
	Coin        string
	ReserveRate float64
}

type apiEnvelope struct {
	Code    any             `json:"code"`
	Msg     string          `json:"msg"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type infoJSON struct {
	AuditID      string         `json:"auditId"`
	SnapshotTime flexInt64      `json:"snapshotTime"`
	MerkleRoot   string         `json:"merkleRoot"`
	TreeHeight   int            `json:"treeHeight"`
	ReserveInfo  []reserveToken `json:"reserveInfo"`
}

type reserveToken struct {
	TokenID      string  `json:"tokenId"`
	ReserveRatio float64 `json:"reserveRatio"`
}

type flexInt64 int64

func (f *flexInt64) UnmarshalJSON(raw []byte) error {
	raw = []byte(strings.TrimSpace(string(raw)))
	if len(raw) == 0 || string(raw) == "null" {
		*f = 0
		return nil
	}
	if raw[0] == '"' {
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return err
		}
		s = strings.TrimSpace(s)
		if s == "" {
			*f = 0
			return nil
		}
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		*f = flexInt64(v)
		return nil
	}
	var v int64
	if err := json.Unmarshal(raw, &v); err != nil {
		return err
	}
	*f = flexInt64(v)
	return nil
}

// SummaryURL returns the public summary API used for artifact provenance.
func SummaryURL() string { return URLLatestAudit }

// FetchPublicSummary downloads the latest public audit info and archives it.
func FetchPublicSummary(ctx context.Context) ([]byte, ParsedSummary, error) {
	raw, err := getJSON(ctx, URLLatestAudit)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	data, err := unwrapData(raw)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	return buildBundle(data, URLLatestAudit)
}

// FetchPublicSummaryByAuditID downloads a specific audit via ?auditId=.
func FetchPublicSummaryByAuditID(ctx context.Context, auditID string) ([]byte, ParsedSummary, error) {
	auditID = strings.TrimSpace(auditID)
	if auditID == "" {
		return FetchPublicSummary(ctx)
	}
	u := URLLatestAudit + "?auditId=" + url.QueryEscape(auditID)
	raw, err := getJSON(ctx, u)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	data, err := unwrapData(raw)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	return buildBundle(data, u)
}

// ParseSummaryBytes parses an archived SummaryBundle or a raw latest-audit data object.
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

// BuildSummaryBundle wraps a browser-captured latest-audit JSON as an archive bundle.
func BuildSummaryBundle(infoFile string) ([]byte, error) {
	raw, err := os.ReadFile(infoFile)
	if err != nil {
		return nil, err
	}
	if data, err := unwrapData(raw); err == nil {
		raw = data
	}
	out, _, err := buildBundle(raw, URLLatestAudit)
	return out, err
}

func buildBundle(infoRaw json.RawMessage, source string) ([]byte, ParsedSummary, error) {
	parsed, err := parseInfo(infoRaw)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	bundle := SummaryBundle{
		FetchedAt:       time.Now().UTC(),
		Source:          source,
		OfficialAuditID: parsed.OfficialAuditID,
		Info:            infoRaw,
	}
	out, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	return out, parsed, nil
}

func parseInfo(infoRaw json.RawMessage) (ParsedSummary, error) {
	var info infoJSON
	if err := json.Unmarshal(infoRaw, &info); err != nil {
		return ParsedSummary{}, fmt.Errorf("decode toobit audit info: %w", err)
	}
	root := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(info.MerkleRoot), "0x"))
	if root == "" && len(info.ReserveInfo) == 0 {
		return ParsedSummary{}, fmt.Errorf("toobit summary missing merkle root and reserve ratios")
	}

	ts := time.Time{}
	tsRaw := ""
	if info.SnapshotTime > 0 {
		ms := int64(info.SnapshotTime)
		// Accept seconds or milliseconds.
		if ms < 1_000_000_000_000 {
			ts = time.Unix(ms, 0).UTC()
		} else {
			ts = time.UnixMilli(ms).UTC()
		}
		tsRaw = ts.Format(time.RFC3339)
	}

	official := strings.TrimSpace(info.AuditID)
	auditID := ""
	if !ts.IsZero() {
		auditID = ts.Format("20060102")
	} else if official != "" {
		auditID = strings.ReplaceAll(official, ":", "")
	}

	rows := make([]CoinRow, 0, len(info.ReserveInfo))
	minRate := 0.0
	for i, r := range info.ReserveInfo {
		coin := strings.TrimSpace(r.TokenID)
		if coin == "" {
			continue
		}
		rows = append(rows, CoinRow{Coin: coin, ReserveRate: r.ReserveRatio})
		if i == 0 || r.ReserveRatio < minRate {
			minRate = r.ReserveRatio
		}
	}
	if len(rows) == 0 {
		return ParsedSummary{}, fmt.Errorf("toobit summary missing reserveInfo rows")
	}

	return ParsedSummary{
		AuditID:          auditID,
		OfficialAuditID:  official,
		AuditTime:        ts,
		AuditTimeRaw:     tsRaw,
		MerkleRoot:       root,
		TreeHeight:       info.TreeHeight,
		TotalReserveRate: minRate,
		CoinRows:         rows,
	}, nil
}

func unwrapData(raw []byte) (json.RawMessage, error) {
	var env apiEnvelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, fmt.Errorf("decode toobit envelope: %w", err)
	}
	if len(env.Data) == 0 {
		return json.RawMessage(raw), nil
	}
	if !codeOK(env.Code) {
		msg := env.Msg
		if msg == "" {
			msg = env.Message
		}
		return nil, fmt.Errorf("toobit API code=%v message=%s", env.Code, msg)
	}
	return env.Data, nil
}

func codeOK(code any) bool {
	switch v := code.(type) {
	case float64:
		return v == 0 || v == 200
	case int:
		return v == 0 || v == 200
	case string:
		return v == "0" || v == "200" || v == "" || strings.EqualFold(v, "ok") || strings.EqualFold(v, "success")
	case nil:
		return true
	default:
		return false
	}
}

func getJSON(ctx context.Context, rawURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", URLPoRPage)
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
