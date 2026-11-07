// Package poloniexapi fetches and parses Poloniex Proof-of-Reserves summary artifacts.
// See docs/poloniex-por-data-guide.md for data availability.
package poloniexapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	userAgent = "ardmere/0.1 (+https://ardmere.org)"

	URLPoRPage         = "https://www.poloniex.com/topic/security/proof-of-reserves"
	URLVersionInfo     = "https://www.poloniex.com/proxy/sapi/merkleProof/public/versionInfo"
	URLSnapshotBalance = "https://www.poloniex.com/proxy/sapi/merkleProof/public/snapshotBalance"
	URLOwnershipOSS    = "https://github.com/poloniex/tools-nodejs-address-verify"
)

// SummaryBundle is the combined summary we archive as one artifact.
type SummaryBundle struct {
	FetchedAt time.Time       `json:"fetchedAt"`
	Source    string          `json:"source,omitempty"`
	Version   string          `json:"version,omitempty"`
	Info      json.RawMessage `json:"info"`
	CoinList  json.RawMessage `json:"coinList,omitempty"`
}

// ParsedSummary is normalized fields extracted from Poloniex summary JSON.
// Reserve rates are stored as percentages (100 = 100%) for solvency-claim.
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
	CurrencyGroup   string
}

type apiEnvelope struct {
	Code    any             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type versionInfoData struct {
	VersionList    []string `json:"versionList"`
	CurrentVersion string   `json:"currentVersion"`
}

type snapshotData struct {
	TotalReserveRate flexFloat64      `json:"totalReserveRate"`
	MerkleHash       string           `json:"merkleHash"`
	MerkleRoot       string           `json:"merkleRoot"`
	MerkleRootHash   string           `json:"merkleRootHash"`
	BalanceList      []balanceRowJSON `json:"balanceList"`
}

type balanceRowJSON struct {
	Currency      string      `json:"currency"`
	CurrencyGroup string      `json:"currencyGroup"`
	UserBalance   flexFloat64 `json:"userBalance"`
	WalletBalance flexFloat64 `json:"walletBalance"`
	ReserveRate   flexFloat64 `json:"reserveRate"`
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
		s = strings.TrimSpace(strings.TrimSuffix(s, "%"))
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

// FetchPublicSummary downloads versionInfo + snapshotBalance for the current (or requested) version.
func FetchPublicSummary(ctx context.Context, version string) ([]byte, ParsedSummary, error) {
	versions, err := fetchVersionInfo(ctx)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	if version == "" {
		version = versions.CurrentVersion
	}
	if version == "" {
		return nil, ParsedSummary{}, fmt.Errorf("poloniex versionInfo missing currentVersion")
	}

	snapRaw, err := getJSON(ctx, URLSnapshotBalance+"?version="+version)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	parsed, err := parseSnapshotPayload(version, snapRaw)
	if err != nil {
		return nil, ParsedSummary{}, err
	}

	versionRaw, _ := json.Marshal(versions)
	bundle := SummaryBundle{
		FetchedAt: time.Now().UTC(),
		Source:    "proxy-sapi-public",
		Version:   version,
		Info:      json.RawMessage(mustMergeInfo(version, versionRaw, snapRaw)),
		CoinList:  json.RawMessage(snapRaw),
	}
	out, err := json.Marshal(bundle)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	return out, parsed, nil
}

func fetchVersionInfo(ctx context.Context) (versionInfoData, error) {
	raw, err := getJSON(ctx, URLVersionInfo)
	if err != nil {
		return versionInfoData{}, err
	}
	payload, err := unwrapAPIPayload(raw)
	if err != nil {
		return versionInfoData{}, err
	}
	var info versionInfoData
	if err := json.Unmarshal(payload, &info); err != nil {
		return versionInfoData{}, fmt.Errorf("decode versionInfo: %w", err)
	}
	return info, nil
}

func mustMergeInfo(version string, versionRaw, snapRaw []byte) []byte {
	info := map[string]any{
		"version": version,
	}
	var versionObj any
	if err := json.Unmarshal(versionRaw, &versionObj); err == nil {
		info["versionInfo"] = versionObj
	}
	payload, err := unwrapAPIPayload(snapRaw)
	if err == nil {
		var snap any
		if json.Unmarshal(payload, &snap) == nil {
			info["snapshot"] = snap
		}
	}
	b, _ := json.Marshal(info)
	return b
}

// ParseSummaryBytes parses a saved SummaryBundle or raw API response.
func ParseSummaryBytes(raw []byte) (ParsedSummary, error) {
	var bundle SummaryBundle
	if err := json.Unmarshal(raw, &bundle); err == nil && len(bundle.Info) > 0 {
		version := bundle.Version
		if version == "" {
			var meta struct {
				Version string `json:"version"`
			}
			_ = json.Unmarshal(bundle.Info, &meta)
			version = meta.Version
		}
		coinRaw := bundle.CoinList
		if len(coinRaw) == 0 {
			coinRaw = bundle.Info
		}
		return parseSnapshotPayload(version, coinRaw)
	}
	return parseSnapshotPayload("", raw)
}

// BuildSummaryBundle merges browser-captured versionInfo + snapshotBalance files.
func BuildSummaryBundle(infoPath, coinsPath, version string) ([]byte, error) {
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
	if version == "" {
		parsed, err := parseSnapshotPayload("", coinsRaw)
		if err == nil {
			version = parsed.AuditID
		}
	}
	bundle := SummaryBundle{
		FetchedAt: time.Now().UTC(),
		Source:    "browser-api",
		Version:   version,
		Info:      json.RawMessage(infoRaw),
		CoinList:  json.RawMessage(coinsRaw),
	}
	return json.Marshal(bundle)
}

func parseSnapshotPayload(version string, raw []byte) (ParsedSummary, error) {
	payload, err := unwrapAPIPayload(raw)
	if err != nil {
		return ParsedSummary{}, err
	}

	// Prefer nested snapshot object from our archived info blob.
	var nested struct {
		Version  string          `json:"version"`
		Snapshot json.RawMessage `json:"snapshot"`
	}
	if err := json.Unmarshal(payload, &nested); err == nil && len(nested.Snapshot) > 0 {
		if version == "" {
			version = nested.Version
		}
		payload = nested.Snapshot
	}

	var snap snapshotData
	if err := json.Unmarshal(payload, &snap); err != nil {
		return ParsedSummary{}, fmt.Errorf("decode snapshotBalance: %w", err)
	}

	root := firstNonEmpty(snap.MerkleHash, snap.MerkleRoot, snap.MerkleRootHash)
	total := asPercent(float64(snap.TotalReserveRate))
	out := ParsedSummary{
		AuditID:          version,
		MerkleRoot:       root,
		TotalReserveRate: total,
	}
	for _, row := range snap.BalanceList {
		coin := strings.TrimSpace(row.Currency)
		if coin == "" {
			continue
		}
		out.CoinRows = append(out.CoinRows, CoinRow{
			Coin:            coin,
			ReserveRate:     asPercent(float64(row.ReserveRate)),
			ReserveAmount:   float64(row.WalletBalance),
			LiabilityAmount: float64(row.UserBalance),
			CurrencyGroup:   row.CurrencyGroup,
		})
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
	if out.AuditID == "" {
		// Try parse from parent bundle version field callers already set; leave empty for normalize fallback.
	}
	if t, err := time.Parse("20060102", out.AuditID); err == nil {
		out.AuditTime = t.UTC()
		out.AuditTimeRaw = t.UTC().Format("2006-01-02") + " (snapshot date; time not published)"
	}
	if out.MerkleRoot == "" && out.TotalReserveRate == 0 && len(out.CoinRows) == 0 {
		return ParsedSummary{}, fmt.Errorf("poloniex summary missing merkle root and reserve fields")
	}
	if out.AuditID == "" {
		return ParsedSummary{}, fmt.Errorf("poloniex summary missing snapshot version id")
	}
	return out, nil
}

// asPercent converts Poloniex API ratios (1.00 = 100%) into percentage units.
// Values already >= 10 are treated as percentages.
func asPercent(v float64) float64 {
	if v == 0 {
		return 0
	}
	if v < 10 {
		return v * 100
	}
	return v
}

func unwrapAPIPayload(raw []byte) ([]byte, error) {
	var wrapped apiEnvelope
	if err := json.Unmarshal(raw, &wrapped); err == nil {
		codeOK := false
		switch c := wrapped.Code.(type) {
		case float64:
			codeOK = c == 0 || c == 200
		case string:
			codeOK = c == "0" || c == "200"
		}
		if wrapped.Code != nil && !codeOK {
			return nil, fmt.Errorf("api code=%v msg=%q", wrapped.Code, wrapped.Message)
		}
		if codeOK && len(wrapped.Data) > 0 {
			return wrapped.Data, nil
		}
	}
	return raw, nil
}

func getJSON(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Origin", "https://www.poloniex.com")
	req.Header.Set("Referer", URLPoRPage)

	resp, err := httpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", url, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: HTTP %d: %s", url, resp.StatusCode, truncate(string(body), 200))
	}
	return body, nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func httpClient() *http.Client { return &http.Client{Timeout: 2 * time.Minute} }

// SummaryURL returns the public PoR landing page.
func SummaryURL() string { return URLPoRPage }

// SnapshotBalanceURL returns the public snapshotBalance endpoint for a version.
func SnapshotBalanceURL(version string) string {
	if version == "" {
		return URLSnapshotBalance
	}
	return URLSnapshotBalance + "?version=" + version
}
