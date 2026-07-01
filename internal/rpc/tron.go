package rpc

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/mr-tron/base58"
)

const NetTron Network = "TRX"

func tronWalletBase(url string) string {
	u := strings.TrimRight(url, "/")
	return strings.TrimSuffix(u, "/jsonrpc")
}

// TronNativeBalance returns TRX balance in SUN (6 decimals) via wallet/getaccount.
// Public Tron nodes expose limited historical native balance; blockNum is best-effort.
func (c *Client) TronNativeBalance(ctx context.Context, holderBase58 string, blockNum int64) (*big.Int, string, error) {
	if blockNum > 0 {
		if bal, used, err := c.tronNativeBalanceAtBlock(ctx, holderBase58, blockNum); err == nil {
			return bal, used, nil
		}
	}
	holderHex, err := TronBase58ToHex(holderBase58)
	if err != nil {
		return nil, "", fmt.Errorf("holder address: %w", err)
	}
	body, _ := json.Marshal(map[string]any{
		"address": holderHex,
		"visible": false,
	})
	return c.tronGetAccountBalance(ctx, body)
}

func (c *Client) tronNativeBalanceAtBlock(ctx context.Context, holderBase58 string, blockNum int64) (*big.Int, string, error) {
	hash, used, err := c.tronBlockHash(ctx, blockNum)
	if err != nil {
		return nil, used, err
	}
	body, _ := json.Marshal(map[string]any{
		"account_identifier": map[string]string{"address": holderBase58},
		"block_identifier":   map[string]any{"hash": hash, "number": blockNum},
		"visible":            true,
	})
	providers := ProvidersFor(c.providers, NetTron, true)
	var lastErr error
	for _, p := range providers {
		u := tronWalletBase(p.URL)
		if c.isDisabled(u) {
			continue
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u+"/wallet/getaccountbalance", bytes.NewReader(body))
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
		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("%s -> HTTP %d", u, resp.StatusCode)
			continue
		}
		var out struct {
			Balance int64  `json:"balance"`
			Error   string `json:"Error"`
		}
		if err := json.Unmarshal(raw, &out); err != nil {
			lastErr = err
			continue
		}
		if out.Error != "" {
			lastErr = fmt.Errorf("%s -> %s", u, out.Error)
			continue
		}
		return big.NewInt(out.Balance), u, nil
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("no usable tron provider for getaccountbalance")
	}
	return nil, "", lastErr
}

func (c *Client) tronBlockHash(ctx context.Context, blockNum int64) (string, string, error) {
	body, _ := json.Marshal(map[string]any{"num": blockNum})
	providers := ProvidersFor(c.providers, NetTron, true)
	var lastErr error
	for _, p := range providers {
		u := tronWalletBase(p.URL)
		if c.isDisabled(u) {
			continue
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u+"/wallet/getblockbynum", bytes.NewReader(body))
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
		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("%s -> HTTP %d", u, resp.StatusCode)
			continue
		}
		var out struct {
			BlockID string `json:"blockID"`
		}
		if err := json.Unmarshal(raw, &out); err != nil {
			lastErr = err
			continue
		}
		if out.BlockID == "" {
			lastErr = fmt.Errorf("%s -> empty blockID", u)
			continue
		}
		return out.BlockID, u, nil
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("no usable tron provider for getblockbynum")
	}
	return "", "", lastErr
}

func (c *Client) tronGetAccountBalance(ctx context.Context, body []byte) (*big.Int, string, error) {
	providers := ProvidersFor(c.providers, NetTron, false)
	var lastErr error
	for _, p := range providers {
		u := tronWalletBase(p.URL)
		if c.isDisabled(u) {
			continue
		}
		if p.RateLimitMs > 0 {
			time.Sleep(time.Duration(p.RateLimitMs) * time.Millisecond)
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u+"/wallet/getaccount", bytes.NewReader(body))
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
		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("%s -> HTTP %d", u, resp.StatusCode)
			continue
		}
		var out struct {
			Balance int64 `json:"balance"`
		}
		if err := json.Unmarshal(raw, &out); err != nil {
			lastErr = err
			continue
		}
		return big.NewInt(out.Balance), u, nil
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("no usable tron provider")
	}
	return nil, "", lastErr
}

// TRC20BalanceOf queries balanceOf(address) for a TRC20 contract at a historical block.
func (c *Client) TRC20BalanceOf(ctx context.Context, contractBase58, holderBase58 string, blockNum int64) (*big.Int, string, error) {
	contractHex, err := TronBase58ToHex(contractBase58)
	if err != nil {
		return nil, "", fmt.Errorf("contract address: %w", err)
	}
	holderHex, err := TronBase58ToHex(holderBase58)
	if err != nil {
		return nil, "", fmt.Errorf("holder address: %w", err)
	}
	if !strings.HasPrefix(holderHex, "41") || len(holderHex) != 42 {
		return nil, "", fmt.Errorf("bad holder hex: %s", holderHex)
	}
	param := strings.Repeat("0", 24) + holderHex[2:]

	providers := ProvidersFor(c.providers, NetTron, blockNum > 0)
	if len(providers) == 0 {
		return nil, "", fmt.Errorf("network %s not configured", NetTron)
	}

	body, _ := json.Marshal(map[string]any{
		"owner_address":    holderHex,
		"contract_address": contractHex,
		"function_selector": "balanceOf(address)",
		"parameter":         param,
		"visible":           false,
		"block_num":         blockNum,
	})

	var lastErr error
	for _, p := range providers {
		u := strings.TrimRight(p.URL, "/")
		if c.isDisabled(u) {
			continue
		}
		if p.RateLimitMs > 0 {
			time.Sleep(time.Duration(p.RateLimitMs) * time.Millisecond)
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u+"/wallet/triggerconstantcontract", bytes.NewReader(body))
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
		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("%s -> HTTP %d", u, resp.StatusCode)
			continue
		}

		var out struct {
			Result struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"result"`
			ConstantResult []string `json:"constant_result"`
		}
		if err := json.Unmarshal(raw, &out); err != nil {
			lastErr = err
			continue
		}
		if out.Result.Code != "" && out.Result.Code != "SUCCESS" {
			lastErr = fmt.Errorf("%s -> %s: %s", u, out.Result.Code, out.Result.Message)
			continue
		}
		if len(out.ConstantResult) == 0 || out.ConstantResult[0] == "" {
			lastErr = fmt.Errorf("%s -> empty constant_result", u)
			continue
		}
		bi := new(big.Int)
		if _, ok := bi.SetString(out.ConstantResult[0], 16); !ok {
			lastErr = fmt.Errorf("%s -> bad constant_result %q", u, out.ConstantResult[0])
			continue
		}
		return bi, u, nil
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("no usable tron provider")
	}
	return nil, "", lastErr
}

// TronBase58ToHex decodes a Tron base58check address to 41-prefixed hex.
func TronBase58ToHex(addr string) (string, error) {
	raw, err := base58.Decode(addr)
	if err != nil {
		return "", err
	}
	if len(raw) < 5 {
		return "", fmt.Errorf("base58 too short")
	}
	payload := raw[:len(raw)-4]
	checksum := raw[len(raw)-4:]
	h := sha256.Sum256(payload)
	h2 := sha256.Sum256(h[:])
	if !bytes.Equal(checksum, h2[:4]) {
		return "", fmt.Errorf("bad base58 checksum")
	}
	if len(payload) != 21 || payload[0] != 0x41 {
		return "", fmt.Errorf("unexpected tron payload len=%d", len(payload))
	}
	return hex.EncodeToString(payload), nil
}
