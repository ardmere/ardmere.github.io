// Package gateapi is a thin client for Gate.com public Proof-of-Reserves web APIs.
// See docs/gate-por-data-guide.md for data availability and verification boundaries.
package gateapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	userAgent = "ardmere/0.1 (+https://ardmere.org)"

	// Public endpoints used by https://www.gate.com/proof-of-reserves (Akamai may block datacenter IPs).
	URLProofInfo  = "https://www.gate.com/api/web/v1/proof-of-reserves/getProofOfReservesInfo"
	URLProofList  = "https://www.gate.com/api/web/v1/proof-of-reserves/getProofOfReservesList"
	URLProofCoins = "https://www.gate.com/api/web/v1/proof-of-reserves/getProofOfReservesCoinList"
)

// SummaryBundle is the combined public summary we archive as one artifact.
type SummaryBundle struct {
	FetchedAt time.Time       `json:"fetchedAt"`
	Info      json.RawMessage `json:"info"`
	CoinList  json.RawMessage `json:"coinList,omitempty"`
	List      json.RawMessage `json:"list,omitempty"`
}

// ParsedSummary is the normalized fields we extract from Gate's summary API.
type ParsedSummary struct {
	AuditID          string
	AuditTime        time.Time
	AuditTimeRaw     string
	MerkleRoot       string
	TotalReserveRate float64 // e.g. 115.69 means 115.69%
	TotalReserveUSD  float64
	CustomerNetUSD   float64
	ExcessReserveUSD float64
	CoinRows         []CoinRow
}

type CoinRow struct {
	Coin         string
	ReserveRate  float64
	ReserveAmount float64
	LiabilityAmount float64
}

type webResp struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// infoData matches the current Gate PoR dashboard payload (field names may evolve).
type infoData struct {
	AuditID            string  `json:"audit_id"`
	AuditTime          string  `json:"audit_time"`
	MerkleRoot         string  `json:"merkle_root_hash"`
	TotalReserveRate   float64 `json:"total_reserve_rate"`
	TotalReserveAmount float64 `json:"total_reserve_amount"`
	CustomerNetBalance float64 `json:"customer_net_balance"`
	ExcessReserveValue float64 `json:"excess_reserve_value"`
}

type coinListData struct {
	List []struct {
		Coin              string  `json:"coin"`
		ReserveRate       float64 `json:"reserve_rate"`
		ReserveAmount     float64 `json:"reserve_amount"`
		CustomerLiability float64 `json:"customer_liability"`
	} `json:"list"`
}

// FetchPublicSummary downloads the public dashboard summary + coin list.
func FetchPublicSummary(ctx context.Context) ([]byte, ParsedSummary, error) {
	infoRaw, err := get(ctx, URLProofInfo)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	coinsRaw, err := get(ctx, URLProofCoins)
	if err != nil {
		coinsRaw = nil // coin list is optional for normalize
	}
	listRaw, _ := get(ctx, URLProofList+"?page=1&page_size=5")

	parsed, err := ParseSummaryBundle(infoRaw, coinsRaw)
	if err != nil {
		return nil, ParsedSummary{}, err
	}

	bundle := SummaryBundle{
		FetchedAt: time.Now().UTC(),
		Info:      json.RawMessage(infoRaw),
		CoinList:  json.RawMessage(coinsRaw),
		List:      json.RawMessage(listRaw),
	}
	out, err := json.Marshal(bundle)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	return out, parsed, nil
}

// ParseSummaryBytes parses a saved SummaryBundle or raw info API response.
func ParseSummaryBytes(raw []byte) (ParsedSummary, error) {
	var bundle SummaryBundle
	if err := json.Unmarshal(raw, &bundle); err == nil && len(bundle.Info) > 0 {
		return ParseSummaryBundle(bundle.Info, bundle.CoinList)
	}
	return ParseSummaryBundle(raw, nil)
}

func ParseSummaryBundle(infoRaw, coinsRaw []byte) (ParsedSummary, error) {
	var wr webResp
	if err := json.Unmarshal(infoRaw, &wr); err != nil {
		return ParsedSummary{}, fmt.Errorf("decode info wrapper: %w", err)
	}
	if wr.Code != 0 || len(wr.Data) == 0 {
		return ParsedSummary{}, fmt.Errorf("gate info API code=%d msg=%q (if Akamai blocked, capture JSON in browser DevTools and use gate-fetch -info-file)", wr.Code, wr.Message)
	}
	var info infoData
	if err := json.Unmarshal(wr.Data, &info); err != nil {
		return ParsedSummary{}, fmt.Errorf("decode info data: %w", err)
	}
	if info.MerkleRoot == "" {
		return ParsedSummary{}, fmt.Errorf("gate info missing merkle_root_hash")
	}

	out := ParsedSummary{
		AuditID:          info.AuditID,
		AuditTimeRaw:     info.AuditTime,
		MerkleRoot:       info.MerkleRoot,
		TotalReserveRate: info.TotalReserveRate,
		TotalReserveUSD:  info.TotalReserveAmount,
		CustomerNetUSD:   info.CustomerNetBalance,
		ExcessReserveUSD: info.ExcessReserveValue,
	}
	if t, err := parseAuditTime(info.AuditTime); err == nil {
		out.AuditTime = t
	}

	if len(coinsRaw) > 0 {
		var cr webResp
		if err := json.Unmarshal(coinsRaw, &cr); err == nil && cr.Code == 0 {
			var coins coinListData
			if err := json.Unmarshal(cr.Data, &coins); err == nil {
				for _, c := range coins.List {
					out.CoinRows = append(out.CoinRows, CoinRow{
						Coin:            c.Coin,
						ReserveRate:     c.ReserveRate,
						ReserveAmount:   c.ReserveAmount,
						LiabilityAmount: c.CustomerLiability,
					})
				}
			}
		}
	}
	if out.AuditID == "" && !out.AuditTime.IsZero() {
		out.AuditID = out.AuditTime.Format("20060102")
	}
	return out, nil
}

func parseAuditTime(s string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			return t.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognized audit_time %q", s)
}

func get(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", "https://www.gate.com/proof-of-reserves")

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
		return nil, fmt.Errorf("GET %s: Akamai blocked (HTML response); capture API JSON in browser and use gate-fetch -info-file", url)
	}
	return body, nil
}

func httpClient() *http.Client { return &http.Client{Timeout: 2 * time.Minute} }

// SummaryURL returns the primary public summary endpoint.
func SummaryURL() string { return URLProofInfo }
