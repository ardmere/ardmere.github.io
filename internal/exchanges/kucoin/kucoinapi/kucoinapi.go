// Package kucoinapi parses KuCoin Proof-of-Reserves summary artifacts.
// See docs/kucoin-por-data-guide.md for data availability.
package kucoinapi

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

	URLPoRPage = "https://www.kucoin.com/proof-of-reserves"
)

// SummaryBundle is the combined summary we archive as one artifact.
type SummaryBundle struct {
	FetchedAt time.Time       `json:"fetchedAt"`
	Source    string          `json:"source,omitempty"`
	Info      json.RawMessage `json:"info"`
	CoinList  json.RawMessage `json:"coinList,omitempty"`
}

// ParsedSummary is normalized fields extracted from KuCoin summary JSON.
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
	LatestAuditDate  int64   `json:"latestAuditDate"`
	MerkleRoot       string  `json:"merkle_root_hash"`
	MerkleRootAlt    string  `json:"merkleRoot"`
	MerkleRootHash   string  `json:"merkleRootHash"`
	TotalReserveRate float64 `json:"total_reserve_rate"`
	TotalReserveAlt  float64 `json:"totalReserveRate"`
	ReserveRate      float64 `json:"reserveRate"`
}

type coinListData struct {
	List         []coinJSON `json:"list"`
	Rows         []coinJSON `json:"rows"`
	Data         []coinJSON `json:"data"`
	Coins        []coinJSON `json:"coins"`
	ReserveAsset []coinJSON `json:"reserveAsset"`
}

type coinJSON struct {
	Coin              string      `json:"coin"`
	CoinName          string      `json:"coinName"`
	Asset             string      `json:"asset"`
	Currency          string      `json:"currency"`
	ReserveRate       flexFloat64 `json:"reserve_rate"`
	ReserveRateCamel  flexFloat64 `json:"reserveRate"`
	ReserveAmount     flexFloat64 `json:"reserve_amount"`
	ReserveAmountAlt  flexFloat64 `json:"reserveAmount"`
	PlatformAssets    flexFloat64 `json:"platformAssets"`
	WalletAsset       flexFloat64 `json:"walletAsset"`
	UserAssets        flexFloat64 `json:"userAssets"`
	UserAsset         flexFloat64 `json:"userAsset"`
	CustomerLiability flexFloat64 `json:"customer_liability"`
	LiabilityAmount   flexFloat64 `json:"liabilityAmount"`
}

// flexFloat64 accepts JSON numbers or numeric strings.
type flexFloat64 float64

func (f *flexFloat64) UnmarshalJSON(raw []byte) error {
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
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		*f = flexFloat64(v)
		return nil
	}
	var v float64
	if err := json.Unmarshal(raw, &v); err != nil {
		return err
	}
	*f = flexFloat64(v)
	return nil
}

var nextDataRE = regexp.MustCompile(`(?s)<script id="__NEXT_DATA__"[^>]*>(.*?)</script>`)

// FetchPublicSummary scrapes the public PoR page SSR payload (__NEXT_DATA__).
func FetchPublicSummary(ctx context.Context) ([]byte, ParsedSummary, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URLPoRPage, nil)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/json")

	resp, err := httpClient().Do(req)
	if err != nil {
		return nil, ParsedSummary{}, fmt.Errorf("fetch kucoin por page: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, ParsedSummary{}, fmt.Errorf("GET %s: HTTP %d: %s", URLPoRPage, resp.StatusCode, string(body))
	}
	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	return ParsePoRPageHTML(html)
}

// ParsePoRPageHTML extracts initialReserve from KuCoin's __NEXT_DATA__ blob.
func ParsePoRPageHTML(html []byte) ([]byte, ParsedSummary, error) {
	m := nextDataRE.FindSubmatch(html)
	if m == nil {
		return nil, ParsedSummary{}, fmt.Errorf("kucoin por page missing __NEXT_DATA__; capture JSON manually and use -summary-path")
	}
	var next struct {
		Props struct {
			PageProps struct {
				InitialReserve json.RawMessage `json:"initialReserve"`
			} `json:"pageProps"`
		} `json:"props"`
	}
	if err := json.Unmarshal(m[1], &next); err != nil {
		return nil, ParsedSummary{}, fmt.Errorf("decode __NEXT_DATA__: %w", err)
	}
	if len(next.Props.PageProps.InitialReserve) == 0 {
		return nil, ParsedSummary{}, fmt.Errorf("kucoin __NEXT_DATA__ missing initialReserve")
	}

	parsed, err := ParseSummaryBundle(next.Props.PageProps.InitialReserve, next.Props.PageProps.InitialReserve)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	bundle := SummaryBundle{
		FetchedAt: time.Now().UTC(),
		Source:    "por-page-ssr",
		Info:      next.Props.PageProps.InitialReserve,
		CoinList:  next.Props.PageProps.InitialReserve,
	}
	out, err := json.Marshal(bundle)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	return out, parsed, nil
}

// ParseSummaryBytes parses a saved SummaryBundle or raw API/SSR response.
func ParseSummaryBytes(raw []byte) (ParsedSummary, error) {
	var bundle SummaryBundle
	if err := json.Unmarshal(raw, &bundle); err == nil && len(bundle.Info) > 0 {
		return ParseSummaryBundle(bundle.Info, bundle.CoinList)
	}
	return ParseSummaryBundle(raw, raw)
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
	} else {
		coinsRaw = infoRaw
	}
	bundle := SummaryBundle{
		FetchedAt: time.Now().UTC(),
		Source:    "browser-api",
		Info:      json.RawMessage(infoRaw),
		CoinList:  json.RawMessage(coinsRaw),
	}
	return json.Marshal(bundle)
}

// ParseSummaryBundle decodes KuCoin-style info + coin list payloads.
func ParseSummaryBundle(infoRaw, coinsRaw []byte) (ParsedSummary, error) {
	payload, err := unwrapAPIPayload(infoRaw)
	if err != nil {
		return ParsedSummary{}, err
	}

	var info infoData
	if err := json.Unmarshal(payload, &info); err != nil {
		return ParsedSummary{}, fmt.Errorf("decode info data: %w", err)
	}

	out := ParsedSummary{
		AuditID:          firstNonEmpty(info.AuditID, info.AuditIDCamel),
		AuditTimeRaw:     firstNonEmpty(info.AuditTime, info.AuditTimeCamel, info.SnapshotTime),
		MerkleRoot:       firstNonEmpty(info.MerkleRoot, info.MerkleRootAlt, info.MerkleRootHash),
		TotalReserveRate: firstNonZero(info.TotalReserveRate, info.TotalReserveAlt, info.ReserveRate),
	}
	if info.LatestAuditDate > 0 {
		ms := info.LatestAuditDate
		if ms > 1e12 {
			ms = ms / 1000
		}
		out.AuditTime = time.Unix(ms, 0).UTC()
		if out.AuditTimeRaw == "" {
			out.AuditTimeRaw = out.AuditTime.Format("2006-01-02 15:04:05")
		}
		if out.AuditID == "" {
			out.AuditID = out.AuditTime.Format("20060102")
		}
	}
	if t, err := parseAuditTime(out.AuditTimeRaw); err == nil && out.AuditTime.IsZero() {
		out.AuditTime = t
	}

	coinPayload := coinsRaw
	if len(coinPayload) == 0 {
		coinPayload = payload
	}
	if len(coinPayload) > 0 {
		if unwrapped, err := unwrapAPIPayload(coinPayload); err == nil {
			appendCoinRows(&out, unwrapped)
		}
	}

	if out.TotalReserveRate == 0 && len(out.CoinRows) > 0 {
		minRate := 0.0
		for _, c := range out.CoinRows {
			if c.ReserveRate == 0 {
				continue
			}
			if minRate == 0 || c.ReserveRate < minRate {
				minRate = c.ReserveRate
			}
		}
		out.TotalReserveRate = minRate
	}

	if out.AuditID == "" && !out.AuditTime.IsZero() {
		out.AuditID = out.AuditTime.Format("20060102")
	}
	if out.MerkleRoot == "" && out.TotalReserveRate == 0 && len(out.CoinRows) == 0 {
		return ParsedSummary{}, fmt.Errorf("kucoin summary missing reserve fields")
	}
	return out, nil
}

func appendCoinRows(out *ParsedSummary, payload []byte) {
	var coins coinListData
	if err := json.Unmarshal(payload, &coins); err != nil {
		return
	}
	rows := append(append(append(append(coins.List, coins.Rows...), coins.Data...), coins.Coins...), coins.ReserveAsset...)
	for _, row := range rows {
		coin := firstNonEmpty(row.Coin, row.CoinName, row.Asset, row.Currency)
		if coin == "" {
			continue
		}
		out.CoinRows = append(out.CoinRows, CoinRow{
			Coin:            coin,
			ReserveRate:     float64(firstFlex(row.ReserveRate, row.ReserveRateCamel)),
			ReserveAmount:   float64(firstFlex(row.ReserveAmount, row.ReserveAmountAlt, row.PlatformAssets, row.WalletAsset)),
			LiabilityAmount: float64(firstFlex(row.CustomerLiability, row.LiabilityAmount, row.UserAssets, row.UserAsset)),
		})
	}
}

func unwrapAPIPayload(raw []byte) ([]byte, error) {
	var wrapped struct {
		Code    any             `json:"code"`
		Success *bool           `json:"success"`
		Message string          `json:"message"`
		Msg     string          `json:"msg"`
		Data    json.RawMessage `json:"data"`
		Result  json.RawMessage `json:"result"`
	}
	if err := json.Unmarshal(raw, &wrapped); err == nil {
		codeOK := false
		switch c := wrapped.Code.(type) {
		case float64:
			codeOK = c == 0 || c == 200
		case string:
			codeOK = c == "0" || c == "200" || c == "200000"
		case nil:
			codeOK = wrapped.Success != nil && *wrapped.Success
		}
		if codeOK && len(wrapped.Data) > 0 {
			return wrapped.Data, nil
		}
		if codeOK && len(wrapped.Result) > 0 {
			return wrapped.Result, nil
		}
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

func firstNonZero(vals ...float64) float64 {
	for _, v := range vals {
		if v != 0 {
			return v
		}
	}
	return 0
}

func firstFlex(vals ...flexFloat64) flexFloat64 {
	for _, v := range vals {
		if v != 0 {
			return v
		}
	}
	return 0
}

func parseAuditTime(s string) (time.Time, error) {
	for _, layout := range []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02", "20060102"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognized audit time %q", s)
}

func httpClient() *http.Client { return &http.Client{Timeout: 2 * time.Minute} }

// SummaryURL returns the public PoR page (SSR source for reserve ratios).
func SummaryURL() string { return URLPoRPage }
