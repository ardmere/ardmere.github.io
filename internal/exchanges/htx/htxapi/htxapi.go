// Package htxapi parses HTX (Huobi) Proof-of-Reserves public artifacts.
// See docs/htx-por-data-guide.md for data availability.
package htxapi

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

	// HTX does not document a stable public PoR REST API; browser capture or zk bundle import is typical.
	URLMerklePage = "https://www.htx.com/zh-cn/finance/merkle/"
)

// SummaryBundle is the combined summary we archive as one artifact.
type SummaryBundle struct {
	FetchedAt time.Time       `json:"fetchedAt"`
	Source    string          `json:"source,omitempty"`
	Info      json.RawMessage `json:"info"`
	CoinList  json.RawMessage `json:"coinList,omitempty"`
	ZkDerived bool            `json:"zkDerived,omitempty"`
}

// ParsedSummary is normalized fields extracted from HTX summary JSON.
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

type webResp struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type infoData struct {
	AuditID          string  `json:"audit_id"`
	AuditTime        string  `json:"audit_time"`
	MerkleRoot       string  `json:"merkle_root_hash"`
	TotalReserveRate float64 `json:"total_reserve_rate"`
}

type coinListData struct {
	List []struct {
		Coin              string  `json:"coin"`
		ReserveRate       float64 `json:"reserve_rate"`
		ReserveAmount     float64 `json:"reserve_amount"`
		CustomerLiability float64 `json:"customer_liability"`
	} `json:"list"`
}

// FetchPublicSummary attempts HTX public endpoints; today this usually fails without browser capture.
func FetchPublicSummary(ctx context.Context) ([]byte, ParsedSummary, error) {
	_ = ctx
	return nil, ParsedSummary{}, fmt.Errorf("no stable public HTX PoR API; import -summary-path from browser DevTools or use -zk-bundle with GitHub public-data.zip")
}

// ParseSummaryBytes parses a saved SummaryBundle or raw info JSON.
func ParseSummaryBytes(raw []byte) (ParsedSummary, error) {
	var bundle SummaryBundle
	if err := json.Unmarshal(raw, &bundle); err == nil && len(bundle.Info) > 0 {
		return ParseSummaryBundle(bundle.Info, bundle.CoinList)
	}
	return ParseSummaryBundle(raw, nil)
}

// ParseSummaryBundle decodes HTX-style info + optional coin list payloads.
func ParseSummaryBundle(infoRaw, coinsRaw []byte) (ParsedSummary, error) {
	var wr webResp
	if err := json.Unmarshal(infoRaw, &wr); err != nil {
		return ParsedSummary{}, fmt.Errorf("decode info wrapper: %w", err)
	}
	if wr.Code != 0 && wr.Code != 200 {
		return ParsedSummary{}, fmt.Errorf("info api code=%d msg=%s", wr.Code, wr.Message)
	}
	payload := infoRaw
	if len(wr.Data) > 0 {
		payload = wr.Data
	}

	var info infoData
	if err := json.Unmarshal(payload, &info); err != nil {
		return ParsedSummary{}, fmt.Errorf("decode info data: %w", err)
	}

	out := ParsedSummary{
		AuditID:          info.AuditID,
		AuditTimeRaw:     info.AuditTime,
		MerkleRoot:       info.MerkleRoot,
		TotalReserveRate: info.TotalReserveRate,
	}
	if info.AuditTime != "" {
		for _, layout := range []string{
			"2006-01-02 15:04:05",
			time.RFC3339,
			"2006-01-02",
		} {
			if t, err := time.Parse(layout, info.AuditTime); err == nil {
				out.AuditTime = t.UTC()
				break
			}
		}
	}

	if len(coinsRaw) > 0 {
		var cwr webResp
		coinsPayload := coinsRaw
		if err := json.Unmarshal(coinsRaw, &cwr); err == nil && len(cwr.Data) > 0 {
			coinsPayload = cwr.Data
		}
		var coins coinListData
		if err := json.Unmarshal(coinsPayload, &coins); err == nil {
			for _, row := range coins.List {
				out.CoinRows = append(out.CoinRows, CoinRow{
					Coin:            row.Coin,
					ReserveRate:     row.ReserveRate,
					ReserveAmount:   row.ReserveAmount,
					LiabilityAmount: row.CustomerLiability,
				})
			}
		}
	}
	return out, nil
}

// BuildSummaryBundle marshals info + coinList into a SummaryBundle.
func BuildSummaryBundle(infoRaw, coinsRaw json.RawMessage, source string, zkDerived bool) ([]byte, error) {
	b := SummaryBundle{
		FetchedAt: time.Now().UTC(),
		Source:    source,
		Info:      infoRaw,
		CoinList:  coinsRaw,
		ZkDerived: zkDerived,
	}
	return json.Marshal(b)
}

// SummaryURL returns the human-facing PoR page.
func SummaryURL() string { return URLMerklePage }

func get(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
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
		return nil, fmt.Errorf("http %d: %s", resp.StatusCode, trim(string(body)))
	}
	return body, nil
}

func trim(s string) string {
	if len(s) > 200 {
		return s[:200] + "..."
	}
	return s
}
