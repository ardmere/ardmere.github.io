// Package rpc provides minimal JSON-RPC clients for public EVM nodes.
// All endpoints listed here are *public, free* providers — no API key,
// per ADR-005 (no self-hosted archive infrastructure).
package rpc

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Network identifies a chain we know how to query.
type Network string

const (
	NetEthereum  Network = "ETH"
	NetBSC       Network = "BSC"
	NetArbitrum  Network = "ARBITRUM"
	NetOptimism  Network = "OPTIMISM"
	NetBase      Network = "BASE"
	NetPolygon   Network = "MATIC"
	NetAvalanche Network = "AVAXC"
	NetOpBNB     Network = "OPBNB"
	NetSonic     Network = "SONIC"
	NetWorld     Network = "WLD"
	NetCelo      Network = "CELO"
	NetPlasma    Network = "PLASMA"
	NetKaia      Network = "KAIA"
	NetKavaEVM   Network = "KAVAEVM"
	NetZkSync    Network = "ZKSYNCERA"
	NetRon       Network = "RON"
	NetSeiEVM    Network = "SEIEVM"
	NetManta     Network = "MANTA"
	NetXLayer    Network = "XLAYER"
	NetFEVM      Network = "FEVM"
	NetSolana    Network = "SOL"
	NetAB        Network = "AB"
)

// PublicProviders lists fallback RPC endpoints for each supported network.
// Populated from config/rpc-providers.json when available.
var PublicProviders = map[Network][]string{}

func init() {
	cfg, err := LoadProviderConfig()
	if err != nil {
		cfg = DefaultProviderConfig()
	}
	for net, providers := range cfg {
		urls := make([]string, 0, len(providers))
		for _, p := range providers {
			urls = append(urls, p.URL)
		}
		PublicProviders[net] = urls
	}
}

// Client is a tiny, safe-for-concurrent-use JSON-RPC pool with provider failover.
type Client struct {
	httpc     *http.Client
	providers map[Network][]Provider
	cache     *ResultCache

	mu       sync.Mutex
	disabled map[string]time.Time // url -> until
	cooldown time.Duration
}

// New returns a Client with sane defaults.
func New() *Client {
	providers, err := LoadProviderConfig()
	if err != nil {
		providers = DefaultProviderConfig()
	}
	return &Client{
		httpc:     &http.Client{Timeout: 20 * time.Second},
		providers: providers,
		cache:     NewResultCache(""),
		disabled:  map[string]time.Time{},
		cooldown:  5 * time.Minute,
	}
}

type rpcReq struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      int    `json:"id"`
}

type rpcResp struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// ContractCall is one eth_call in a JSON-RPC batch.
type ContractCall struct {
	ID   int
	To   string
	Data []byte
}

const defaultBatchSize = 40

// call dispatches one JSON-RPC request, trying providers in order until one
// returns a non-error response.
func (c *Client) call(ctx context.Context, net Network, method string, params []any, historical bool) (json.RawMessage, string, error) {
	providers := ProvidersFor(c.providers, net, historical)
	if len(providers) == 0 {
		return nil, "", fmt.Errorf("network %s not configured", net)
	}
	body, _ := json.Marshal(rpcReq{Jsonrpc: "2.0", Method: method, Params: params, ID: 1})

	var lastErr error
	for _, p := range providers {
		u := p.URL
		if c.isDisabled(u) {
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
			c.disable(u)
			lastErr = err
			continue
		}
		raw, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode == http.StatusTooManyRequests {
			lastErr = fmt.Errorf("%s -> HTTP 429", u)
			time.Sleep(2 * time.Second)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			c.disable(u)
			lastErr = fmt.Errorf("%s -> HTTP %d", u, resp.StatusCode)
			continue
		}
		var rr rpcResp
		if err := json.Unmarshal(raw, &rr); err != nil {
			c.disable(u)
			lastErr = err
			continue
		}
		if rr.Error != nil {
			lastErr = fmt.Errorf("%s -> rpc error %d: %s", u, rr.Error.Code, rr.Error.Message)
			continue
		}
		return rr.Result, u, nil
	}
	if lastErr == nil {
		lastErr = errors.New("no usable provider")
	}
	return nil, "", lastErr
}

func (c *Client) callBatch(ctx context.Context, net Network, calls []rpcReq, historical bool) ([]rpcResp, string, error) {
	if len(calls) == 0 {
		return nil, "", nil
	}
	providers := ProvidersFor(c.providers, net, historical)
	if len(providers) == 0 {
		return nil, "", fmt.Errorf("network %s not configured", net)
	}
	body, _ := json.Marshal(calls)

	var lastErr error
	for _, p := range providers {
		u := p.URL
		if c.isDisabled(u) {
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
			c.disable(u)
			lastErr = err
			continue
		}
		raw, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode == http.StatusTooManyRequests {
			lastErr = fmt.Errorf("%s -> HTTP 429", u)
			time.Sleep(2 * time.Second)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			c.disable(u)
			lastErr = fmt.Errorf("%s -> HTTP %d", u, resp.StatusCode)
			continue
		}
		var batch []rpcResp
		if err := json.Unmarshal(raw, &batch); err != nil {
			var single rpcResp
			if err2 := json.Unmarshal(raw, &single); err2 != nil {
				c.disable(u)
				lastErr = err
				continue
			}
			batch = []rpcResp{single}
		}
		if len(batch) != len(calls) {
			lastErr = fmt.Errorf("%s -> batch size mismatch: sent %d got %d", u, len(calls), len(batch))
			continue
		}
		for _, rr := range batch {
			if rr.Error != nil {
				lastErr = fmt.Errorf("%s -> rpc error %d: %s", u, rr.Error.Code, rr.Error.Message)
				goto nextProvider
			}
		}
		if historical {
			for _, rr := range batch {
				out, err := decodeHexBytes(rr.Result)
				out = normalizeEthCallResult(out)
				if err != nil || len(out) < 32 {
					n := 0
					if out != nil {
						n = len(out)
					}
					lastErr = fmt.Errorf("%s -> short eth_call in batch: %d bytes", u, n)
					goto nextProvider
				}
			}
		}
		return batch, u, nil
	nextProvider:
	}
	if lastErr == nil {
		lastErr = errors.New("no usable provider")
	}
	return nil, "", lastErr
}

func (c *Client) isDisabled(u string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	until, ok := c.disabled[u]
	if !ok {
		return false
	}
	if time.Now().After(until) {
		delete(c.disabled, u)
		return false
	}
	return true
}

func (c *Client) disable(u string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.disabled[u] = time.Now().Add(c.cooldown)
}

// GetBalance returns the native (ETH/BNB) balance of address at the given block height.
func (c *Client) GetBalance(ctx context.Context, net Network, address string, height int64) (*big.Int, string, error) {
	blk := blockTag(height)
	if c.cache != nil {
		if raw, ok := c.cache.Get(net, "eth_getBalance", address, blk, height); ok {
			bi, err := decodeHexBigInt(raw)
			if err == nil {
				return bi, "cache", nil
			}
		}
	}
	historical := height > 0
	res, used, err := c.call(ctx, net, "eth_getBalance", []any{address, blk}, historical)
	if err != nil && historical && isEVMHistoricalUnavailable(err) {
		res, used, err = c.call(ctx, net, "eth_getBalance", []any{address, "latest"}, false)
	}
	if err != nil {
		return nil, used, err
	}
	if c.cache != nil {
		c.cache.Put(net, "eth_getBalance", address, blk, height, res)
	}
	bi, err := decodeHexBigInt(res)
	if err != nil {
		return nil, used, err
	}
	return bi, used, nil
}

// CallContract executes eth_call against a contract at the given historical block.
func (c *Client) CallContract(ctx context.Context, net Network, to string, data []byte, height int64) ([]byte, string, error) {
	dataHex := "0x" + hex.EncodeToString(data)
	if c.cache != nil {
		if raw, ok := c.cache.Get(net, "eth_call", to, dataHex, height); ok {
			out, err := decodeHexBytes(raw)
			if err == nil && len(out) >= 32 {
				return out, "cache", nil
			}
		}
	}
	callObj := map[string]string{
		"to":   to,
		"data": dataHex,
	}
	blk := blockTag(height)
	historical := height > 0
	providers := ProvidersFor(c.providers, net, historical)
	if len(providers) == 0 {
		return nil, "", fmt.Errorf("network %s not configured", net)
	}
	body, _ := json.Marshal(rpcReq{Jsonrpc: "2.0", Method: "eth_call", Params: []any{callObj, blk}, ID: 1})

	var lastErr error
	for _, p := range providers {
		u := p.URL
		if c.isDisabled(u) {
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
			c.disable(u)
			lastErr = err
			continue
		}
		raw, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode == http.StatusTooManyRequests {
			lastErr = fmt.Errorf("%s -> HTTP 429", u)
			time.Sleep(2 * time.Second)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			c.disable(u)
			lastErr = fmt.Errorf("%s -> HTTP %d", u, resp.StatusCode)
			continue
		}
		var rr rpcResp
		if err := json.Unmarshal(raw, &rr); err != nil {
			c.disable(u)
			lastErr = err
			continue
		}
		if rr.Error != nil {
			lastErr = fmt.Errorf("%s -> rpc error %d: %s", u, rr.Error.Code, rr.Error.Message)
			continue
		}
		out, err := decodeHexBytes(rr.Result)
		if err != nil {
			lastErr = err
			continue
		}
		out = normalizeEthCallResult(out)
		if historical && len(out) < 32 {
			lastErr = fmt.Errorf("%s -> short eth_call: %d bytes", u, len(out))
			continue
		}
		if c.cache != nil && len(out) >= 32 {
			c.cache.Put(net, "eth_call", to, dataHex, height, ethCallResultJSON(out))
		}
		return out, u, nil
	}
	if lastErr == nil {
		lastErr = errors.New("no usable provider")
	}
	if historical && isEVMHistoricalUnavailable(lastErr) {
		body, _ = json.Marshal(rpcReq{Jsonrpc: "2.0", Method: "eth_call", Params: []any{callObj, "latest"}, ID: 1})
		liveProviders := ProvidersFor(c.providers, net, false)
		for _, p := range liveProviders {
			u := p.URL
			if c.isDisabled(u) {
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
				lastErr = err
				continue
			}
			raw, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				lastErr = fmt.Errorf("%s -> HTTP %d", u, resp.StatusCode)
				continue
			}
			var rr rpcResp
			if err := json.Unmarshal(raw, &rr); err != nil {
				lastErr = err
				continue
			}
			if rr.Error != nil {
				lastErr = fmt.Errorf("%s -> rpc error %d: %s", u, rr.Error.Code, rr.Error.Message)
				continue
			}
			out, err := decodeHexBytes(rr.Result)
			if err != nil {
				lastErr = err
				continue
			}
			out = normalizeEthCallResult(out)
			if len(out) < 32 {
				lastErr = fmt.Errorf("%s -> short eth_call: %d bytes", u, len(out))
				continue
			}
			return out, u, nil
		}
	}
	return nil, "", lastErr
}

// CallContractBatch executes multiple eth_call requests in JSON-RPC batches.
func (c *Client) CallContractBatch(ctx context.Context, net Network, calls []ContractCall, height int64) (map[int][]byte, string, error) {
	out := make(map[int][]byte, len(calls))
	if len(calls) == 0 {
		return out, "", nil
	}
	blk := blockTag(height)
	used := ""
	for start := 0; start < len(calls); start += defaultBatchSize {
		end := start + defaultBatchSize
		if end > len(calls) {
			end = len(calls)
		}
		chunk := calls[start:end]
		var pending []ContractCall
		var rpcCalls []rpcReq
		for _, call := range chunk {
			dataHex := "0x" + hex.EncodeToString(call.Data)
			if c.cache != nil {
				if raw, ok := c.cache.Get(net, "eth_call", call.To, dataHex, height); ok {
					decoded, err := decodeHexBytes(raw)
					if err == nil && len(decoded) >= 32 {
						out[call.ID] = decoded
						continue
					}
				}
			}
			pending = append(pending, call)
			callObj := map[string]string{
				"to":   call.To,
				"data": dataHex,
			}
			rpcCalls = append(rpcCalls, rpcReq{
				Jsonrpc: "2.0",
				Method:  "eth_call",
				Params:  []any{callObj, blk},
				ID:      call.ID,
			})
		}
		if len(rpcCalls) == 0 {
			continue
		}
		responses, batchUsed, err := c.callBatch(ctx, net, rpcCalls, height > 0)
		if err != nil {
			return nil, batchUsed, err
		}
		if used == "" {
			used = batchUsed
		}
		respByID := map[int]rpcResp{}
		for _, rr := range responses {
			respByID[rr.ID] = rr
		}
		for _, call := range pending {
			rr, ok := respByID[call.ID]
			if !ok {
				return nil, used, fmt.Errorf("missing batch response for id %d", call.ID)
			}
			decoded, err := decodeHexBytes(rr.Result)
			if err != nil {
				return nil, used, fmt.Errorf("decode eth_call %d: %w", call.ID, err)
			}
			decoded = normalizeEthCallResult(decoded)
			out[call.ID] = decoded
			if c.cache != nil && len(decoded) >= 32 {
				dataHex := "0x" + hex.EncodeToString(call.Data)
				c.cache.Put(net, "eth_call", call.To, dataHex, height, ethCallResultJSON(decoded))
			}
		}
		if height > 0 {
			for _, call := range pending {
				if len(out[call.ID]) >= 32 {
					continue
				}
				retried, retryUsed, err := c.CallContract(ctx, net, call.To, call.Data, height)
				if err != nil {
					continue
				}
				out[call.ID] = retried
				if used == "" {
					used = retryUsed
				}
			}
		}
	}
	return out, used, nil
}

func decodeHexBigInt(raw json.RawMessage) (*big.Int, error) {
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, err
	}
	bi := new(big.Int)
	if _, ok := bi.SetString(strings.TrimPrefix(s, "0x"), 16); !ok {
		return nil, fmt.Errorf("bad balance hex: %s", s)
	}
	return bi, nil
}

func decodeHexBytes(raw json.RawMessage) ([]byte, error) {
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, err
	}
	out, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	if err != nil {
		return nil, fmt.Errorf("bad eth_call hex: %s", s)
	}
	return out, nil
}

// normalizeEthCallResult treats empty eth_call output as a 32-byte zero word (common for zero ERC20 balances).
func normalizeEthCallResult(out []byte) []byte {
	if len(out) == 0 {
		return make([]byte, 32)
	}
	return out
}

func ethCallResultJSON(out []byte) json.RawMessage {
	raw, _ := json.Marshal("0x" + hex.EncodeToString(out))
	return raw
}

func blockTag(height int64) string {
	blk := "0x" + strings.TrimLeft(fmt.Sprintf("%x", height), "0")
	if blk == "0x" {
		return "0x0"
	}
	return blk
}

func isEVMHistoricalUnavailable(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "missing trie") ||
		strings.Contains(msg, "pruned") ||
		strings.Contains(msg, "header not found") ||
		strings.Contains(msg, "historical state") ||
		strings.Contains(msg, "state is not available")
}

// GetBlockTime returns the unix timestamp of the given block on the given network.
func (c *Client) GetBlockTime(ctx context.Context, net Network, height int64) (int64, string, error) {
	blk := fmt.Sprintf("0x%x", height)
	res, used, err := c.call(ctx, net, "eth_getBlockByNumber", []any{blk, false}, height > 0)
	if err != nil {
		return 0, used, err
	}
	var blkJSON struct {
		Timestamp string `json:"timestamp"`
	}
	if err := json.Unmarshal(res, &blkJSON); err != nil {
		return 0, used, err
	}
	if blkJSON.Timestamp == "" {
		return 0, used, fmt.Errorf("block %d not found", height)
	}
	bi := new(big.Int)
	bi.SetString(strings.TrimPrefix(blkJSON.Timestamp, "0x"), 16)
	return bi.Int64(), used, nil
}

// To32Bytes pads bytes to 32-byte big-endian; useful when constructing eth_call data.
func To32Bytes(b []byte) []byte {
	if len(b) > 32 {
		panic("input > 32 bytes")
	}
	out := make([]byte, 32)
	copy(out[32-len(b):], b)
	return out
}

// EncodeHexAddress strips "0x" and lowercases, returning the 20-byte payload.
func EncodeHexAddress(addr string) ([]byte, error) {
	addr = strings.TrimPrefix(strings.ToLower(addr), "0x")
	if len(addr) != 40 {
		return nil, fmt.Errorf("not a 20-byte hex address: %s", addr)
	}
	return hex.DecodeString(addr)
}
