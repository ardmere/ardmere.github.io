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

// StarknetERC20BalanceAtBlock returns ERC20 balance at block number via starknet_call.
func (c *LedgerClient) StarknetERC20BalanceAtBlock(ctx context.Context, token, holder string, blockNum int64) (*big.Int, string, error) {
	cacheChain := "starknet:" + token
	if raw, ok := c.cacheGet(cacheChain, "balance", holder, blockNum); ok {
		v, _ := new(big.Int).SetString(string(raw), 10)
		return v, "cache", nil
	}
	endpoint := starknetRPC()
	selector := "0x2e4263af247b6bd0575f54389450d502e7800d903a585d0f31506541755fda2" // balanceOf selector
	calldata := []string{holder}
	body, used, err := c.httpPostJSON(ctx, endpoint, map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "starknet_call",
		"params": map[string]any{
			"request": map[string]any{
				"contract_address": token,
				"entry_point_selector": selector,
				"calldata": calldata,
			},
			"block_id": map[string]any{"block_number": blockNum},
		},
	})
	if err != nil {
		return nil, used, err
	}
	var out struct {
		Result []string `json:"result"`
		Error  *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, used, err
	}
	if out.Error != nil {
		return nil, used, fmt.Errorf("starknet: %s", out.Error.Message)
	}
	if len(out.Result) == 0 {
		return big.NewInt(0), used, nil
	}
	low, ok := new(big.Int).SetString(strings.TrimPrefix(out.Result[0], "0x"), 16)
	if !ok {
		return nil, used, fmt.Errorf("starknet balance decode")
	}
	if len(out.Result) > 1 {
		if hi, ok := new(big.Int).SetString(strings.TrimPrefix(out.Result[1], "0x"), 16); ok {
			low.Add(low, new(big.Int).Lsh(hi, 128))
		}
	}
	c.cachePut(cacheChain, "balance", holder, blockNum, []byte(low.String()))
	return low, used, nil
}

// StarknetNativeETHAtBlock returns native ETH (WETH contract) balance on Starknet.
func (c *LedgerClient) StarknetNativeETHAtBlock(ctx context.Context, holder string, blockNum int64) (*big.Int, string, error) {
	const weth = "0x049d36570d4e46f48e996dafd94c4525b4747612bb"
	return c.StarknetERC20BalanceAtBlock(ctx, weth, holder, blockNum)
}

func starknetRPC() string {
	if u := os.Getenv("STARKNET_RPC"); u != "" {
		return u
	}
	return "https://starknet-mainnet.public.blastapi.io/rpc/v0_7"
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
	u := fmt.Sprintf("%s/accounts/%s/jettons/%s", base, owner, jettonMaster)
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
