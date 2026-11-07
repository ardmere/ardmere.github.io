// Package ourbitapi fetches and parses Ourbit Proof-of-Reserves summary artifacts.
// See docs/ourbit-por-data-guide.md for data availability.
package ourbitapi

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

	URLPoRPage       = "https://www.ourbit.com/assets/reserve-proof"
	URLReserveRatio  = "https://www.ourbit.com/assets/reserve-ratio"
	URLWalletPage    = "https://www.ourbit.com/assets/reserve-proof"
	URLMerklePage    = "https://www.ourbit.com/assets/merkle"
	URLStockInfo     = "https://www.ourbit.com/api/assetbussiness/merkle/stock_info"
	URLWalletDetail  = "https://www.ourbit.com/api/assetbussiness/merkle/wallet_detail"
	URLCurrencies    = "https://www.ourbit.com/api/assetbussiness/merkle/currencies"
	URLMerkleValidate = "https://github.com/OurbitOfficial/merkle_validate"
)

// SummaryBundle is the combined summary we archive as one artifact.
type SummaryBundle struct {
	FetchedAt   time.Time       `json:"fetchedAt"`
	Source      string          `json:"source,omitempty"`
	Info        json.RawMessage `json:"info"`
	CoinList    json.RawMessage `json:"coinList,omitempty"`
	Wallets     json.RawMessage `json:"wallets,omitempty"`
	Currencies  json.RawMessage `json:"currencies,omitempty"`
}

// ParsedSummary is normalized fields extracted from Ourbit summary JSON.
// Reserve rates are stored as percentages (100 = 100%) for solvency-claim.
type ParsedSummary struct {
	AuditID          string
	AuditTime        time.Time
	AuditTimeRaw     string
	MerkleRoot       string
	TotalReserveRate float64
	CoinRows         []CoinRow
	WalletRows       []WalletRow
}

type CoinRow struct {
	Coin            string
	ReserveRate     float64
	ReserveAmount   float64
	LiabilityAmount float64
}

type WalletRow struct {
	Currency    string
	Chain       string
	Address     string
	Amount      float64
	BlockHeight int64
	BalanceType int
}

type apiEnvelope struct {
	Code      any             `json:"code"`
	Msg       string          `json:"msg"`
	Message   string          `json:"message"`
	Timestamp int64           `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

type stockInfoData struct {
	SnapshotDate flexInt64      `json:"snapshotDate"`
	Details      []stockCoinJSON `json:"details"`
}

type stockCoinJSON struct {
	Currency         string      `json:"currency"`
	WalletTotalAsset flexFloat64 `json:"walletTotalAsset"`
	UserSpotAsset    flexFloat64 `json:"userSpotAsset"`
	StockRate        flexFloat64 `json:"stockRate"`
}

type walletDetailData struct {
	Total  int              `json:"total"`
	Result []walletRowJSON  `json:"result"`
}

type walletRowJSON struct {
	Currency    string      `json:"currency"`
	Chain       string      `json:"chain"`
	Address     string      `json:"address"`
	Amount      flexFloat64 `json:"amount"`
	Height      flexInt64   `json:"height"`
	BalanceType int         `json:"balanceType"`
}

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

// FetchPublicSummary downloads stock_info + wallet_detail + currencies without login.
func FetchPublicSummary(ctx context.Context) ([]byte, ParsedSummary, error) {
	stockRaw, err := getJSON(ctx, URLStockInfo)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	walletRaw, err := getJSON(ctx, URLWalletDetail)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	currRaw, err := getJSON(ctx, URLCurrencies)
	if err != nil {
		// currencies is optional enrichment
		currRaw = []byte(`{"code":0,"data":[]}`)
	}

	stockPayload, err := unwrapAPIPayload(stockRaw)
	if err != nil {
		return nil, ParsedSummary{}, fmt.Errorf("stock_info: %w", err)
	}
	walletPayload, err := unwrapAPIPayload(walletRaw)
	if err != nil {
		return nil, ParsedSummary{}, fmt.Errorf("wallet_detail: %w", err)
	}
	currPayload, _ := unwrapAPIPayload(currRaw)

	var stock stockInfoData
	if err := json.Unmarshal(stockPayload, &stock); err != nil {
		return nil, ParsedSummary{}, fmt.Errorf("decode stock_info: %w", err)
	}
	parsed, err := parseStock(stock)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	var wallets walletDetailData
	if err := json.Unmarshal(walletPayload, &wallets); err == nil {
		parsed.WalletRows = parseWallets(wallets)
	}

	detailsRaw, _ := json.Marshal(stock.Details)
	bundle := SummaryBundle{
		FetchedAt:  time.Now().UTC(),
		Source:     "assetbussiness-merkle-public",
		Info:       stockRaw,
		CoinList:   detailsRaw,
		Wallets:    walletPayload,
		Currencies: currPayload,
	}
	out, err := json.Marshal(bundle)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	return out, parsed, nil
}

// ParseSummaryBytes parses a saved SummaryBundle or raw stock_info response.
func ParseSummaryBytes(raw []byte) (ParsedSummary, error) {
	var bundle SummaryBundle
	if err := json.Unmarshal(raw, &bundle); err == nil && len(bundle.Info) > 0 {
		parsed, err := parseInfoRaw(bundle.Info, bundle.CoinList)
		if err != nil {
			return ParsedSummary{}, err
		}
		if len(bundle.Wallets) > 0 {
			var wallets walletDetailData
			if err := json.Unmarshal(bundle.Wallets, &wallets); err == nil {
				parsed.WalletRows = parseWallets(wallets)
			}
		}
		return parsed, nil
	}
	return parseInfoRaw(raw, nil)
}

func parseInfoRaw(infoRaw, coinList json.RawMessage) (ParsedSummary, error) {
	payload, err := unwrapAPIPayload(infoRaw)
	if err != nil {
		// maybe already unwrapped stockInfoData
		payload = infoRaw
	}
	var stock stockInfoData
	if err := json.Unmarshal(payload, &stock); err != nil {
		return ParsedSummary{}, fmt.Errorf("decode ourbit stock info: %w", err)
	}
	if len(stock.Details) == 0 && len(coinList) > 0 {
		var details []stockCoinJSON
		if err := json.Unmarshal(coinList, &details); err == nil {
			stock.Details = details
		}
	}
	return parseStock(stock)
}

func parseStock(stock stockInfoData) (ParsedSummary, error) {
	if len(stock.Details) == 0 {
		return ParsedSummary{}, fmt.Errorf("ourbit stock_info missing details")
	}
	ms := int64(stock.SnapshotDate)
	if ms <= 0 {
		return ParsedSummary{}, fmt.Errorf("ourbit stock_info missing snapshotDate")
	}
	// API uses epoch millis; exchange PoR UI is UTC+8 calendar day.
	auditTime := time.UnixMilli(ms).UTC()
	exchangeDay := auditTime.In(time.FixedZone("UTC+8", 8*3600))
	auditID := exchangeDay.Format("20060102")

	out := ParsedSummary{
		AuditID:      auditID,
		AuditTime:    auditTime,
		AuditTimeRaw: auditTime.Format(time.RFC3339),
	}
	minRate := 0.0
	for i, d := range stock.Details {
		ratePct := float64(d.StockRate) * 100 // API fraction → percentage units
		row := CoinRow{
			Coin:            strings.ToUpper(strings.TrimSpace(d.Currency)),
			ReserveRate:     ratePct,
			ReserveAmount:   float64(d.WalletTotalAsset),
			LiabilityAmount: float64(d.UserSpotAsset),
		}
		out.CoinRows = append(out.CoinRows, row)
		if i == 0 || ratePct < minRate {
			minRate = ratePct
		}
	}
	// stock_info has no aggregate total; use min per-coin rate as conservative floor.
	out.TotalReserveRate = minRate
	return out, nil
}

func parseWallets(wallets walletDetailData) []WalletRow {
	out := make([]WalletRow, 0, len(wallets.Result))
	for _, r := range wallets.Result {
		out = append(out, WalletRow{
			Currency:    strings.ToUpper(strings.TrimSpace(r.Currency)),
			Chain:       r.Chain,
			Address:     r.Address,
			Amount:      float64(r.Amount),
			BlockHeight: int64(r.Height),
			BalanceType: r.BalanceType,
		})
	}
	return out
}

func unwrapAPIPayload(raw []byte) (json.RawMessage, error) {
	var env apiEnvelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, err
	}
	if len(env.Data) == 0 {
		return raw, nil
	}
	if codeFail(env.Code) {
		msg := env.Msg
		if msg == "" {
			msg = env.Message
		}
		return nil, fmt.Errorf("ourbit api code=%v msg=%s", env.Code, msg)
	}
	return env.Data, nil
}

func codeFail(code any) bool {
	switch v := code.(type) {
	case nil:
		return false
	case float64:
		return int(v) != 0
	case int:
		return v != 0
	case string:
		return v != "" && v != "0" && !strings.EqualFold(v, "success")
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
	resp, err := httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d from %s: %s", resp.StatusCode, url, truncate(string(body), 160))
	}
	return body, nil
}

func httpClient() *http.Client {
	return &http.Client{Timeout: 45 * time.Second}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

// BuildSummaryBundle merges browser-captured stock_info (+ optional wallet/currencies files).
func BuildSummaryBundle(infoPath, walletsPath, currenciesPath string) ([]byte, error) {
	infoRaw, err := os.ReadFile(infoPath)
	if err != nil {
		return nil, err
	}
	var walletsRaw, currRaw []byte
	if walletsPath != "" {
		walletsRaw, err = os.ReadFile(walletsPath)
		if err != nil {
			return nil, err
		}
	}
	if currenciesPath != "" {
		currRaw, err = os.ReadFile(currenciesPath)
		if err != nil {
			return nil, err
		}
	}

	stockPayload, err := unwrapAPIPayload(infoRaw)
	if err != nil {
		stockPayload = infoRaw
	}
	var stock stockInfoData
	_ = json.Unmarshal(stockPayload, &stock)
	detailsRaw, _ := json.Marshal(stock.Details)

	var walletPayload json.RawMessage
	if len(walletsRaw) > 0 {
		if p, err := unwrapAPIPayload(walletsRaw); err == nil {
			walletPayload = p
		} else {
			walletPayload = walletsRaw
		}
	}
	var currPayload json.RawMessage
	if len(currRaw) > 0 {
		if p, err := unwrapAPIPayload(currRaw); err == nil {
			currPayload = p
		} else {
			currPayload = currRaw
		}
	}

	bundle := SummaryBundle{
		FetchedAt:  time.Now().UTC(),
		Source:     "browser-api",
		Info:       infoRaw,
		CoinList:   detailsRaw,
		Wallets:    walletPayload,
		Currencies: currPayload,
	}
	return json.Marshal(bundle)
}

// SummaryURL returns the primary public summary endpoint.
func SummaryURL() string { return URLStockInfo }

// WalletDetailURL returns the public multi-chain wallet table endpoint.
func WalletDetailURL() string { return URLWalletDetail }
