package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// SolanaSPLBalanceAtSlot returns SPL balance at owner+mint. When a historical
// provider is configured and slot > 0, it queries slot-indexed balance first;
// otherwise it falls back to live getTokenAccountsByOwner.
func (c *LedgerClient) SolanaSPLBalanceAtSlot(ctx context.Context, owner, mint string, slot int64) (amount uint64, used string, historical bool, err error) {
	if slot > 0 {
		if raw, u, ok, herr := c.solanaHistorySPLAtSlot(ctx, owner, mint, slot); ok {
			return raw, u, true, herr
		}
	}
	raw, u, err := c.SolanaSPLBalance(ctx, owner, mint)
	return raw, u, false, err
}

type solanaHistoryConfig struct {
	BaseURL string
	APIKey  string
}

func loadSolanaHistoryConfig() (solanaHistoryConfig, bool) {
	key := strings.TrimSpace(os.Getenv("SOLANA_INDEX_API_KEY"))
	if key == "" {
		key = strings.TrimSpace(os.Getenv("SOLANA_HISTORY_API_KEY"))
	}
	if key == "" {
		return solanaHistoryConfig{}, false
	}
	base := strings.TrimRight(strings.TrimSpace(os.Getenv("SOLANA_HISTORY_API_URL")), "/")
	if base == "" {
		base = "https://solanaindex.top"
	}
	return solanaHistoryConfig{BaseURL: base, APIKey: key}, true
}

func (c *LedgerClient) solanaHistorySPLAtSlot(ctx context.Context, owner, mint string, slot int64) (uint64, string, bool, error) {
	cfg, ok := loadSolanaHistoryConfig()
	if !ok {
		return 0, "", false, nil
	}
	cacheChain := "solana:history:solanaindex"
	cacheKey := owner + "|" + mint
	if raw, ok := c.cacheGet(cacheChain, "balance", cacheKey, slot); ok {
		v, _ := strconv.ParseUint(string(raw), 10, 64)
		return v, "cache", true, nil
	}
	path := fmt.Sprintf("%s/api/v1/solana/token-balance/%s/%s/%d",
		cfg.BaseURL, url.PathEscape(owner), url.PathEscape(mint), slot)

	var body []byte
	var status int
	var reqErr error
	for attempt := 0; attempt < 4; attempt++ {
		c.throttleSolanaHistory()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		if err != nil {
			return 0, "", false, err
		}
		req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
		req.Header.Set("User-Agent", "ardmere/0.1")
		req.Header.Set("Accept", "application/json")

		resp, err := c.httpc.Do(req)
		if err != nil {
			reqErr = fmt.Errorf("solana history: %w", err)
			time.Sleep(time.Duration(attempt+1) * 500 * time.Millisecond)
			continue
		}
		body, _ = io.ReadAll(resp.Body)
		resp.Body.Close()
		status = resp.StatusCode
		if status == http.StatusTooManyRequests {
			time.Sleep(time.Duration(attempt+1) * time.Second)
			continue
		}
		reqErr = nil
		break
	}
	if reqErr != nil {
		return 0, "", false, nil
	}
	if status == http.StatusUnauthorized || status == http.StatusForbidden {
		return 0, cfg.BaseURL, true, fmt.Errorf("solana history HTTP %d: unauthorized (check SOLANA_INDEX_API_KEY)", status)
	}
	if status == http.StatusTooManyRequests {
		return 0, "", false, nil
	}
	if status != http.StatusOK {
		return 0, cfg.BaseURL, true, fmt.Errorf("solana history HTTP %d: %s", status, truncate(string(body), 120))
	}
	raw, err := parseSolanaIndexBalanceRaw(body)
	if err != nil {
		return 0, cfg.BaseURL, true, err
	}
	c.cachePut(cacheChain, "balance", cacheKey, slot, []byte(strconv.FormatUint(raw, 10)))
	return raw, cfg.BaseURL + "/api/v1", true, nil
}

var solHistoryThrottle sync.Mutex
var solHistoryLast time.Time

func (c *LedgerClient) throttleSolanaHistory() {
	ms := 350
	if v := strings.TrimSpace(os.Getenv("SOLANA_HISTORY_RATE_LIMIT_MS")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			ms = n
		}
	}
	solHistoryThrottle.Lock()
	defer solHistoryThrottle.Unlock()
	if wait := time.Duration(ms)*time.Millisecond - time.Since(solHistoryLast); wait > 0 {
		time.Sleep(wait)
	}
	solHistoryLast = time.Now()
}

func parseSolanaIndexBalanceRaw(body []byte) (uint64, error) {
	var out struct {
		BalanceRaw string `json:"balanceRaw"`
		Error      string `json:"error"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return 0, fmt.Errorf("solana history decode: %w", err)
	}
	if out.Error != "" {
		return 0, fmt.Errorf("solana history: %s", out.Error)
	}
	if out.BalanceRaw == "" {
		return 0, fmt.Errorf("solana history: empty balanceRaw")
	}
	v, err := strconv.ParseUint(out.BalanceRaw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("solana history balanceRaw: %w", err)
	}
	return v, nil
}
