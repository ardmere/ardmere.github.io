// Package bitgetapi parses Bitget Proof-of-Reserves summary artifacts.
// Bitget's public page is browser-oriented; ardmere primarily imports captured JSON.
package bitgetapi

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

	URLPoRPage     = "https://www.bitget.com/proof-of-reserves"
	URLReservePage = "https://www.bitget.com/proof-of-reserves"
)

type SummaryBundle struct {
	FetchedAt time.Time       `json:"fetchedAt"`
	Source    string          `json:"source,omitempty"`
	Info      json.RawMessage `json:"info"`
	CoinList  json.RawMessage `json:"coinList,omitempty"`
}

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
	MerkelRoot       string  `json:"merkelRoot"`
	MerkleRootHash   string  `json:"merkleRootHash"`
	MerkelRootHash   string  `json:"merkelRootHash"`
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
	ReserveRate       float64 `json:"reserve_rate"`
	ReserveRateCamel  float64 `json:"reserveRate"`
	ReserveAmount     float64 `json:"reserve_amount"`
	ReserveAmountAlt  float64 `json:"reserveAmount"`
	PlatformAssets    float64 `json:"platformAssets"`
	UserAssets        float64 `json:"userAssets"`
	CustomerLiability float64 `json:"customer_liability"`
	LiabilityAmount   float64 `json:"liabilityAmount"`
}

func FetchPublicSummary(ctx context.Context) ([]byte, ParsedSummary, error) {
	return nil, ParsedSummary{}, fmt.Errorf("no stable public Bitget PoR JSON endpoint; capture the page API in browser DevTools or use -summary-path")
}

func ParseSummaryBytes(raw []byte) (ParsedSummary, error) {
	var bundle SummaryBundle
	if err := json.Unmarshal(raw, &bundle); err == nil && len(bundle.Info) > 0 {
		return ParseSummaryBundle(bundle.Info, bundle.CoinList)
	}
	return ParseSummaryBundle(raw, nil)
}

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

func ParseSummaryBundle(infoRaw, coinsRaw []byte) (ParsedSummary, error) {
	payload, err := unwrapAPIPayload(infoRaw)
	if err != nil {
		return ParsedSummary{}, err
	}

	var info infoData
	if err := json.Unmarshal(payload, &info); err != nil {
		return ParsedSummary{}, fmt.Errorf("decode info data: %w", err)
	}

	merkleRoot := firstNonEmpty(info.MerkleRoot, info.MerkelRoot, info.MerkleRootAlt, info.MerkleRootHash, info.MerkelRootHash)
	totalRate := firstNonZero(info.TotalReserveRate, info.TotalReserveAlt, info.ReserveRate)
	out := ParsedSummary{
		AuditID:          firstNonEmpty(info.AuditID, info.AuditIDCamel),
		AuditTimeRaw:     firstNonEmpty(info.AuditTime, info.AuditTimeCamel, info.SnapshotTime),
		MerkleRoot:       merkleRoot,
		TotalReserveRate: totalRate,
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

	if out.AuditID == "" && !out.AuditTime.IsZero() {
		out.AuditID = out.AuditTime.Format("20060102")
	}
	if out.MerkleRoot == "" && out.TotalReserveRate == 0 && len(out.CoinRows) == 0 {
		return ParsedSummary{}, fmt.Errorf("bitget summary missing merkle root and reserve fields")
	}
	return out, nil
}

func appendCoinRows(out *ParsedSummary, payload []byte) {
	var coins coinListData
	if err := json.Unmarshal(payload, &coins); err != nil {
		return
	}
	for _, row := range append(append(append(coins.List, coins.Rows...), coins.Data...), coins.Coins...) {
		coin := firstNonEmpty(row.Coin, row.CoinName, row.Asset)
		if coin == "" {
			continue
		}
		reserveAmount := firstNonZero(row.ReserveAmount, row.ReserveAmountAlt, row.PlatformAssets)
		liabilityAmount := firstNonZero(row.CustomerLiability, row.LiabilityAmount, row.UserAssets)
		reserveRate := firstNonZero(row.ReserveRate, row.ReserveRateCamel)
		out.CoinRows = append(out.CoinRows, CoinRow{
			Coin:            coin,
			ReserveRate:     reserveRate,
			ReserveAmount:   reserveAmount,
			LiabilityAmount: liabilityAmount,
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

func get(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", URLPoRPage)

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
	return body, nil
}

func httpClient() *http.Client { return &http.Client{Timeout: 2 * time.Minute} }

func SummaryURL() string { return URLReservePage }
