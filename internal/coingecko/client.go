// Package coingecko fetches historical USD prices for PoR reserve coverage.
//
// Responses are cached under a gitignored directory (default .cache/coingecko).
// The public API is rate-limited (~10–30 calls/min without a key). On HTTP 429
// the client sleeps and retries a few times; set COINGECKO_API_KEY for the
// pro/demo host when available.
package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.coingecko.com/api/v3"
	defaultCache   = ".cache/coingecko"
	userAgent      = "ardmere-por/coingecko-1"
)

// Client reads historical USD prices, with on-disk JSON caching.
type Client struct {
	BaseURL    string
	CacheDir   string
	HTTPClient *http.Client
	APIKey     string
	Now        func() time.Time
}

func New() *Client {
	base := strings.TrimRight(os.Getenv("COINGECKO_API_BASE"), "/")
	if base == "" {
		base = defaultBaseURL
	}
	cache := os.Getenv("COINGECKO_CACHE_DIR")
	if cache == "" {
		cache = defaultCache
	}
	return &Client{
		BaseURL:    base,
		CacheDir:   cache,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		APIKey:     os.Getenv("COINGECKO_API_KEY"),
		Now:        time.Now,
	}
}

// PriceAt returns the USD price of coinID nearest to at (UTC).
// Prefer market_chart/range; fall back to /history for the UTC calendar day.
func (c *Client) PriceAt(ctx context.Context, coinID string, at time.Time) (price float64, asOf time.Time, err error) {
	if coinID == "" {
		return 0, time.Time{}, fmt.Errorf("empty coin id")
	}
	at = at.UTC()
	if p, t, ok := c.loadCache(coinID, at); ok {
		return p, t, nil
	}
	price, asOf, err = c.priceFromRange(ctx, coinID, at)
	if err != nil {
		price, asOf, err = c.priceFromHistory(ctx, coinID, at)
		if err != nil {
			return 0, time.Time{}, err
		}
	}
	_ = c.storeCache(coinID, at, price, asOf)
	return price, asOf, nil
}

func (c *Client) priceFromRange(ctx context.Context, coinID string, at time.Time) (float64, time.Time, error) {
	// ±36h window around snapshotTime.
	from := at.Add(-36 * time.Hour).Unix()
	to := at.Add(36 * time.Hour).Unix()
	path := fmt.Sprintf("/coins/%s/market_chart/range?vs_currency=usd&from=%d&to=%d",
		url.PathEscape(coinID), from, to)
	body, err := c.get(ctx, path)
	if err != nil {
		return 0, time.Time{}, err
	}
	var parsed struct {
		Prices [][]float64 `json:"prices"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return 0, time.Time{}, err
	}
	if len(parsed.Prices) == 0 {
		return 0, time.Time{}, fmt.Errorf("coingecko range empty for %s", coinID)
	}
	targetMs := float64(at.UnixMilli())
	best := parsed.Prices[0]
	bestDelta := math.Abs(best[0] - targetMs)
	for _, row := range parsed.Prices[1:] {
		if len(row) < 2 {
			continue
		}
		d := math.Abs(row[0] - targetMs)
		if d < bestDelta {
			best = row
			bestDelta = d
		}
	}
	asOf := time.UnixMilli(int64(best[0])).UTC()
	return best[1], asOf, nil
}

func (c *Client) priceFromHistory(ctx context.Context, coinID string, at time.Time) (float64, time.Time, error) {
	// CoinGecko /history uses dd-mm-yyyy.
	date := at.Format("02-01-2006")
	path := fmt.Sprintf("/coins/%s/history?date=%s&localization=false", url.PathEscape(coinID), date)
	body, err := c.get(ctx, path)
	if err != nil {
		return 0, time.Time{}, err
	}
	var parsed struct {
		MarketData struct {
			CurrentPrice map[string]float64 `json:"current_price"`
		} `json:"market_data"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return 0, time.Time{}, err
	}
	p, ok := parsed.MarketData.CurrentPrice["usd"]
	if !ok {
		return 0, time.Time{}, fmt.Errorf("coingecko history missing usd for %s on %s", coinID, date)
	}
	// History is daily — report noon UTC of that day as as-of.
	asOf := time.Date(at.Year(), at.Month(), at.Day(), 12, 0, 0, 0, time.UTC)
	return p, asOf, nil
}

func (c *Client) get(ctx context.Context, path string) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt < 4; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+path, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("Accept", "application/json")
		if c.APIKey != "" {
			req.Header.Set("x-cg-pro-api-key", c.APIKey)
			req.Header.Set("x-cg-demo-api-key", c.APIKey)
		}
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(attempt+1) * time.Second)
			continue
		}
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
		resp.Body.Close()
		if readErr != nil {
			lastErr = readErr
			continue
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			lastErr = fmt.Errorf("coingecko rate limited (429)")
			// Honor Retry-After when present.
			wait := time.Duration(attempt+1) * 8 * time.Second
			if ra := resp.Header.Get("Retry-After"); ra != "" {
				if sec, err := strconv.Atoi(ra); err == nil && sec > 0 {
					wait = time.Duration(sec) * time.Second
				}
			}
			time.Sleep(wait)
			continue
		}
		if resp.StatusCode >= 300 {
			return nil, fmt.Errorf("coingecko HTTP %d: %s", resp.StatusCode, truncate(string(body), 200))
		}
		return body, nil
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("coingecko request failed")
	}
	return nil, lastErr
}

func (c *Client) cachePath(coinID string, at time.Time) string {
	// Bucket by hour so range/history lookups share a file nearby.
	key := fmt.Sprintf("%s_%s.json", coinID, at.UTC().Format("2006010215"))
	key = strings.ReplaceAll(key, "/", "_")
	return filepath.Join(c.CacheDir, key)
}

func (c *Client) loadCache(coinID string, at time.Time) (float64, time.Time, bool) {
	path := c.cachePath(coinID, at)
	raw, err := os.ReadFile(path)
	if err != nil {
		return 0, time.Time{}, false
	}
	var doc struct {
		Price float64 `json:"price"`
		AsOf  string  `json:"asOf"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return 0, time.Time{}, false
	}
	asOf, err := time.Parse(time.RFC3339, doc.AsOf)
	if err != nil {
		return 0, time.Time{}, false
	}
	return doc.Price, asOf, true
}

func (c *Client) storeCache(coinID string, at time.Time, price float64, asOf time.Time) error {
	if err := os.MkdirAll(c.CacheDir, 0o755); err != nil {
		return err
	}
	doc := map[string]any{
		"coinId":        coinID,
		"snapshotBucket": at.UTC().Format(time.RFC3339),
		"price":         price,
		"asOf":          asOf.UTC().Format(time.RFC3339),
		"cachedAt":      c.Now().UTC().Format(time.RFC3339),
	}
	raw, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.cachePath(coinID, at), raw, 0o644)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
