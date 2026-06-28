package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"testing"
	"time"
)

func TestTronProviderProbe(t *testing.T) {
	ch, err := TronBase58ToHex("TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t")
	if err != nil {
		t.Fatal(err)
	}
	hh, err := TronBase58ToHex("TJ5usJLLwjwn7Pw3TPbdzreG7dvgKzfQ5y")
	if err != nil {
		t.Fatal(err)
	}
	param := "000000000000000000000000" + hh[2:]
	body, _ := json.Marshal(map[string]any{
		"owner_address":     hh,
		"contract_address":    ch,
		"function_selector": "balanceOf(address)",
		"parameter":         param,
		"visible":           false,
		"block_num":         83201055,
	})

	for _, base := range []string{
		"https://api.trongrid.io",
		"https://tron.drpc.org",
		"https://tron-rpc.publicnode.com",
		"https://rpc.ankr.com/tron",
	} {
		t.Run(base, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"/wallet/triggerconstantcontract", bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Skip(err)
			}
			raw, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Skipf("HTTP %d: %s", resp.StatusCode, string(raw))
			}
			var out struct {
				ConstantResult []string `json:"constant_result"`
				Result         struct {
					Code string `json:"code"`
				} `json:"result"`
			}
			if err := json.Unmarshal(raw, &out); err != nil {
				t.Fatal(err)
			}
			if len(out.ConstantResult) == 0 {
				t.Skipf("empty constant_result: %s", string(raw))
			}
			bal, ok := new(big.Int).SetString(out.ConstantResult[0], 16)
			if !ok {
				t.Fatalf("bad result %q", out.ConstantResult[0])
			}
			t.Logf("balance=%d code=%s", bal, out.Result.Code)
		})
	}
}

// TestTronBlockNumIgnored documents whether block_num affects TRC20 balance reads.
func TestTronBlockNumIgnored(t *testing.T) {
	c := New()
	contract := "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"
	holder := "TJ5usJLLwjwn7Pw3TPbdzreG7dvgKzfQ5y"
	ctx := context.Background()
	atSnap, _, err := c.TRC20BalanceOf(ctx, contract, holder, 83201055)
	if err != nil {
		t.Skip(err)
	}
	atZero, _, err := c.TRC20BalanceOf(ctx, contract, holder, 0)
	if err != nil {
		t.Skip(err)
	}
	if atSnap.Cmp(atZero) == 0 {
		t.Logf("block_num likely ignored (snap=%s latest=%s)", atSnap, atZero)
	} else {
		t.Logf("block_num may work (snap=%s latest=%s)", atSnap, atZero)
	}
}
