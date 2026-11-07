// Package bitunixapi fetches and parses Bitunix Proof-of-Reserves summary artifacts.
// See docs/bitunix-por-data-guide.md for data availability.
package bitunixapi

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

	URLPoRPage      = "https://www.bitunix.com/proof-of-reserve"
	URLHelpVerify   = "https://www.bitunix.com/hub/helpcenter/article/how-to-verify-my-assets-on-bitunix?id=166"
	URLOSSVerifier  = "https://github.com/BitunixOfficial/proof-of-reserves"
	URLAPIBase      = "https://api.bitunix.com/asset-hub-api/api/public/v1/reserve/certificate"
	URLAuditTimes   = URLAPIBase + "/queryAuditTimeList"
	URLAuditInfo    = URLAPIBase + "/queryAuditInfoByTime"
	URLUserAuditIDs = URLAPIBase + "/queryUserAuditIdList"
)

// SummaryBundle is the archived summarySnapshot for one Bitunix PoR period.
type SummaryBundle struct {
	FetchedAt time.Time       `json:"fetchedAt"`
	Source    string          `json:"source,omitempty"`
	AuditTime string          `json:"auditTime"`
	Info      json.RawMessage `json:"info"`
}

// ParsedSummary is normalized fields extracted from Bitunix audit info.
// Reserve rates are percentages (100 = 100%) for solvency-claim.
type ParsedSummary struct {
	AuditID          string
	AuditTime        time.Time
	AuditTimeRaw     string
	MerkleRoot       string
	TotalReserveRate float64
	WalletAddress    string
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
	Success bool            `json:"success"`
}

type infoJSON struct {
	BTCReserveRatio  flexString `json:"btcReserveRatio"`
	ETHReserveRatio  flexString `json:"ethReserveRatio"`
	USDTReserveRatio flexString `json:"usdtReserveRatio"`
	MerkleRootHash   string     `json:"merkleRootHash"`
	WalletAddress    string     `json:"walletAddress"`
}

type flexString string

func (f *flexString) UnmarshalJSON(raw []byte) error {
	raw = []byte(strings.TrimSpace(string(raw)))
	if len(raw) == 0 || string(raw) == "null" {
		*f = ""
		return nil
	}
	if raw[0] == '"' {
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return err
		}
		*f = flexString(strings.TrimSpace(strings.TrimSuffix(s, "%")))
		return nil
	}
	*f = flexString(string(raw))
	return nil
}

// SummaryURL returns the public summary API used for artifact provenance.
func SummaryURL() string { return URLAuditInfo }

// FetchPublicSummary downloads the latest public audit info and archives it.
func FetchPublicSummary(ctx context.Context) ([]byte, ParsedSummary, error) {
	timesRaw, err := getJSON(ctx, URLAuditTimes)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	times, err := unwrapStringList(timesRaw)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	if len(times) == 0 {
		return nil, ParsedSummary{}, fmt.Errorf("bitunix audit time list empty")
	}
	auditTime := times[0]
	infoURL := URLAuditInfo + "?auditTime=" + url.QueryEscape(auditTime)
	infoRaw, err := getJSON(ctx, infoURL)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	infoBody, err := unwrapData(infoRaw)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	return buildBundle(auditTime, infoBody, infoURL)
}

// ParseSummaryBytes parses an archived SummaryBundle or a raw audit-info object.
func ParseSummaryBytes(raw []byte) (ParsedSummary, error) {
	var bundle SummaryBundle
	if err := json.Unmarshal(raw, &bundle); err == nil && len(bundle.Info) > 0 {
		return parseInfo(bundle.AuditTime, bundle.Info)
	}
	// Accept a bare info object (needs auditTime sibling or embedded).
	var top map[string]json.RawMessage
	if err := json.Unmarshal(raw, &top); err == nil {
		if info, ok := top["info"]; ok {
			at := ""
			if v, ok := top["auditTime"]; ok {
				_ = json.Unmarshal(v, &at)
			}
			return parseInfo(at, info)
		}
	}
	return parseInfo("", json.RawMessage(raw))
}

// BuildSummaryBundle wraps a browser-captured audit-info JSON as an archive bundle.
func BuildSummaryBundle(auditTime, infoFile string) ([]byte, error) {
	raw, err := readFile(infoFile)
	if err != nil {
		return nil, err
	}
	if data, err := unwrapData(raw); err == nil {
		raw = data
	}
	if auditTime == "" {
		var probe struct {
			AuditTime string `json:"auditTime"`
		}
		_ = json.Unmarshal(raw, &probe)
		auditTime = probe.AuditTime
	}
	out, _, err := buildBundle(auditTime, raw, URLAuditInfo)
	return out, err
}

func buildBundle(auditTime string, infoRaw json.RawMessage, source string) ([]byte, ParsedSummary, error) {
	parsed, err := parseInfo(auditTime, infoRaw)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	bundle := SummaryBundle{
		FetchedAt: time.Now().UTC(),
		Source:    source,
		AuditTime: auditTime,
		Info:      infoRaw,
	}
	out, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	return out, parsed, nil
}

func parseInfo(auditTime string, infoRaw json.RawMessage) (ParsedSummary, error) {
	var info infoJSON
	if err := json.Unmarshal(infoRaw, &info); err != nil {
		return ParsedSummary{}, fmt.Errorf("decode bitunix audit info: %w", err)
	}
	root := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(info.MerkleRootHash), "0x"))
	btc := mustFloat(string(info.BTCReserveRatio))
	eth := mustFloat(string(info.ETHReserveRatio))
	usdt := mustFloat(string(info.USDTReserveRatio))
	if root == "" && btc == 0 && eth == 0 && usdt == 0 {
		return ParsedSummary{}, fmt.Errorf("bitunix summary missing merkle root and reserve ratios")
	}

	ts, tsRaw := parseAuditTime(auditTime)
	auditID := ts.Format("20060102")
	if ts.IsZero() && auditTime != "" {
		auditID = strings.ReplaceAll(auditTime, ":", "")
	}

	rows := []CoinRow{
		{Coin: "BTC", ReserveRate: btc},
		{Coin: "ETH", ReserveRate: eth},
		{Coin: "USDT", ReserveRate: usdt},
	}
	minRate := btc
	for _, r := range rows[1:] {
		if r.ReserveRate < minRate {
			minRate = r.ReserveRate
		}
	}
	return ParsedSummary{
		AuditID:          auditID,
		AuditTime:        ts,
		AuditTimeRaw:     tsRaw,
		MerkleRoot:       root,
		TotalReserveRate: minRate,
		WalletAddress:    strings.TrimSpace(info.WalletAddress),
		CoinRows:         rows,
	}, nil
}

func parseAuditTime(s string) (time.Time, string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, ""
	}
	layouts := []string{
		time.RFC3339Nano,
		"2006-01-02T15:04:05.000000Z",
		"2006-01-02T15:04:05Z",
		time.RFC3339,
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			t = t.UTC()
			return t, t.Format(time.RFC3339)
		}
	}
	return time.Time{}, s
}

func unwrapStringList(raw []byte) ([]string, error) {
	data, err := unwrapData(raw)
	if err != nil {
		return nil, err
	}
	var times []string
	if err := json.Unmarshal(data, &times); err != nil {
		return nil, fmt.Errorf("decode bitunix audit times: %w", err)
	}
	return times, nil
}

func unwrapData(raw []byte) (json.RawMessage, error) {
	var env apiEnvelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, fmt.Errorf("decode bitunix envelope: %w", err)
	}
	if len(env.Data) == 0 {
		// Already unwrapped.
		return json.RawMessage(raw), nil
	}
	if !codeOK(env.Code) {
		msg := env.Msg
		if msg == "" {
			msg = env.Message
		}
		return nil, fmt.Errorf("bitunix API code=%v message=%s", env.Code, msg)
	}
	return env.Data, nil
}

func codeOK(code any) bool {
	switch v := code.(type) {
	case float64:
		return v == 0
	case int:
		return v == 0
	case string:
		return v == "0" || v == "" || strings.EqualFold(v, "ok") || strings.EqualFold(v, "success")
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

func readFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func mustFloat(s string) float64 {
	s = strings.TrimSpace(strings.TrimSuffix(s, "%"))
	if s == "" {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
