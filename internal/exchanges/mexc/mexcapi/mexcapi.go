// Package mexcapi parses MEXC Proof-of-Reserves summary artifacts.
// See docs/mexc-por-data-guide.md for data availability.
package mexcapi

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

	URLPoRSupport = "https://www.mexc.com/support/article/how-to-verify-user-assets-on-mexc-17827791528836"
	URLBlogPoR    = "https://blog.mexc.com/press-release/mexcs-july-2026-proof-of-reserves-report/"
	URLVerifier   = "https://github.com/mexcmerkle/mexc-merkle-proof-of-reserves"
)

// SummaryBundle is the combined summary we archive as one artifact.
type SummaryBundle struct {
	FetchedAt time.Time       `json:"fetchedAt"`
	Source    string          `json:"source,omitempty"`
	Info      json.RawMessage `json:"info"`
	CoinList  json.RawMessage `json:"coinList,omitempty"`
}

// ParsedSummary is normalized fields extracted from MEXC summary JSON.
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
	ProofID          string  `json:"proofId"`
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
	List  []coinJSON `json:"list"`
	Rows  []coinJSON `json:"rows"`
	Data  []coinJSON `json:"data"`
	Coins []coinJSON `json:"coins"`
}

type coinJSON struct {
	Coin              string  `json:"coin"`
	CoinName          string  `json:"coinName"`
	Asset             string  `json:"asset"`
	Currency          string  `json:"currency"`
	ReserveRate       float64 `json:"reserve_rate"`
	ReserveRateCamel  float64 `json:"reserveRate"`
	ReserveAmount     float64 `json:"reserve_amount"`
	ReserveAmountAlt  float64 `json:"reserveAmount"`
	PlatformAssets    float64 `json:"platformAssets"`
	WalletAsset       float64 `json:"walletAsset"`
	UserAssets        float64 `json:"userAssets"`
	UserAsset         float64 `json:"userAsset"`
	CustomerLiability float64 `json:"customer_liability"`
	LiabilityAmount   float64 `json:"liabilityAmount"`
}

// FetchPublicSummary documents that MEXC PoR JSON is not stably permissionless from datacenters.
func FetchPublicSummary(ctx context.Context) ([]byte, ParsedSummary, error) {
	// Probe the support article; Akamai often returns 403 from datacenter IPs.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URLPoRSupport, nil)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/json")
	resp, err := httpClient().Do(req)
	if err != nil {
		return nil, ParsedSummary{}, fmt.Errorf("mexc por probe: %w (WAF may block; use -summary-path)", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
	return nil, ParsedSummary{}, fmt.Errorf("no stable public MEXC PoR JSON endpoint (HTTP %d from support page: %s); capture dashboard API in browser or use -summary-path / press-release fixture", resp.StatusCode, truncate(string(body), 120))
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

// ParseSummaryBundle decodes MEXC-style info + optional coin list payloads.
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
		AuditID:          firstNonEmpty(info.AuditID, info.AuditIDCamel, info.ProofID),
		AuditTimeRaw:     firstNonEmpty(info.AuditTime, info.AuditTimeCamel, info.SnapshotTime),
		MerkleRoot:       firstNonEmpty(info.MerkleRoot, info.MerkleRootAlt, info.MerkleRootHash),
		TotalReserveRate: firstNonZero(info.TotalReserveRate, info.TotalReserveAlt, info.ReserveRate),
	}
	if t, err := parseAuditTime(out.AuditTimeRaw); err == nil {
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
		out.AuditID = out.AuditTime.Format("200601")
	}
	if out.MerkleRoot == "" && out.TotalReserveRate == 0 && len(out.CoinRows) == 0 {
		return ParsedSummary{}, fmt.Errorf("mexc summary missing merkle root and reserve fields")
	}
	return out, nil
}

func appendCoinRows(out *ParsedSummary, payload []byte) {
	var coins coinListData
	if err := json.Unmarshal(payload, &coins); err != nil {
		return
	}
	for _, row := range append(append(append(coins.List, coins.Rows...), coins.Data...), coins.Coins...) {
		coin := firstNonEmpty(row.Coin, row.CoinName, row.Asset, row.Currency)
		if coin == "" {
			continue
		}
		out.CoinRows = append(out.CoinRows, CoinRow{
			Coin:            coin,
			ReserveRate:     firstNonZero(row.ReserveRate, row.ReserveRateCamel),
			ReserveAmount:   firstNonZero(row.ReserveAmount, row.ReserveAmountAlt, row.PlatformAssets, row.WalletAsset),
			LiabilityAmount: firstNonZero(row.CustomerLiability, row.LiabilityAmount, row.UserAssets, row.UserAsset),
		})
	}
}

func unwrapAPIPayload(raw []byte) ([]byte, error) {
	var wrapped struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Msg     string          `json:"msg"`
		Data    json.RawMessage `json:"data"`
		Result  json.RawMessage `json:"result"`
	}
	if err := json.Unmarshal(raw, &wrapped); err == nil {
		if (wrapped.Code == 0 || wrapped.Code == 200) && len(wrapped.Data) > 0 {
			return wrapped.Data, nil
		}
		if (wrapped.Code == 0 || wrapped.Code == 200) && len(wrapped.Result) > 0 {
			return wrapped.Result, nil
		}
		if wrapped.Code != 0 && wrapped.Code != 200 {
			return nil, fmt.Errorf("api code=%d msg=%q", wrapped.Code, firstNonEmpty(wrapped.Message, wrapped.Msg))
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

func parseAuditTime(s string) (time.Time, error) {
	for _, layout := range []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02", "200601"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognized audit time %q", s)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func httpClient() *http.Client { return &http.Client{Timeout: 2 * time.Minute} }

// SummaryURL returns the primary public PoR reference page (JSON API is WAF/login gated).
func SummaryURL() string { return URLPoRSupport }
