// Package bybitapi parses Bybit Proof-of-Reserves summary artifacts.
// See docs/bybit-por-data-guide.md for data availability.
package bybitapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	userAgent = "ardmere/0.1 (+https://ardmere.org)"

	URLPoRPage = "https://www.bybit.com/en/proof-of-reserves"

	// Undocumented x-api path used by the PoR dashboard (WAF blocks most datacenter IPs).
	URLReserveRatio = "https://www.bybit.com/x-api/por/public/v1/reserve-ratio/latest"
)

// SummaryBundle is the combined summary we archive as one artifact.
type SummaryBundle struct {
	FetchedAt time.Time       `json:"fetchedAt"`
	Source    string          `json:"source,omitempty"`
	Info      json.RawMessage `json:"info"`
	CoinList  json.RawMessage `json:"coinList,omitempty"`
}

// ParsedSummary is normalized fields extracted from Bybit summary JSON.
type ParsedSummary struct {
	AuditID          string
	AuditTime        time.Time
	AuditTimeRaw     string
	MerkleRoot       string
	TotalReserveRate float64
	CoinRows         []CoinRow
}

type CoinRow struct {
	Coin            string
	ReserveRate     float64
	ReserveAmount   float64
	LiabilityAmount float64
}

type infoData struct {
	AuditID          string  `json:"audit_id"`
	AuditIDCamel     string  `json:"auditId"`
	AuditTime        string  `json:"audit_time"`
	AuditTimeCamel   string  `json:"auditTime"`
	SnapshotTime     string  `json:"snapshotTime"`
	MerkleRoot       string  `json:"merkle_root_hash"`
	MerkleRootAlt    string  `json:"merkleRoot"`
	MerkleRootHash   string  `json:"merkleRootHash"`
	TotalReserveRate float64 `json:"total_reserve_rate"`
	TotalReserveAlt  float64 `json:"totalReserveRate"`
	ReserveRate      float64 `json:"reserveRate"`
}

type coinListData struct {
	List []struct {
		Coin              string  `json:"coin"`
		Asset             string  `json:"asset"`
		ReserveRate       float64 `json:"reserve_rate"`
		ReserveRateCamel  float64 `json:"reserveRate"`
		ReserveAmount     float64 `json:"reserve_amount"`
		ReserveAmountAlt  float64 `json:"reserveAmount"`
		CustomerLiability float64 `json:"customer_liability"`
		LiabilityAmount   float64 `json:"liabilityAmount"`
	} `json:"list"`
	Coins []struct {
		Coin            string  `json:"coin"`
		ReserveRate     float64 `json:"reserveRate"`
		ReserveAmount   float64 `json:"reserveAmount"`
		LiabilityAmount float64 `json:"liabilityAmount"`
	} `json:"coins"`
}

// FetchPublicSummary attempts the PoR dashboard API; usually blocked without browser session.
func FetchPublicSummary(ctx context.Context) ([]byte, ParsedSummary, error) {
	infoRaw, err := get(ctx, URLReserveRatio)
	if err != nil {
		return nil, ParsedSummary{}, fmt.Errorf("fetch reserve ratio (WAF may block; use -summary-path): %w", err)
	}
	parsed, err := ParseSummaryBundle(infoRaw, nil)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	bundle := SummaryBundle{
		FetchedAt: time.Now().UTC(),
		Info:      json.RawMessage(infoRaw),
	}
	out, err := json.Marshal(bundle)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	return out, parsed, nil
}

// ParseSummaryBytes parses a saved SummaryBundle or raw API response.
func ParseSummaryBytes(raw []byte) (ParsedSummary, error) {
	var bundle SummaryBundle
	if err := json.Unmarshal(raw, &bundle); err == nil && len(bundle.Info) > 0 {
		return ParseSummaryBundle(bundle.Info, bundle.CoinList)
	}
	return ParseSummaryBundle(raw, nil)
}

// BuildSummaryBundle merges browser-captured info + optional coin list files.
func BuildSummaryBundle(infoPath, coinsPath string) ([]byte, error) {
	infoRaw, err := os.ReadFile(infoPath)
	if err != nil {
		return nil, err
	}
	var coinsRaw []byte
	if coinsPath != "" {
		coinsRaw, err = os.ReadFile(coinsPath)
		if err != nil {
			return nil, err
		}
	}
	bundle := SummaryBundle{
		FetchedAt: time.Now().UTC(),
		Source:    "browser-api",
		Info:      json.RawMessage(infoRaw),
		CoinList:  json.RawMessage(coinsRaw),
	}
	return json.Marshal(bundle)
}

// ParseSummaryBundle decodes Bybit-style info + optional coin list payloads.
func ParseSummaryBundle(infoRaw, coinsRaw []byte) (ParsedSummary, error) {
	payload, err := unwrapAPIPayload(infoRaw)
	if err != nil {
		return ParsedSummary{}, err
	}

	var info infoData
	if err := json.Unmarshal(payload, &info); err != nil {
		return ParsedSummary{}, fmt.Errorf("decode info data: %w", err)
	}

	auditID := firstNonEmpty(info.AuditID, info.AuditIDCamel)
	auditTimeRaw := firstNonEmpty(info.AuditTime, info.AuditTimeCamel, info.SnapshotTime)
	merkleRoot := firstNonEmpty(info.MerkleRoot, info.MerkleRootAlt, info.MerkleRootHash)
	totalRate := info.TotalReserveRate
	if totalRate == 0 {
		totalRate = info.TotalReserveAlt
	}
	if totalRate == 0 {
		totalRate = info.ReserveRate
	}

	out := ParsedSummary{
		AuditID:          auditID,
		AuditTimeRaw:     auditTimeRaw,
		MerkleRoot:       merkleRoot,
		TotalReserveRate: totalRate,
	}
	if t, err := parseAuditTime(auditTimeRaw); err == nil {
		out.AuditTime = t
	}

	if len(coinsRaw) > 0 {
		coinPayload, err := unwrapAPIPayload(coinsRaw)
		if err == nil {
			appendCoinRows(&out, coinPayload)
		}
	}

	if out.AuditID == "" && !out.AuditTime.IsZero() {
		out.AuditID = out.AuditTime.Format("20060102")
	}
	if out.MerkleRoot == "" && out.TotalReserveRate == 0 && len(out.CoinRows) == 0 {
		return ParsedSummary{}, fmt.Errorf("bybit summary missing merkle root and reserve fields")
	}
	return out, nil
}

func appendCoinRows(out *ParsedSummary, payload []byte) {
	var coins coinListData
	if err := json.Unmarshal(payload, &coins); err != nil {
		return
	}
	for _, c := range coins.List {
		coin := firstNonEmpty(c.Coin, c.Asset)
		rate := c.ReserveRate
		if rate == 0 {
			rate = c.ReserveRateCamel
		}
		amt := c.ReserveAmount
		if amt == 0 {
			amt = c.ReserveAmountAlt
		}
		liab := c.CustomerLiability
		if liab == 0 {
			liab = c.LiabilityAmount
		}
		out.CoinRows = append(out.CoinRows, CoinRow{
			Coin:            coin,
			ReserveRate:     rate,
			ReserveAmount:   amt,
			LiabilityAmount: liab,
		})
	}
	for _, c := range coins.Coins {
		out.CoinRows = append(out.CoinRows, CoinRow{
			Coin:            c.Coin,
			ReserveRate:     c.ReserveRate,
			ReserveAmount:   c.ReserveAmount,
			LiabilityAmount: c.LiabilityAmount,
		})
	}
}

func unwrapAPIPayload(raw []byte) ([]byte, error) {
	var gate struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(raw, &gate); err == nil && (gate.Code == 0 || gate.Code == 200) && len(gate.Data) > 0 {
		return gate.Data, nil
	}
	var bybit struct {
		RetCode int             `json:"retCode"`
		RetMsg  string          `json:"retMsg"`
		Result  json.RawMessage `json:"result"`
	}
	if err := json.Unmarshal(raw, &bybit); err == nil && bybit.RetCode == 0 && len(bybit.Result) > 0 {
		return bybit.Result, nil
	}
	if gate.Code != 0 && gate.Code != 200 {
		return nil, fmt.Errorf("api code=%d msg=%q (capture JSON in browser DevTools and use -info-file)", gate.Code, gate.Message)
	}
	if bybit.RetCode != 0 {
		return nil, fmt.Errorf("api retCode=%d retMsg=%q", bybit.RetCode, bybit.RetMsg)
	}
	return raw, nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func parseAuditTime(s string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
		"2006010215",
	}
	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			return t.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognized audit time %q", s)
}

func get(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", URLPoRPage)
	req.Header.Set("Origin", "https://www.bybit.com")

	resp, err := httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("GET %s: HTTP %d: %s", url, resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(body) > 0 && body[0] == '<' {
		return nil, fmt.Errorf("GET %s: WAF blocked (HTML response); capture API JSON in browser", url)
	}
	return body, nil
}

func httpClient() *http.Client { return &http.Client{Timeout: 2 * time.Minute} }

// SummaryURL returns the primary public summary endpoint (may be WAF-gated).
func SummaryURL() string { return URLReserveRatio }
