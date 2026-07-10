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

const (
	solanaHistoryProviderHelius      = "helius"
	solanaHistoryProviderSolanaIndex = "solanaindex"
	heliusNativeMint                 = "So11111111111111111111111111111111111111111"
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

// SolanaNativeBalanceAtSlot returns lamports at owner address. When a historical
// provider is configured and slot > 0, it queries slot-indexed balance first;
// otherwise it falls back to live getBalance.
func (c *LedgerClient) SolanaNativeBalanceAtSlot(ctx context.Context, addr string, slot int64) (lamports uint64, used string, historical bool, err error) {
	if slot > 0 {
		if raw, u, ok, herr := c.solanaHistoryNativeAtSlot(ctx, addr, slot); ok {
			return raw, u, true, herr
		}
	}
	raw, u, err := c.SolanaNativeBalance(ctx, addr)
	return raw, u, false, err
}

func (c *LedgerClient) solanaHistoryNativeAtSlot(ctx context.Context, addr string, slot int64) (uint64, string, bool, error) {
	cfg, ok := loadSolanaHistoryConfig()
	if !ok {
		return 0, "", false, nil
	}
	cacheChain := "solana:history:" + cfg.Provider + ":native"
	if raw, ok := c.cacheGet(cacheChain, "balance", addr, slot); ok {
		v, _ := strconv.ParseUint(string(raw), 10, 64)
		return v, cfg.displayBase(), true, nil
	}
	switch cfg.Provider {
	case solanaHistoryProviderHelius:
		raw, used, err := c.heliusBalanceAtSlot(ctx, cfg, addr, heliusNativeMint, slot)
		if err != nil {
			return 0, used, true, err
		}
		c.cachePut(cacheChain, "balance", addr, slot, []byte(strconv.FormatUint(raw, 10)))
		return raw, used, true, nil
	default:
		paths := []string{
			fmt.Sprintf("%s/api/v1/solana/native-balance/%s/%d", cfg.BaseURL, url.PathEscape(addr), slot),
			fmt.Sprintf("%s/api/v1/solana/balance/%s/%d", cfg.BaseURL, url.PathEscape(addr), slot),
		}
		for _, path := range paths {
			raw, used, err := c.solanaIndexHistoryGET(ctx, cfg, path)
			if err != nil {
				if strings.Contains(err.Error(), "HTTP 404") {
					continue
				}
				return 0, used, true, err
			}
			c.cachePut(cacheChain, "balance", addr, slot, []byte(strconv.FormatUint(raw, 10)))
			return raw, used, true, nil
		}
		return 0, "", false, nil
	}
}

func (c *LedgerClient) solanaHistorySPLAtSlot(ctx context.Context, owner, mint string, slot int64) (uint64, string, bool, error) {
	cfg, ok := loadSolanaHistoryConfig()
	if !ok {
		return 0, "", false, nil
	}
	cacheChain := "solana:history:" + cfg.Provider + ":spl"
	cacheKey := owner + "|" + mint
	if raw, ok := c.cacheGet(cacheChain, "balance", cacheKey, slot); ok {
		v, _ := strconv.ParseUint(string(raw), 10, 64)
		return v, cfg.displayBase(), true, nil
	}
	var raw uint64
	var used string
	var err error
	switch cfg.Provider {
	case solanaHistoryProviderHelius:
		raw, used, err = c.heliusBalanceAtSlot(ctx, cfg, owner, mint, slot)
	default:
		path := fmt.Sprintf("%s/api/v1/solana/token-balance/%s/%s/%d",
			cfg.BaseURL, url.PathEscape(owner), url.PathEscape(mint), slot)
		raw, used, err = c.solanaIndexHistoryGET(ctx, cfg, path)
	}
	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			return 0, "", false, nil
		}
		return 0, used, true, err
	}
	c.cachePut(cacheChain, "balance", cacheKey, slot, []byte(strconv.FormatUint(raw, 10)))
	return raw, used, true, nil
}

func (c *LedgerClient) heliusBalanceAtSlot(ctx context.Context, cfg solanaHistoryConfig, wallet, mint string, slot int64) (uint64, string, error) {
	u := fmt.Sprintf("%s/v1/wallet/%s/balance-at?mint=%s&slot=%d",
		cfg.BaseURL, url.PathEscape(wallet), url.QueryEscape(mint), slot)
	return c.solanaHistoryGET(ctx, cfg, u, cfg.displayBase()+"/v1/wallet")
}

func (c *LedgerClient) solanaIndexHistoryGET(ctx context.Context, cfg solanaHistoryConfig, path string) (uint64, string, error) {
	raw, used, err := c.solanaHistoryGET(ctx, cfg, path, cfg.BaseURL+"/api/v1")
	if err != nil {
		return 0, used, err
	}
	return raw, used, nil
}

func (c *LedgerClient) solanaHistoryGET(ctx context.Context, cfg solanaHistoryConfig, reqURL, usedLabel string) (uint64, string, error) {
	var body []byte
	var status int
	var reqErr error
	for attempt := 0; attempt < 4; attempt++ {
		c.throttleSolanaHistory()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
		if err != nil {
			return 0, "", err
		}
		switch cfg.Provider {
		case solanaHistoryProviderHelius:
			req.Header.Set("X-Api-Key", cfg.APIKey)
		default:
			req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
		}
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
		return 0, "", reqErr
	}
	if status == http.StatusNotFound {
		return 0, usedLabel, fmt.Errorf("solana history HTTP 404")
	}
	if status == http.StatusUnauthorized || status == http.StatusForbidden {
		keyHint := "SOLANA_INDEX_API_KEY"
		if cfg.Provider == solanaHistoryProviderHelius {
			keyHint = "HELIUS_API_KEY"
		}
		return 0, usedLabel, fmt.Errorf("solana history HTTP %d: unauthorized (check %s)", status, keyHint)
	}
	if status != http.StatusOK {
		return 0, usedLabel, fmt.Errorf("solana history HTTP %d: %s", status, truncate(string(body), 120))
	}
	raw, err := parseSolanaHistoryBalanceRaw(body)
	if err != nil {
		return 0, usedLabel, err
	}
	return raw, usedLabel, nil
}

type solanaHistoryConfig struct {
	Provider string
	BaseURL  string
	APIKey   string
}

func (c solanaHistoryConfig) displayBase() string {
	return c.BaseURL
}

func loadSolanaHistoryConfig() (solanaHistoryConfig, bool) {
	if key := strings.TrimSpace(os.Getenv("HELIUS_API_KEY")); key != "" {
		base := strings.TrimRight(strings.TrimSpace(os.Getenv("HELIUS_API_URL")), "/")
		if base == "" {
			base = "https://api.helius.xyz"
		}
		return solanaHistoryConfig{
			Provider: solanaHistoryProviderHelius,
			BaseURL:  base,
			APIKey:   key,
		}, true
	}
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
	return solanaHistoryConfig{
		Provider: solanaHistoryProviderSolanaIndex,
		BaseURL:  base,
		APIKey:   key,
	}, true
}

// SolanaHistoryProviderFromUsed maps a history RPC label/URL to a provider name.
func SolanaHistoryProviderFromUsed(used string) string {
	u := strings.ToLower(used)
	switch {
	case strings.Contains(u, "helius"):
		return solanaHistoryProviderHelius
	case strings.Contains(u, "solanaindex"):
		return solanaHistoryProviderSolanaIndex
	default:
		return used
	}
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

func parseSolanaHistoryBalanceRaw(body []byte) (uint64, error) {
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
