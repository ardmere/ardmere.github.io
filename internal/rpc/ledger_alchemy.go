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
	"time"
)

// alchemyUTXOHosts maps ledger slug to Alchemy UTXO API origin (overridable in tests).
var alchemyUTXOHosts = map[string]string{
	"bitcoin":      "https://bitcoin-mainnet.g.alchemy.com",
	"bitcoin-cash": "https://bitcoincash-mainnet.g.alchemy.com",
}

// AlchemyBalanceAtHeight returns address balance in satoshis at block height via Alchemy UTXO API.
// chain is "bitcoin" or "bitcoin-cash". Requires ALCHEMY_KEY.
func (c *LedgerClient) AlchemyBalanceAtHeight(ctx context.Context, chain, addr string, height int64) (int64, string, error) {
	key := c.alchemyKey
	if key == "" {
		key = os.Getenv("ALCHEMY_KEY")
	}
	if key == "" {
		return 0, "", fmt.Errorf("ALCHEMY_KEY not set")
	}
	host, ok := alchemyUTXOHosts[chain]
	if !ok {
		return 0, "", fmt.Errorf("unsupported alchemy utxo chain %q", chain)
	}

	cacheChain := "alchemy:" + chain
	if raw, ok := c.cacheGet(cacheChain, "balance", addr, height); ok {
		v, _ := strconv.ParseInt(string(raw), 10, 64)
		return v, "cache", nil
	}

	u := fmt.Sprintf("%s/v2/%s/api/v2/address/%s?to=%d&details=basic",
		host, url.PathEscape(key), url.PathEscape(addr), height)
	bal, used, err := c.alchemyParseAddress(ctx, u)
	if err != nil {
		return 0, redactAlchemyKey(used, key), err
	}
	c.cachePut(cacheChain, "balance", addr, height, []byte(strconv.FormatInt(bal, 10)))
	return bal, redactAlchemyKey(used, key), nil
}

func (c *LedgerClient) alchemyParseAddress(ctx context.Context, u string) (int64, string, error) {
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * 2 * time.Second)
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
		if err != nil {
			return 0, u, err
		}
		req.Header.Set("User-Agent", "ardmere/0.1")
		resp, err := c.httpc.Do(req)
		if err != nil {
			return 0, u, err
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode == http.StatusTooManyRequests {
			lastErr = fmt.Errorf("alchemy HTTP 429")
			continue
		}
		if resp.StatusCode != http.StatusOK {
			return 0, u, fmt.Errorf("alchemy HTTP %d", resp.StatusCode)
		}
		var out struct {
			Balance string `json:"balance"`
			Error   string `json:"error"`
		}
		if err := json.Unmarshal(body, &out); err != nil {
			return 0, u, err
		}
		if out.Error != "" {
			return 0, u, fmt.Errorf("alchemy: %s", out.Error)
		}
		if out.Balance == "" {
			return 0, u, fmt.Errorf("alchemy: empty balance")
		}
		bal, err := strconv.ParseInt(out.Balance, 10, 64)
		if err != nil {
			return 0, u, fmt.Errorf("alchemy: parse balance %q: %w", out.Balance, err)
		}
		return bal, u, nil
	}
	if lastErr != nil {
		return 0, u, lastErr
	}
	return 0, u, fmt.Errorf("alchemy: request failed")
}

func redactAlchemyKey(u, key string) string {
	if key == "" {
		return u
	}
	return strings.Replace(u, key, "${ALCHEMY_KEY}", 1)
}
