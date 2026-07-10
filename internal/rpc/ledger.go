package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// LedgerClient queries non-EVM chain balances (UTXO, Solana, XRPL).
type LedgerClient struct {
	httpc           *http.Client
	cache           *ResultCache
	blockchairKey   string
	alchemyKey      string
	solanaProviders []Provider
	xrplRPC         string

	mu          sync.Mutex
	solDisabled map[string]time.Time
	solCooldown time.Duration
}

// NewLedger returns a LedgerClient with public defaults.
func NewLedger() *LedgerClient {
	return &LedgerClient{
		httpc:           &http.Client{Timeout: 90 * time.Second},
		cache:           NewResultCache(""),
		blockchairKey:   os.Getenv("BLOCKCHAIR_API_KEY"),
		alchemyKey:      os.Getenv("ALCHEMY_KEY"),
		solanaProviders: loadSolanaProviders(),
		xrplRPC:         loadXRPLRPC(),
		solDisabled:     map[string]time.Time{},
		solCooldown:     5 * time.Minute,
	}
}

func loadSolanaProviders() []Provider {
	if u := os.Getenv("SOLANA_RPC"); u != "" {
		ms := 100
		if v := os.Getenv("SOLANA_RATE_LIMIT_MS"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n >= 0 {
				ms = n
			}
		}
		return []Provider{{URL: u, Weight: 100, RateLimitMs: ms}}
	}
	cfg, err := LoadProviderConfig()
	if err != nil {
		cfg = DefaultProviderConfig()
	}
	if p := ProvidersFor(cfg, NetSolana, false); len(p) > 0 {
		return p
	}
	return DefaultProviderConfig()[NetSolana]
}

func loadXRPLRPC() string {
	xrpl := "https://xrplcluster.com/"
	if u := os.Getenv("XRPL_RPC"); u != "" {
		xrpl = u
	}
	return xrpl
}

func (c *LedgerClient) isSolDisabled(u string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	until, ok := c.solDisabled[u]
	return ok && time.Now().Before(until)
}

func (c *LedgerClient) disableSol(u string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.solDisabled[u] = time.Now().Add(c.solCooldown)
}

func (c *LedgerClient) cacheGet(chain, method, target string, height int64) ([]byte, bool) {
	return c.cache.Get(Network(chain), method, target, "", height)
}

func (c *LedgerClient) cachePut(chain, method, target string, height int64, result []byte) {
	c.cache.Put(Network(chain), method, target, "", height, result)
}

// EsploraBalanceAtHeight sums confirmed UTXOs for addr up to block height via Esplora tx pagination.
func (c *LedgerClient) EsploraBalanceAtHeight(ctx context.Context, apiBase, addr string, height int64) (int64, string, error) {
	var lastErr error
	for _, base := range esploraBasesToTry(apiBase) {
		bal, used, err := c.esploraBalanceAtHeightOnce(ctx, base, addr, height)
		if err == nil {
			return bal, used, nil
		}
		lastErr = err
		if !isEsploraRetryable(err) {
			return 0, used, err
		}
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("esplora: no usable base for %s", apiBase)
	}
	return 0, "", lastErr
}

func esploraBasesToTry(primary string) []string {
	out := []string{primary}
	if strings.Contains(primary, "blockstream") || strings.Contains(primary, "mempool.space") {
		for _, b := range EsploraBases {
			if b != primary {
				out = append(out, b)
			}
		}
	}
	return out
}

func isEsploraRetryable(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "HTTP 429") ||
		strings.Contains(msg, "HTTP 502") ||
		strings.Contains(msg, "HTTP 503") ||
		strings.Contains(msg, "HTTP 504")
}

func (c *LedgerClient) esploraBalanceAtHeightOnce(ctx context.Context, apiBase, addr string, height int64) (int64, string, error) {
	chain := "esplora:" + apiBase
	if raw, ok := c.cacheGet(chain, "balance", addr, height); ok {
		v, _ := strconv.ParseInt(string(raw), 10, 64)
		return v, "cache", nil
	}

	base := strings.TrimRight(apiBase, "/")
	var balance int64
	lastTxid := ""
	pages := 0
	used := base
	for pages < 400 {
		path := fmt.Sprintf("%s/address/%s/txs/chain", base, url.PathEscape(addr))
		if lastTxid != "" {
			path += "/" + lastTxid
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		if err != nil {
			return 0, used, err
		}
		req.Header.Set("User-Agent", "ardmere/0.1")
		resp, err := c.httpc.Do(req)
		if err != nil {
			return 0, used, err
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return 0, used, fmt.Errorf("%s -> HTTP %d", path, resp.StatusCode)
		}
		var txs []esploraTx
		if err := json.Unmarshal(body, &txs); err != nil {
			return 0, used, err
		}
		if len(txs) == 0 {
			break
		}
		pages++
		for _, tx := range txs {
			if tx.Status.BlockHeight == nil || *tx.Status.BlockHeight > height {
				continue
			}
			for _, vin := range tx.Vin {
				if vin.Prevout != nil && vin.Prevout.ScriptPubkeyAddress == addr {
					balance -= vin.Prevout.Value
				}
			}
			for _, vout := range tx.Vout {
				if vout.ScriptPubkeyAddress == addr {
					balance += vout.Value
				}
			}
		}
		lastTxid = txs[len(txs)-1].Txid
		if len(txs) < 25 {
			break
		}
	}
	if pages >= 400 {
		return 0, used, fmt.Errorf("esplora tx pagination exceeded 400 pages for %s", addr)
	}
	c.cachePut(chain, "balance", addr, height, []byte(strconv.FormatInt(balance, 10)))
	return balance, used, nil
}

type esploraTx struct {
	Txid   string `json:"txid"`
	Status struct {
		BlockHeight *int64 `json:"block_height"`
	} `json:"status"`
	Vin []struct {
		Prevout *struct {
			ScriptPubkeyAddress string `json:"scriptpubkey_address"`
			Value               int64  `json:"value"`
		} `json:"prevout"`
	} `json:"vin"`
	Vout []struct {
		ScriptPubkeyAddress string `json:"scriptpubkey_address"`
		Value               int64  `json:"value"`
	} `json:"vout"`
}

// BlockchairBalanceAtHeight returns address balance at block state when API key is configured.
// Without BLOCKCHAIR_API_KEY it falls back to the latest Blockchair dashboard (live snapshot).
func (c *LedgerClient) BlockchairBalanceAtHeight(ctx context.Context, chain, addr string, height int64) (int64, string, bool, error) {
	if c.blockchairKey == "" {
		return c.blockchairBalanceLive(ctx, chain, addr)
	}
	cacheChain := "blockchair:" + chain
	if raw, ok := c.cacheGet(cacheChain, "balance", addr, height); ok {
		v, _ := strconv.ParseInt(string(raw), 10, 64)
		return v, "cache", false, nil
	}
	u := fmt.Sprintf("https://api.blockchair.com/%s/dashboards/address/%s?state=%d&key=%s",
		chain, url.PathEscape(addr), height, url.QueryEscape(c.blockchairKey))
	bal, used, err := c.blockchairParseDashboard(ctx, u)
	if err != nil {
		return 0, used, false, err
	}
	c.cachePut(cacheChain, "balance", addr, height, []byte(strconv.FormatInt(bal, 10)))
	return bal, used, false, nil
}

func (c *LedgerClient) blockchairBalanceLive(ctx context.Context, chain, addr string) (int64, string, bool, error) {
	if c.blockchairKey == "" {
		time.Sleep(750 * time.Millisecond)
	}
	u := fmt.Sprintf("https://api.blockchair.com/%s/dashboards/address/%s", chain, url.PathEscape(addr))
	bal, used, err := c.blockchairParseDashboard(ctx, u)
	if err != nil {
		return 0, used, true, err
	}
	return bal, used, true, nil
}

func (c *LedgerClient) blockchairParseDashboard(ctx context.Context, u string) (int64, string, error) {
	var lastErr error
	maxAttempts := 3
	if c.blockchairKey == "" {
		maxAttempts = 5
	}
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * 2 * time.Second)
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
		if err != nil {
			return 0, "", err
		}
		resp, err := c.httpc.Do(req)
		if err != nil {
			return 0, "", err
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode == 430 {
			lastErr = fmt.Errorf("blockchair HTTP 430")
			continue
		}
		if resp.StatusCode != http.StatusOK {
			return 0, u, fmt.Errorf("blockchair HTTP %d", resp.StatusCode)
		}
		var out struct {
			Data map[string]struct {
				Address struct {
					Balance int64 `json:"balance"`
				} `json:"address"`
			} `json:"data"`
			Context struct {
				Error string `json:"error"`
			} `json:"context"`
		}
		if err := json.Unmarshal(body, &out); err != nil {
			return 0, u, err
		}
		if out.Context.Error != "" {
			return 0, u, fmt.Errorf("blockchair: %s", out.Context.Error)
		}
		for _, v := range out.Data {
			return v.Address.Balance, u, nil
		}
		return 0, u, fmt.Errorf("blockchair: empty data")
	}
	if lastErr != nil {
		return 0, u, lastErr
	}
	return 0, u, fmt.Errorf("blockchair: request failed")
}

// BlockcypherBalanceBefore returns balance before block height (approximate for high-volume addresses).
func (c *LedgerClient) BlockcypherBalanceBefore(ctx context.Context, chain, addr string, height int64) (int64, int, string, error) {
	cacheChain := "blockcypher:" + chain
	if raw, ok := c.cacheGet(cacheChain, "balance", addr, height); ok {
		parts := strings.Split(string(raw), "|")
		if len(parts) == 2 {
			v, _ := strconv.ParseInt(parts[0], 10, 64)
			n, _ := strconv.Atoi(parts[1])
			return v, n, "cache", nil
		}
	}
	u := fmt.Sprintf("https://api.blockcypher.com/v1/%s/main/addrs/%s?before=%d", chain, url.PathEscape(addr), height+1)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return 0, 0, "", err
	}
	resp, err := c.httpc.Do(req)
	if err != nil {
		return 0, 0, "", err
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	status := resp.StatusCode
	if status == 429 || strings.Contains(string(body), "Limits reached") {
		time.Sleep(2 * time.Second)
		req2, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
		if resp2, err2 := c.httpc.Do(req2); err2 == nil {
			body, _ = io.ReadAll(resp2.Body)
			resp2.Body.Close()
			status = resp2.StatusCode
		}
	}
	if status != http.StatusOK && !strings.Contains(string(body), `"balance"`) {
		return 0, 0, u, fmt.Errorf("blockcypher HTTP %d: %s", status, truncate(string(body), 80))
	}
	var out struct {
		Balance  int64 `json:"balance"`
		FinalNTx int   `json:"final_n_tx"`
		NTx      int   `json:"n_tx"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return 0, 0, u, err
	}
	n := out.FinalNTx
	if n == 0 {
		n = out.NTx
	}
	c.cachePut(cacheChain, "balance", addr, height, []byte(fmt.Sprintf("%d|%d", out.Balance, n)))
	return out.Balance, n, u, nil
}

// SolanaNativeBalance returns lamports for a base58 address (live node; not historical slot).
func (c *LedgerClient) SolanaNativeBalance(ctx context.Context, addr string) (uint64, string, error) {
	var result struct {
		Value uint64 `json:"value"`
	}
	used, err := c.solanaCall(ctx, "getBalance", []any{addr, map[string]string{"commitment": "confirmed"}}, &result)
	return result.Value, used, err
}

// SolanaSPLBalance sums all token accounts for owner+mint (live node).
func (c *LedgerClient) SolanaSPLBalance(ctx context.Context, owner, mint string) (uint64, string, error) {
	var result struct {
		Value []struct {
			Account struct {
				Data struct {
					Parsed struct {
						Info struct {
							TokenAmount struct {
								Amount string `json:"amount"`
							} `json:"tokenAmount"`
						} `json:"info"`
					} `json:"parsed"`
				} `json:"data"`
			} `json:"account"`
		} `json:"value"`
	}
	filter := map[string]string{"mint": mint}
	used, err := c.solanaCall(ctx, "getTokenAccountsByOwner", []any{owner, filter, map[string]string{"encoding": "jsonParsed"}}, &result)
	if err != nil {
		return 0, used, err
	}
	sum := new(big.Int)
	for _, acc := range result.Value {
		amt := acc.Account.Data.Parsed.Info.TokenAmount.Amount
		if amt == "" {
			continue
		}
		v, ok := new(big.Int).SetString(amt, 10)
		if ok {
			sum.Add(sum, v)
		}
	}
	if !sum.IsUint64() {
		return 0, used, fmt.Errorf("spl balance overflow")
	}
	return sum.Uint64(), used, nil
}

func (c *LedgerClient) solanaCall(ctx context.Context, method string, params []any, result any) (string, error) {
	target := solanaCacheTarget(method, params)
	if raw, ok := c.cacheGet("solana", method, target, 0); ok {
		if err := json.Unmarshal(raw, result); err == nil {
			return "cache", nil
		}
	}

	body, _ := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": 1, "method": method, "params": params})
	var lastErr error
	for _, p := range c.solanaProviders {
		u := p.URL
		if c.isSolDisabled(u) {
			continue
		}
		if p.RateLimitMs > 0 {
			time.Sleep(time.Duration(p.RateLimitMs) * time.Millisecond)
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
		if err != nil {
			lastErr = err
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "ardmere/0.1")

		resp, err := c.httpc.Do(req)
		if err != nil {
			c.disableSol(u)
			lastErr = err
			continue
		}
		raw, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode == http.StatusTooManyRequests {
			lastErr = fmt.Errorf("solana HTTP 429")
			time.Sleep(2 * time.Second)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			c.disableSol(u)
			lastErr = fmt.Errorf("solana HTTP %d", resp.StatusCode)
			continue
		}
		var envelope struct {
			Result json.RawMessage `json:"result"`
			Error  *struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal(raw, &envelope); err != nil {
			c.disableSol(u)
			lastErr = err
			continue
		}
		if envelope.Error != nil {
			lastErr = fmt.Errorf("solana %s: %s", method, envelope.Error.Message)
			if envelope.Error.Code == -32005 || strings.Contains(strings.ToLower(envelope.Error.Message), "rate") {
				time.Sleep(2 * time.Second)
			}
			continue
		}
		if err := json.Unmarshal(envelope.Result, result); err != nil {
			lastErr = err
			continue
		}
		c.cachePut("solana", method, target, 0, envelope.Result)
		return u, nil
	}
	if lastErr == nil {
		lastErr = errors.New("no usable solana provider")
	}
	return "", lastErr
}

func solanaCacheTarget(method string, params []any) string {
	switch method {
	case "getBalance":
		if len(params) > 0 {
			return fmt.Sprint(params[0])
		}
	case "getTokenAccountsByOwner":
		if len(params) >= 2 {
			return fmt.Sprint(params[0]) + "|" + fmt.Sprint(params[1])
		}
	}
	b, _ := json.Marshal(params)
	return method + ":" + string(b)
}

// XRPLBalanceAtLedger returns XRP balance in drops at ledger_index.
func (c *LedgerClient) XRPLBalanceAtLedger(ctx context.Context, addr string, ledgerIndex int64) (int64, string, error) {
	cacheChain := "xrpl"
	if raw, ok := c.cacheGet(cacheChain, "balance", addr, ledgerIndex); ok {
		v, _ := strconv.ParseInt(string(raw), 10, 64)
		return v, "cache", nil
	}
	body, _ := json.Marshal(map[string]any{
		"method": "account_info",
		"params": []any{map[string]any{
			"account":       addr,
			"ledger_index":  ledgerIndex,
			"strict":        true,
		}},
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.xrplRPC, bytes.NewReader(body))
	if err != nil {
		return 0, "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpc.Do(req)
	if err != nil {
		return 0, "", err
	}
	raw, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, c.xrplRPC, fmt.Errorf("xrpl HTTP %d", resp.StatusCode)
	}
	var out struct {
		Result struct {
			AccountData struct {
				Balance string `json:"Balance"`
			} `json:"account_data"`
			Status string `json:"status"`
		} `json:"result"`
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return 0, c.xrplRPC, err
	}
	if out.Result.Status != "success" {
		return 0, c.xrplRPC, fmt.Errorf("xrpl account_info failed for %s at %d", addr, ledgerIndex)
	}
	bal, err := strconv.ParseInt(out.Result.AccountData.Balance, 10, 64)
	if err != nil {
		return 0, c.xrplRPC, err
	}
	c.cachePut(cacheChain, "balance", addr, ledgerIndex, []byte(strconv.FormatInt(bal, 10)))
	return bal, c.xrplRPC, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

// CipherScanZECBalance returns transparent ZEC balance in zatoshis via cipherscan.app (live).
func (c *LedgerClient) CipherScanZECBalance(ctx context.Context, addr string) (int64, string, error) {
	cacheChain := "cipherscan:zec"
	if raw, ok := c.cacheGet(cacheChain, "balance", addr, 0); ok {
		v, _ := strconv.ParseInt(string(raw), 10, 64)
		return v, "cache", nil
	}
	u := fmt.Sprintf("https://cipherscan.app/api/address/%s", url.PathEscape(addr))
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
	if resp.StatusCode != http.StatusOK {
		return 0, u, fmt.Errorf("cipherscan HTTP %d", resp.StatusCode)
	}
	var out struct {
		Balance float64 `json:"balance"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return 0, u, err
	}
	zat := int64(out.Balance * 1e8)
	c.cachePut(cacheChain, "balance", addr, 0, []byte(strconv.FormatInt(zat, 10)))
	return zat, u, nil
}
