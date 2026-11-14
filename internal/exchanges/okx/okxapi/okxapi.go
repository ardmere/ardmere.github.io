// Package okxapi fetches OKX Proof-of-Reserves page data and CDN artifacts.
// OKX embeds structured JSON in HTML (id="appState"); there is no public REST API.
package okxapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	URLDetail   = "https://www.okx.com/proof-of-reserves/detail"
	URLDownload = "https://www.okx.com/proof-of-reserves/download"
	userAgent   = "ardmere/0.1 (+https://ardmere.org)"
)

// AuditRootInfo is auditRootInfo from the detail page appState.
type AuditRootInfo struct {
	AuditID                  string            `json:"auditId"`
	MerkleHash               string            `json:"merkleHash"`
	CreateTime               int64             `json:"createTime"`
	TypeNum                  int               `json:"typeNum"`
	CapitalRatio             map[string]string `json:"capitalRatio"`
	LiabilityBalances        map[string]string `json:"liabilityBalances"`
	ReserveBalances          map[string]string `json:"reserveBalances"`
	ExchangeReserveBalances  map[string]string `json:"exchangeReserveBalances"`
	CustodyReserveBalances   map[string]string `json:"custodyReserveBalances"`
	ReserveCurrencies        []string          `json:"reserveCurrencies"`
	LiabilityCurrencies      []string          `json:"liabilityCurrencies"`
}

// AuditRecord is one entry from the download page audit list.
type AuditRecord struct {
	AuditID            string `json:"auditId"`
	CreateTime         int64  `json:"createTime"`
	TypeNum            int    `json:"typeNum"`
	Download           string `json:"download"`
	MerkleTreeDownload string `json:"merkleTreeDownload"`
}

// SummaryBundle is the normalized summary artifact we archive locally.
type SummaryBundle struct {
	FetchedAt string        `json:"fetchedAt"`
	Source    string        `json:"source"`
	DetailURL string        `json:"detailUrl"`
	Audit     AuditRootInfo `json:"audit"`
}

// FetchCurrentSummary loads the latest auditRootInfo from the detail page.
func FetchCurrentSummary(ctx context.Context) (SummaryBundle, []byte, error) {
	html, err := fetchHTML(ctx, URLDetail)
	if err != nil {
		return SummaryBundle{}, nil, err
	}
	audit, err := parseAuditRootInfo(html)
	if err != nil {
		return SummaryBundle{}, nil, err
	}
	bundle := SummaryBundle{
		FetchedAt: time.Now().UTC().Format(time.RFC3339),
		Source:    "page-appState",
		DetailURL: URLDetail,
		Audit:     audit,
	}
	raw, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return SummaryBundle{}, nil, err
	}
	return bundle, raw, nil
}

// FetchAuditList returns historical audit records from the download page.
func FetchAuditList(ctx context.Context) ([]AuditRecord, error) {
	html, err := fetchHTML(ctx, URLDownload)
	if err != nil {
		return nil, err
	}
	return parseAuditList(html)
}

// ParseSummaryBytes decodes a saved summary bundle JSON.
func ParseSummaryBytes(raw []byte) (SummaryBundle, error) {
	var bundle SummaryBundle
	if err := json.Unmarshal(raw, &bundle); err != nil {
		return SummaryBundle{}, fmt.Errorf("decode okx summary: %w", err)
	}
	if bundle.Audit.AuditID == "" {
		return SummaryBundle{}, fmt.Errorf("okx summary missing auditId")
	}
	return bundle, nil
}

// ResolveDownloads picks wallet and liability CDN URLs for an audit id.
func ResolveDownloads(ctx context.Context, auditID string) (walletURL, liabilityURL string, err error) {
	list, err := FetchAuditList(ctx)
	if err != nil {
		return "", "", err
	}
	for _, rec := range list {
		if rec.AuditID == auditID {
			if rec.Download == "" || rec.MerkleTreeDownload == "" {
				return "", "", fmt.Errorf("audit %s missing download URLs", auditID)
			}
			return rec.Download, rec.MerkleTreeDownload, nil
		}
	}
	return "", "", fmt.Errorf("audit %s not found in download list", auditID)
}

func fetchHTML(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GET %s: HTTP %d", url, resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return "", err
	}
	return string(body), nil
}

var appStateRe = regexp.MustCompile(`(?s)<script[^>]*id="appState"[^>]*>(\{.*?\})</script>`)

func parseAppState(html string) (map[string]any, error) {
	m := appStateRe.FindStringSubmatch(html)
	if len(m) < 2 {
		return nil, fmt.Errorf("appState script not found")
	}
	var root map[string]any
	if err := json.Unmarshal([]byte(m[1]), &root); err != nil {
		return nil, fmt.Errorf("decode appState: %w", err)
	}
	return root, nil
}

func parseAuditRootInfo(html string) (AuditRootInfo, error) {
	root, err := parseAppState(html)
	if err != nil {
		return AuditRootInfo{}, err
	}
	raw, err := json.Marshal(dig(root, "appContext", "initialProps", "App", "auditRootInfo"))
	if err != nil {
		return AuditRootInfo{}, err
	}
	var audit AuditRootInfo
	if err := json.Unmarshal(raw, &audit); err != nil {
		return AuditRootInfo{}, fmt.Errorf("decode auditRootInfo: %w", err)
	}
	if audit.AuditID == "" {
		return AuditRootInfo{}, fmt.Errorf("auditRootInfo missing auditId")
	}
	audit.MerkleHash = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(audit.MerkleHash)), "0x")
	return audit, nil
}

func parseAuditList(html string) ([]AuditRecord, error) {
	root, err := parseAppState(html)
	if err != nil {
		return nil, err
	}
	raw, err := json.Marshal(dig(root, "appContext", "initialProps", "App", "auditList"))
	if err != nil {
		return nil, err
	}
	var list []AuditRecord
	if err := json.Unmarshal(raw, &list); err != nil {
		return nil, fmt.Errorf("decode auditList: %w", err)
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("auditList empty")
	}
	return list, nil
}

func dig(m map[string]any, keys ...string) any {
	var cur any = m
	for _, k := range keys {
		obj, ok := cur.(map[string]any)
		if !ok {
			return nil
		}
		cur = obj[k]
	}
	return cur
}
