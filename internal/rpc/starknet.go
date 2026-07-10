package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	starknetETHContract = "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"
	// balanceOf(address) selector on Starknet (not EVM keccak).
	starknetERC20BalanceOfSelector = "0x2e4263afad30923c891518314c3c95dbe830a16874e8abc5777a9a20b54c76e"
)

// StarknetERC20BalanceAtBlock returns ERC20 balance at block number via starknet_call.
func (c *LedgerClient) StarknetERC20BalanceAtBlock(ctx context.Context, token, holder string, blockNum int64) (*big.Int, string, error) {
	cacheChain := "starknet:" + token
	if raw, ok := c.cacheGet(cacheChain, "balance", holder, blockNum); ok {
		v, _ := new(big.Int).SetString(string(raw), 10)
		return v, "cache", nil
	}
	tokenAddr, err := normalizeStarknetFelt(token)
	if err != nil {
		return nil, "", err
	}
	holderAddr, err := normalizeStarknetFelt(holder)
	if err != nil {
		return nil, "", err
	}
	payload := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "starknet_call",
		"params": map[string]any{
			"request": map[string]any{
				"contract_address":     tokenAddr,
				"entry_point_selector": starknetERC20BalanceOfSelector,
				"calldata":             []string{holderAddr},
			},
			"block_id": map[string]any{"block_number": blockNum},
		},
	}

	var lastErr error
	for _, endpoint := range starknetEndpoints() {
		body, used, err := c.httpPostJSON(ctx, endpoint, payload)
		if err != nil {
			lastErr = err
			continue
		}
		low, err := parseStarknetBalanceResult(body)
		if err != nil {
			lastErr = fmt.Errorf("%s: %w", used, err)
			continue
		}
		c.cachePut(cacheChain, "balance", holder, blockNum, []byte(low.String()))
		return low, used, nil
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("no usable starknet provider")
	}
	return nil, "", lastErr
}

// StarknetNativeETHAtBlock returns native ETH (WETH contract) balance on Starknet.
func (c *LedgerClient) StarknetNativeETHAtBlock(ctx context.Context, holder string, blockNum int64) (*big.Int, string, error) {
	return c.StarknetERC20BalanceAtBlock(ctx, starknetETHContract, holder, blockNum)
}

func parseStarknetBalanceResult(body []byte) (*big.Int, error) {
	var out struct {
		Result []string `json:"result"`
		Error  *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	if out.Error != nil {
		return nil, fmt.Errorf("starknet: %s", out.Error.Message)
	}
	if len(out.Result) == 0 {
		return big.NewInt(0), nil
	}
	low, ok := new(big.Int).SetString(strings.TrimPrefix(out.Result[0], "0x"), 16)
	if !ok {
		return nil, fmt.Errorf("starknet balance decode")
	}
	if len(out.Result) > 1 {
		if hi, ok := new(big.Int).SetString(strings.TrimPrefix(out.Result[1], "0x"), 16); ok {
			low.Add(low, new(big.Int).Lsh(hi, 128))
		}
	}
	return low, nil
}

func normalizeStarknetFelt(addr string) (string, error) {
	addr = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(addr)), "0x")
	if addr == "" {
		return "", fmt.Errorf("empty starknet address")
	}
	if len(addr) > 64 {
		return "", fmt.Errorf("starknet address too long")
	}
	return "0x" + strings.Repeat("0", 64-len(addr)) + addr, nil
}

func starknetEndpoints() []string {
	if u := os.Getenv("STARKNET_RPC"); u != "" {
		return []string{u}
	}
	cfg, err := LoadProviderConfig()
	if err != nil {
		cfg = DefaultProviderConfig()
	}
	if providers, ok := cfg[Network("STARKNET")]; ok {
		out := make([]string, 0, len(providers))
		for _, p := range providers {
			if u, ok := expandProviderURL(p.URL); ok {
				out = append(out, u)
			}
		}
		if len(out) > 0 {
			return out
		}
	}
	return []string{
		"https://starknet-mainnet.public.blastapi.io/rpc/v0_9",
		"https://rpc.starknet.lava.build/rpc/v0_9",
	}
}

// TonJettonBalance returns jetton balance via tonapi (live; no historical seqno on public API).
func (c *LedgerClient) TonJettonBalance(ctx context.Context, owner, jettonMaster string, seqno int64) (int64, string, error) {
	cacheChain := "tonapi:" + jettonMaster
	if raw, ok := c.cacheGet(cacheChain, "balance", owner, seqno); ok {
		v, _ := strconv.ParseInt(string(raw), 10, 64)
		return v, "cache", nil
	}
	base := "https://tonapi.io/v2"
	if u := os.Getenv("TONAPI_URL"); u != "" {
		base = strings.TrimRight(u, "/")
	}
	u := fmt.Sprintf("%s/accounts/%s/jettons/%s", base, url.PathEscape(owner), url.PathEscape(jettonMaster))
	body, used, err := c.httpGet(ctx, u)
	if err != nil {
		return 0, used, err
	}
	var out struct {
		Balance string `json:"balance"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return 0, used, err
	}
	v, err := strconv.ParseInt(out.Balance, 10, 64)
	if err != nil {
		return 0, used, err
	}
	c.cachePut(cacheChain, "balance", owner, seqno, []byte(strconv.FormatInt(v, 10)))
	return v, used, nil
}

// TonNativeBalance returns TON balance in nanotons via tonapi (live; no public historical seqno).
func (c *LedgerClient) TonNativeBalance(ctx context.Context, address string, seqno int64) (int64, string, error) {
	cacheChain := "tonapi:native"
	if raw, ok := c.cacheGet(cacheChain, "balance", address, seqno); ok {
		v, _ := strconv.ParseInt(string(raw), 10, 64)
		return v, "cache", nil
	}
	base := "https://tonapi.io/v2"
	if u := os.Getenv("TONAPI_URL"); u != "" {
		base = strings.TrimRight(u, "/")
	}
	u := fmt.Sprintf("%s/accounts/%s", base, url.PathEscape(address))
	body, used, err := c.httpGet(ctx, u)
	if err != nil {
		return 0, used, err
	}
	var out struct {
		Balance int64 `json:"balance"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return 0, used, err
	}
	c.cachePut(cacheChain, "balance", address, seqno, []byte(strconv.FormatInt(out.Balance, 10)))
	return out.Balance, used, nil
}
