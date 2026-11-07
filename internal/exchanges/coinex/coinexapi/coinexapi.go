// Package coinexapi fetches and parses CoinEx Proof-of-Reserves summary artifacts.
// See docs/coinex-por-data-guide.md for data availability.
package coinexapi

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

	URLPoRPage       = "https://www.coinex.com/en/reserve-proof"
	URLReserveAPI    = "https://www.coinex.com/res/account/balance/merkle/reserve"
	URLUserProofAPI  = "https://www.coinex.com/res/account/balance/merkle/proof"
	URLCheckMerkle   = "https://www.coinex.com/en/reserve-proof/help/check-merkle"
	URLVerifyAddress = "https://www.coinex.com/en/reserve-proof/help/verify-address"
	URLOSSVerifier   = "https://github.com/coinexcom/proof-of-reserves"
)

// SummaryBundle is the archived summarySnapshot for one CoinEx PoR period.
type SummaryBundle struct {
	FetchedAt time.Time       `json:"fetchedAt"`
	Source    string          `json:"source,omitempty"`
	Period    json.RawMessage `json:"period"`
}

// ParsedSummary is normalized fields extracted from a CoinEx period.
// Reserve rates are percentages (100 = 100%) for solvency-claim.
type ParsedSummary struct {
	AuditID          string
	AuditTime        time.Time
	AuditTimeRaw     string
	MerkleRoot       string
	TotalReserveRate float64
	WalletTotalUSD   string
	WalletProofFile  string
	CoinRows         []CoinRow
}

type CoinRow struct {
	Coin            string
	ReserveRate     float64
	ReserveAmount   float64
	LiabilityAmount float64
	WalletUSD       string
	Price           string
}

type apiEnvelope struct {
	Code    any             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type periodJSON struct {
	Time            flexInt64     `json:"time"`
	RootHash        string        `json:"root_hash"`
	WalletProofFile string        `json:"wallet_proof_file"`
	WalletTotalUSD  flexString    `json:"wallet_total_usd"`
	Balances        []coinRowJSON `json:"balances"`
}

type coinRowJSON struct {
	Asset        string     `json:"asset"`
	WalletAmount flexString `json:"wallet_amount"`
	UserAmount   flexString `json:"user_amount"`
	ReserveRate  flexString `json:"reserve_rate"`
	WalletUSD    flexString `json:"wallet_usd"`
	Price        flexString `json:"price"`
}

// flexString keeps the original decimal text when present.
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
func SummaryURL() string { return URLReserveAPI }

// FetchPublicSummary downloads the public merkle/reserve list and archives the latest period.
func FetchPublicSummary(ctx context.Context) ([]byte, ParsedSummary, error) {
	raw, err := getJSON(ctx, URLReserveAPI)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	periods, err := unwrapPeriodRaws(raw)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	if len(periods) == 0 {
		return nil, ParsedSummary{}, fmt.Errorf("coinex reserve API returned empty data")
	}
	return buildBundleFromRaw(periods[0])
}

// ParseSummaryBytes parses an archived SummaryBundle or a raw period object.
func ParseSummaryBytes(raw []byte) (ParsedSummary, error) {
	var bundle SummaryBundle
	if err := json.Unmarshal(raw, &bundle); err == nil && len(bundle.Period) > 0 {
		return parsePeriodRaw(bundle.Period)
	}
	return parsePeriodRaw(raw)
}

// BuildSummaryBundle wraps a browser-captured period JSON (or full API envelope) as an archive bundle.
func BuildSummaryBundle(periodFile string) ([]byte, error) {
	raw, err := os.ReadFile(periodFile)
	if err != nil {
		return nil, err
	}
	if periods, err := unwrapPeriodRaws(raw); err == nil && len(periods) > 0 {
		out, _, err := buildBundleFromRaw(periods[0])
		return out, err
	}
	out, _, err := buildBundleFromRaw(json.RawMessage(raw))
	return out, err
}

func buildBundleFromRaw(periodRaw json.RawMessage) ([]byte, ParsedSummary, error) {
	parsed, err := parsePeriodRaw(periodRaw)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	bundle := SummaryBundle{
		FetchedAt: time.Now().UTC(),
		Source:    URLReserveAPI,
		Period:    periodRaw,
	}
	out, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	return out, parsed, nil
}

func parsePeriodRaw(periodRaw json.RawMessage) (ParsedSummary, error) {
	var period periodJSON
	if err := json.Unmarshal(periodRaw, &period); err != nil {
		return ParsedSummary{}, fmt.Errorf("decode coinex period: %w", err)
	}
	return parsePeriod(period)
}

func parsePeriod(period periodJSON) (ParsedSummary, error) {
	if period.Time == 0 {
		return ParsedSummary{}, fmt.Errorf("coinex period missing time")
	}
	ts := time.Unix(int64(period.Time), 0).UTC()
	out := ParsedSummary{
		AuditID:         ts.Format("20060102"),
		AuditTime:       ts,
		AuditTimeRaw:    ts.Format(time.RFC3339),
		MerkleRoot:      strings.ToLower(strings.TrimPrefix(strings.TrimSpace(period.RootHash), "0x")),
		WalletTotalUSD:  string(period.WalletTotalUSD),
		WalletProofFile: strings.TrimSpace(period.WalletProofFile),
	}
	minRate := 0.0
	for i, c := range period.Balances {
		rate := mustFloat(string(c.ReserveRate))
		out.CoinRows = append(out.CoinRows, CoinRow{
			Coin:            strings.ToUpper(strings.TrimSpace(c.Asset)),
			ReserveRate:     rate,
			ReserveAmount:   mustFloat(string(c.WalletAmount)),
			LiabilityAmount: mustFloat(string(c.UserAmount)),
			WalletUSD:       string(c.WalletUSD),
			Price:           string(c.Price),
		})
		if i == 0 || rate < minRate {
			minRate = rate
		}
	}
	out.TotalReserveRate = minRate
	return out, nil
}

func unwrapPeriodRaws(raw []byte) ([]json.RawMessage, error) {
	var env apiEnvelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, fmt.Errorf("decode coinex envelope: %w", err)
	}
	if !codeOK(env.Code) {
		return nil, fmt.Errorf("coinex API code=%v message=%s", env.Code, env.Message)
	}
	var periods []json.RawMessage
	if err := json.Unmarshal(env.Data, &periods); err != nil {
		return nil, fmt.Errorf("decode coinex periods: %w", err)
	}
	return periods, nil
}

func codeOK(code any) bool {
	switch v := code.(type) {
	case float64:
		return v == 0
	case int:
		return v == 0
	case string:
		return v == "0" || v == "" || strings.EqualFold(v, "ok")
	case nil:
		return true
	default:
		return false
	}
}

func getJSON(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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
		return nil, fmt.Errorf("GET %s: HTTP %d: %s", url, resp.StatusCode, truncate(string(body), 200))
	}
	return body, nil
}

func mustFloat(s string) float64 {
	s = strings.TrimSpace(s)
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
