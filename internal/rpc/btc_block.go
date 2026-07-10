package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// EsploraBases are public Esplora-compatible APIs for BTC block metadata.
var EsploraBases = []string{
	"https://mempool.space/api",
	"https://blockstream.info/api",
}

type esploraBlock struct {
	Timestamp int64  `json:"timestamp"`
	Height    int64  `json:"height"`
	Hash      string `json:"id"`
}

// BTCBlockTimeAtHeight returns the UTC block time for a BTC mainnet height
// via Esplora-compatible APIs (mempool.space, blockstream.info).
func (c *LedgerClient) BTCBlockTimeAtHeight(ctx context.Context, height int64) (time.Time, string, error) {
	if height <= 0 {
		return time.Time{}, "", fmt.Errorf("invalid btc block height %d", height)
	}
	chain := "btc:mainnet"
	if raw, ok := c.cacheGet(chain, "blocktime", "", height); ok {
		sec, err := strconv.ParseInt(string(raw), 10, 64)
		if err == nil {
			return time.Unix(sec, 0).UTC(), "cache", nil
		}
	}

	var lastErr error
	for _, base := range EsploraBases {
		t, used, err := c.btcBlockTimeEsplora(ctx, base, height)
		if err != nil {
			lastErr = err
			continue
		}
		c.cachePut(chain, "blocktime", "", height, []byte(strconv.FormatInt(t.Unix(), 10)))
		return t, used, nil
	}
	if lastErr != nil {
		return time.Time{}, "", lastErr
	}
	return time.Time{}, "", fmt.Errorf("no esplora provider configured")
}

func (c *LedgerClient) btcBlockTimeEsplora(ctx context.Context, apiBase string, height int64) (time.Time, string, error) {
	base := strings.TrimRight(apiBase, "/")
	used := base

	hashURL := fmt.Sprintf("%s/block-height/%d", base, height)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, hashURL, nil)
	if err != nil {
		return time.Time{}, used, err
	}
	req.Header.Set("User-Agent", "ardmere/0.1")
	resp, err := c.httpc.Do(req)
	if err != nil {
		return time.Time{}, used, err
	}
	hashBody, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return time.Time{}, used, fmt.Errorf("%s -> HTTP %d", hashURL, resp.StatusCode)
	}
	hash := strings.TrimSpace(string(hashBody))
	if hash == "" {
		return time.Time{}, used, fmt.Errorf("empty block hash for height %d", height)
	}

	blockURL := fmt.Sprintf("%s/block/%s", base, hash)
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, blockURL, nil)
	if err != nil {
		return time.Time{}, used, err
	}
	req.Header.Set("User-Agent", "ardmere/0.1")
	resp, err = c.httpc.Do(req)
	if err != nil {
		return time.Time{}, used, err
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return time.Time{}, used, fmt.Errorf("%s -> HTTP %d", blockURL, resp.StatusCode)
	}
	var blk esploraBlock
	if err := json.Unmarshal(body, &blk); err != nil {
		return time.Time{}, used, err
	}
	if blk.Timestamp <= 0 {
		return time.Time{}, used, fmt.Errorf("block %d missing timestamp", height)
	}
	return time.Unix(blk.Timestamp, 0).UTC(), used, nil
}
