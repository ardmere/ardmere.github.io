package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// Chainlist Tron Mainnet (chainId 728126428) as of 2026-06.
var chainlistTronRPCs = []string{
	"https://tron.drpc.org",
	"https://tron.therpc.io/jsonrpc",
	"https://tron.api.pocket.network",
	"https://rpc.ankr.com/tron_jsonrpc",
	"https://api.trongrid.io/jsonrpc",
	// not on chainlist but working in our config
	"https://api.trongrid.io",
	"https://tron-rpc.publicnode.com",
	"https://rpc.ankr.com/tron",
}

func walletBase(url string) string {
	u := strings.TrimRight(url, "/")
	u = strings.TrimSuffix(u, "/jsonrpc")
	return u
}

func TestChainlistTronNativeProbe(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	holder := "TJ5usJLLwjwn7Pw3TPbdzreG7dvgKzfQ5y"
	hh, err := TronBase58ToHex(holder)
	if err != nil {
		t.Fatal(err)
	}
	ch, err := TronBase58ToHex("TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t")
	if err != nil {
		t.Fatal(err)
	}
	param := "000000000000000000000000" + hh[2:]
	body, _ := json.Marshal(map[string]any{
		"owner_address":     hh,
		"contract_address":  ch,
		"function_selector": "balanceOf(address)",
		"parameter":         param,
		"visible":           false,
		"block_num":         83201055,
	})

	t.Log("probe /wallet/triggerconstantcontract (native API required by verifier)")
	for _, raw := range chainlistTronRPCs {
		base := walletBase(raw)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"/wallet/triggerconstantcontract", bytes.NewReader(body))
		if err != nil {
			t.Logf("SKIP %s: %v", raw, err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Logf("FAIL %s -> %v", raw, err)
			continue
		}
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			msg := string(b)
			if len(msg) > 120 {
				msg = msg[:120]
			}
			t.Logf("FAIL %s -> HTTP %d %s", raw, resp.StatusCode, msg)
			continue
		}
		var out struct {
			ConstantResult []string `json:"constant_result"`
		}
		if err := json.Unmarshal(b, &out); err != nil || len(out.ConstantResult) == 0 {
			t.Logf("FAIL %s -> bad body %s", raw, string(b[:min(120, len(b))]))
			continue
		}
		t.Logf("OK   %s bal_hex=%s", raw, out.ConstantResult[0][:min(16, len(out.ConstantResult[0]))]+"...")
	}
}

func TestChainlistTronJSONRPCProbe(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// balanceOf(address) on USDT TRC20 via eth_call — only if endpoint is JSON-RPC compatible
	holder := "TJ5usJLLwjwn7Pw3TPbdzreG7dvgKzfQ5y"
	hh, err := TronBase58ToHex(holder)
	if err != nil {
		t.Fatal(err)
	}
	// Tron JSON-RPC uses 0x41-prefixed 20-byte address in call data
	addr20 := hh[2:]
	data := "0x70a08231" + strings.Repeat("0", 24) + addr20
	contract := "0x41a614f803b6fd780986a42c78ec9c7f77e6ded13c" // USDT hex

	t.Log("probe eth_call via /jsonrpc (Chainlist style)")
	for _, raw := range chainlistTronRPCs {
		if !strings.Contains(raw, "jsonrpc") && raw != "https://tron.drpc.org" && raw != "https://tron.api.pocket.network" {
			continue
		}
		url := raw
		if !strings.HasSuffix(url, "/jsonrpc") {
			url = strings.TrimRight(url, "/") + "/jsonrpc"
		}
		payload, _ := json.Marshal(map[string]any{
			"jsonrpc": "2.0",
			"id":      1,
			"method":  "eth_call",
			"params": []any{
				map[string]string{"to": contract, "data": data},
				"0x4f3f9bf", // 83201055
			},
		})
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
		if err != nil {
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Logf("FAIL %s -> %v", url, err)
			continue
		}
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		msg := string(b)
		if len(msg) > 160 {
			msg = msg[:160]
		}
		if resp.StatusCode != http.StatusOK {
			t.Logf("FAIL %s -> HTTP %d %s", url, resp.StatusCode, msg)
			continue
		}
		var out struct {
			Result string `json:"result"`
			Error  *struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		json.Unmarshal(b, &out)
		if out.Error != nil {
			t.Logf("FAIL %s -> %s", url, out.Error.Message)
			continue
		}
		if out.Result == "" || out.Result == "0x" {
			t.Logf("FAIL %s -> empty result", url)
			continue
		}
		t.Logf("OK   %s result=%s", url, out.Result[:min(18, len(out.Result))]+"...")
	}
}
