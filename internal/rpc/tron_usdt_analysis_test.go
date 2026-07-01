package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"testing"
	"time"
)

// PR01JUN26 USDT|TRX HotCold rows (address, claim, height).
var usdtTrxRows = []struct {
	addr  string
	claim string
}{
	{"TAzsQ9Gx8eqFNFSKbeXrbi45CuVPHzA8wr", "45331338.13"},
	{"TCLgK89AnXbC9rewvhNb9UgXCc2qJJpBXh", "57078438.35"},
	{"TDqSquXBgUCLYvYC4XZgrprLK589dkhSCf", "407520197.9"},
	{"TJ5usJLLwjwn7Pw3TPbdzreG7dvgKzfQ5y", "49025134.17"},
	{"TJDENsfBJs4RFETt1X1W8wMDc8M5XnJhCe", "49185201.95"},
	{"TJqwA7SoZnERE4zW5uDEiPkbz4B66h9TFj", "37833490.54"},
	{"TK4ykR48cQQoyFcZ5N4xZCbsBaHcg6n3gJ", "36817881.03"},
	{"TNXoiAJ3dct8Fjg4M9fkLFh9S2v9TXc32G", "47573015.43"},
	{"TPtW5TEHhouj6KGshVu5ZQSKZA48QPBnXG", "3084166.245"},
	{"TCEn8ogRSiqdqv26UhsJmQQemrgJS56ZBD", "3054944.819"},
	{"TMuA6YqfCeX8EhbfYEg5y7S4DqzSJireY9", "3.1"},
	{"TMwf7KT8CCdUKuZfKNPTTjbYkFb3eGRbzY", "0.13"},
	{"TQq26fUorctUZvrAgKg8Wz6QyYHvYd6xWK", "12296078.18"},
	{"TQrY8tryqsYVCYS3MFbtffiPp2ccyn4STm", "31361706.05"},
	{"TVGDpgtCs45PJE7ZMHhiC76L3v77qAwJW9", "47034307.8"},
	{"TWd4WrZ9wn84f5x1hZhL4DHvk738ns5jwb", "1857436958"},
	{"TYASr5UV6HEcXatwdFQfmLVUqQQQMUxHLS", "53860294.75"},
}

const (
	usdtTrxContract = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"
	usdtTrxHeight   = int64(83201055)
)

func tronBalanceAt(ctx context.Context, base, holder string, block int64) (*big.Int, error) {
	ch, err := TronBase58ToHex(usdtTrxContract)
	if err != nil {
		return nil, err
	}
	hh, err := TronBase58ToHex(holder)
	if err != nil {
		return nil, err
	}
	param := "000000000000000000000000" + hh[2:]
	body, _ := json.Marshal(map[string]any{
		"owner_address":     hh,
		"contract_address":  ch,
		"function_selector": "balanceOf(address)",
		"parameter":         param,
		"visible":           false,
		"block_num":         block,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"/wallet/triggerconstantcontract", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	raw, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(raw))
	}
	var out struct {
		ConstantResult []string `json:"constant_result"`
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if len(out.ConstantResult) == 0 {
		return nil, fmt.Errorf("empty result")
	}
	bi := new(big.Int)
	if _, ok := bi.SetString(out.ConstantResult[0], 16); !ok {
		return nil, fmt.Errorf("bad hex %q", out.ConstantResult[0])
	}
	return bi, nil
}

func usdt6(bi *big.Int) string {
	if bi == nil {
		return "ERR"
	}
	r := new(big.Rat).SetInt(bi)
	r.Quo(r, big.NewRat(1_000_000, 1))
	return r.FloatString(6)
}

func TestUSDTTrxBalanceAnalysis(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	bases := map[string]string{
		"trongrid":   "https://api.trongrid.io",
		"publicnode": "https://tron-rpc.publicnode.com",
	}
	heights := []int64{usdtTrxHeight - 1, usdtTrxHeight, usdtTrxHeight + 1, 0}

	t.Logf("USDT|TRX analysis @ snapshot block %d (block_num may be ignored)", usdtTrxHeight)
	for _, row := range usdtTrxRows {
		t.Logf("--- %s claim=%s", row.addr, row.claim)
		for name, base := range bases {
			for _, h := range heights {
				bal, err := tronBalanceAt(ctx, base, row.addr, h)
				if err != nil {
					t.Logf("  %s H=%d ERR %v", name, h, err)
					continue
				}
				t.Logf("  %s H=%d bal=%s", name, h, usdt6(bal))
			}
		}
	}
}
