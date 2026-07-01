package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
)

const filecoinGlifV0 = "https://api.node.glif.io/rpc/v0"

// FilecoinActorBalance returns native FIL balance in attoFIL via StateGetActor.
// Public Glif nodes expose limited historical tipsets; height is best-effort.
func (c *LedgerClient) FilecoinActorBalance(ctx context.Context, address string, height int64) (*big.Int, string, error) {
	if c.cache != nil {
		if raw, ok := c.cacheGet("filecoin", "StateGetActor", address, height); ok {
			if v, ok := new(big.Int).SetString(string(raw), 10); ok {
				return v, "cache", nil
			}
		}
	}
	body, _ := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"method":  "Filecoin.StateGetActor",
		"params":  []any{address, nil},
		"id":      1,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, filecoinGlifV0, bytes.NewReader(body))
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "ardmere/0.1")
	resp, err := c.httpc.Do(req)
	if err != nil {
		return nil, "", err
	}
	raw, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, filecoinGlifV0, fmt.Errorf("filecoin HTTP %d", resp.StatusCode)
	}
	var out struct {
		Result struct {
			Balance string `json:"Balance"`
		} `json:"result"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, filecoinGlifV0, err
	}
	if out.Error != nil {
		return nil, filecoinGlifV0, fmt.Errorf("filecoin StateGetActor: %s", out.Error.Message)
	}
	bal, ok := new(big.Int).SetString(strings.TrimSpace(out.Result.Balance), 10)
	if !ok {
		return nil, filecoinGlifV0, fmt.Errorf("bad filecoin balance %q", out.Result.Balance)
	}
	if c.cache != nil {
		c.cachePut("filecoin", "StateGetActor", address, height, []byte(bal.String()))
	}
	return bal, filecoinGlifV0, nil
}
