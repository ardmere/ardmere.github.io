package rpc

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/vedhavyas/go-subkey/v2"
)

func (c *LedgerClient) nearRPC() string {
	if u := os.Getenv("NEAR_RPC"); u != "" {
		return u
	}
	return "https://archival-rpc.mainnet.fastnear.com"
}

func (c *LedgerClient) aptosRPC() string {
	if u := os.Getenv("APTOS_RPC"); u != "" {
		return u
	}
	return "https://fullnode.mainnet.aptoslabs.com/v1"
}

func (c *LedgerClient) suiRPC() string {
	if u := os.Getenv("SUI_RPC"); u != "" {
		return u
	}
	return "https://fullnode.mainnet.sui.io:443"
}

func (c *LedgerClient) httpGet(ctx context.Context, u string) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, u, err
	}
	req.Header.Set("User-Agent", "ardmere/0.1")
	resp, err := c.httpc.Do(req)
	if err != nil {
		return nil, u, err
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, u, fmt.Errorf("HTTP %d: %s", resp.StatusCode, truncate(string(body), 120))
	}
	return body, u, nil
}

func (c *LedgerClient) httpPostJSON(ctx context.Context, endpoint string, payload any) ([]byte, string, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, endpoint, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(string(raw)))
	if err != nil {
		return nil, endpoint, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "ardmere/0.1")
	resp, err := c.httpc.Do(req)
	if err != nil {
		return nil, endpoint, err
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, endpoint, fmt.Errorf("HTTP %d: %s", resp.StatusCode, truncate(string(body), 120))
	}
	return body, endpoint, nil
}

// AptosCoinBalanceAtVersion returns native APT in CoinStore octas at ledger version.
// Missing CoinStore (HTTP 404) is reported as zero balance, not an error.
func (c *LedgerClient) AptosCoinBalanceAtVersion(ctx context.Context, addr string, version int64) (int64, string, error) {
	v, used, missing, err := c.aptosCoinStoreOctas(ctx, addr, version)
	if err != nil {
		return 0, used, err
	}
	if missing {
		return 0, used, nil
	}
	return v, used, nil
}

const aptosNativeFAMetadata = "0xa"

// AptosAccountedAPTAtVersion sums legacy CoinStore APT and fungible-asset native APT (0xa).
func (c *LedgerClient) AptosAccountedAPTAtVersion(ctx context.Context, addr string, version int64) (total int64, used string, components map[string]string, err error) {
	components = map[string]string{"mode": "aptos_accounted"}
	coin, u1, missing, err := c.aptosCoinStoreOctas(ctx, addr, version)
	if err != nil {
		return 0, u1, nil, err
	}
	fa, u2, err := c.AptosFABalanceAtVersion(ctx, addr, aptosNativeFAMetadata, version)
	if err != nil {
		return 0, u2, nil, err
	}
	used = u1
	if u2 != "" && u2 != used {
		if used == "" {
			used = u2
		} else if u2 != "cache" {
			used = u2
		}
	}
	components["coin_store"] = strconv.FormatInt(coin, 10)
	components["fa_native"] = strconv.FormatInt(fa, 10)
	if missing {
		components["coin_store_missing"] = "true"
	}
	total = coin + fa
	return total, used, components, nil
}

func (c *LedgerClient) aptosCoinStoreOctas(ctx context.Context, addr string, version int64) (octas int64, used string, missing bool, err error) {
	cacheChain := "aptos:coin"
	if raw, ok := c.cacheGet(cacheChain, "balance", addr, version); ok {
		if string(raw) == "missing" {
			return 0, "cache", true, nil
		}
		v, _ := strconv.ParseInt(string(raw), 10, 64)
		return v, "cache", false, nil
	}
	base := strings.TrimRight(c.aptosRPC(), "/")
	path := fmt.Sprintf("%s/accounts/%s/resource/0x1::coin::CoinStore%%3C0x1::aptos_coin::AptosCoin%%3E?ledger_version=%d",
		base, url.PathEscape(addr), version)
	body, used, err := c.httpGet(ctx, path)
	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			c.cachePut(cacheChain, "balance", addr, version, []byte("missing"))
			return 0, used, true, nil
		}
		return 0, used, false, err
	}
	var out struct {
		Data struct {
			Coin struct {
				Value string `json:"value"`
			} `json:"coin"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return 0, used, false, err
	}
	v, err := strconv.ParseInt(out.Data.Coin.Value, 10, 64)
	if err != nil {
		return 0, used, false, fmt.Errorf("aptos coin value: %w", err)
	}
	c.cachePut(cacheChain, "balance", addr, version, []byte(strconv.FormatInt(v, 10)))
	return v, used, false, nil
}

// parseAptosFABalanceBody handles Aptos /balance responses: bare JSON number or {"balance":"..."}.
func parseAptosFABalanceBody(body []byte) (int64, error) {
	body = bytesTrimSpace(body)
	if len(body) == 0 {
		return 0, fmt.Errorf("aptos fa balance: empty body")
	}
	if body[0] == '{' {
		var out struct {
			Balance string `json:"balance"`
		}
		if err := json.Unmarshal(body, &out); err != nil {
			return 0, err
		}
		return strconv.ParseInt(out.Balance, 10, 64)
	}
	var n json.Number
	if err := json.Unmarshal(body, &n); err != nil {
		return 0, err
	}
	return n.Int64()
}

// AptosFABalanceAtVersion returns fungible-asset balance (metadata address) at ledger version.
func (c *LedgerClient) AptosFABalanceAtVersion(ctx context.Context, addr, faMetadata string, version int64) (int64, string, error) {
	cacheChain := "aptos:fa:" + faMetadata
	if raw, ok := c.cacheGet(cacheChain, "balance", addr, version); ok {
		v, _ := strconv.ParseInt(string(raw), 10, 64)
		return v, "cache", nil
	}
	base := strings.TrimRight(c.aptosRPC(), "/")
	path := fmt.Sprintf("%s/accounts/%s/balance/%s?ledger_version=%d",
		base, url.PathEscape(addr), url.PathEscape(faMetadata), version)
	body, used, err := c.httpGet(ctx, path)
	if err != nil {
		return 0, used, err
	}
	v, err := parseAptosFABalanceBody(body)
	if err != nil {
		return 0, used, fmt.Errorf("aptos fa balance: %w", err)
	}
	c.cachePut(cacheChain, "balance", addr, version, []byte(strconv.FormatInt(v, 10)))
	return v, used, nil
}

// NearNativeBalanceAtHeight queries NEAR yoctoNEAR at block height; falls back to live on archival miss.
func (c *LedgerClient) NearNativeBalanceAtHeight(ctx context.Context, account string, height int64) (int64, string, bool, error) {
	cacheChain := "near:native"
	if raw, ok := c.cacheGet(cacheChain, "balance", account, height); ok {
		parts := strings.Split(string(raw), "|")
		if len(parts) == 2 {
			v, _ := strconv.ParseInt(parts[0], 10, 64)
			return v, "cache", parts[1] == "live", nil
		}
	}
	endpoints := []string{c.nearRPC()}
	if extra := os.Getenv("NEAR_RPC_FALLBACK"); extra != "" {
		endpoints = append(endpoints, strings.Split(extra, ",")...)
	}
	endpoints = append(endpoints, "https://archival-rpc.mainnet.near.org", "https://rpc.mainnet.near.org")

	var lastErr error
	for _, ep := range endpoints {
		ep = strings.TrimSpace(ep)
		if ep == "" {
			continue
		}
		amount, live, err := c.nearViewAccount(ctx, ep, account, height, !strings.Contains(ep, "archival"))
		if err == nil {
			flag := "hist"
			if live {
				flag = "live"
			}
			c.cachePut(cacheChain, "balance", account, height, []byte(fmt.Sprintf("%d|%s", amount, flag)))
			return amount, ep, live, nil
		}
		lastErr = err
		if strings.Contains(err.Error(), "GARBAGE_COLLECTED") {
			continue
		}
	}
	return 0, "", false, lastErr
}

func (c *LedgerClient) nearViewAccount(ctx context.Context, endpoint, account string, height int64, allowLiveFallback bool) (int64, bool, error) {
	query := func(blockID any, live bool) (int64, error) {
		params := map[string]any{
			"request_type": "view_account",
			"account_id":   account,
		}
		if blockID != nil {
			params["finality"] = "none"
			params["block_id"] = blockID
		} else {
			params["finality"] = "optimistic"
		}
		body, _, err := c.httpPostJSON(ctx, endpoint, map[string]any{
			"jsonrpc": "2.0", "id": 1, "method": "query", "params": params,
		})
		if err != nil {
			return 0, err
		}
		var out struct {
			Result struct {
				Amount string `json:"amount"`
			} `json:"result"`
			Error *struct {
				Data string `json:"data"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &out); err != nil {
			return 0, err
		}
		if out.Error != nil {
			return 0, fmt.Errorf("near: %s", out.Error.Data)
		}
		return strconv.ParseInt(out.Result.Amount, 10, 64)
	}
	if v, err := query(height, false); err == nil {
		return v, false, nil
	} else if allowLiveFallback {
		if v, err := query(nil, true); err == nil {
			return v, true, nil
		} else {
			return 0, false, err
		}
	}
	return 0, false, fmt.Errorf("near historical query failed at %d", height)
}

// NearFTBalanceAtHeight returns NEP-141 ft balance at block height (live fallback on archival miss).
func (c *LedgerClient) NearFTBalanceAtHeight(ctx context.Context, ftContract, account string, height int64) (int64, string, bool, error) {
	cacheChain := "near:ft:" + ftContract
	if raw, ok := c.cacheGet(cacheChain, "balance", account, height); ok {
		parts := strings.Split(string(raw), "|")
		if len(parts) == 2 {
			v, _ := strconv.ParseInt(parts[0], 10, 64)
			return v, "cache", parts[1] == "live", nil
		}
	}
	args, _ := json.Marshal(map[string]string{"account_id": account})
	argsB64 := base64.StdEncoding.EncodeToString(args)
	try := func(blockID any, live bool, endpoint string) (int64, bool, error) {
		params := map[string]any{
			"request_type": "call_function",
			"finality":     "none",
			"account_id":   ftContract,
			"method_name":  "ft_balance_of",
			"args_base64":  argsB64,
		}
		if blockID != nil {
			params["block_id"] = blockID
		} else {
			params["finality"] = "optimistic"
		}
		body, _, err := c.httpPostJSON(ctx, endpoint, map[string]any{
			"jsonrpc": "2.0", "id": 1, "method": "query", "params": params,
		})
		if err != nil {
			return 0, live, err
		}
		var out struct {
			Result struct {
				Result []byte `json:"result"`
			} `json:"result"`
			Error *struct{ Data string `json:"data"` } `json:"error"`
		}
		if err := json.Unmarshal(body, &out); err != nil {
			return 0, live, err
		}
		if out.Error != nil {
			return 0, live, fmt.Errorf("near ft: %s", out.Error.Data)
		}
		var balStr string
		if err := json.Unmarshal(out.Result.Result, &balStr); err != nil {
			return 0, live, err
		}
		v, err := strconv.ParseInt(strings.Trim(balStr, `"`), 10, 64)
		return v, live, err
	}

	endpoint := c.nearRPC()
	if v, live, err := try(height, false, endpoint); err == nil {
		flag := "hist"
		if live {
			flag = "live"
		}
		c.cachePut(cacheChain, "balance", account, height, []byte(fmt.Sprintf("%d|%s", v, flag)))
		return v, endpoint, live, nil
	} else if v, _, err2 := try(nil, true, "https://rpc.mainnet.near.org"); err2 == nil {
		c.cachePut(cacheChain, "balance", account, height, []byte(fmt.Sprintf("%d|live", v)))
		return v, "https://rpc.mainnet.near.org", true, nil
	} else {
		return 0, endpoint, false, err
	}
}

// SuiBalanceAtCheckpoint returns MIST balance for coinType at checkpoint (0x2::sui::SUI when coinType empty).
func (c *LedgerClient) SuiBalanceAtCheckpoint(ctx context.Context, addr, coinType string, checkpoint int64) (uint64, string, error) {
	if coinType == "" {
		coinType = "0x2::sui::SUI"
	}
	cacheChain := "sui:" + coinType
	if raw, ok := c.cacheGet(cacheChain, "balance", addr, checkpoint); ok {
		v, _ := strconv.ParseUint(string(raw), 10, 64)
		return v, "cache", nil
	}
	params := []any{addr, coinType, strconv.FormatInt(checkpoint, 10)}
	body, used, err := c.httpPostJSON(ctx, c.suiRPC(), map[string]any{
		"jsonrpc": "2.0", "id": 1, "method": "suix_getBalance", "params": params,
	})
	if err != nil {
		return 0, used, err
	}
	var out struct {
		Result struct {
			TotalBalance string `json:"totalBalance"`
		} `json:"result"`
		Error *struct{ Message string `json:"message"` } `json:"error"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return 0, used, err
	}
	if out.Error != nil {
		return 0, used, fmt.Errorf("sui: %s", out.Error.Message)
	}
	v, err := strconv.ParseUint(out.Result.TotalBalance, 10, 64)
	if err != nil {
		return 0, used, err
	}
	c.cachePut(cacheChain, "balance", addr, checkpoint, []byte(strconv.FormatUint(v, 10)))
	return v, used, nil
}

// HbarBalanceAtBlock returns tinybar balance at mirror-node block number.
func (c *LedgerClient) HbarBalanceAtBlock(ctx context.Context, account string, blockNum int64) (int64, string, error) {
	cacheChain := "hbar"
	if raw, ok := c.cacheGet(cacheChain, "balance", account, blockNum); ok {
		v, _ := strconv.ParseInt(string(raw), 10, 64)
		return v, "cache", nil
	}
	mirror := "https://mainnet-public.mirrornode.hedera.com/api/v1"
	if u := os.Getenv("HBAR_MIRROR"); u != "" {
		mirror = strings.TrimRight(u, "/")
	}
	blkBody, used, err := c.httpGet(ctx, fmt.Sprintf("%s/blocks/%d", mirror, blockNum))
	if err != nil {
		return 0, used, err
	}
	var blk struct {
		Timestamp struct {
			From string `json:"from"`
		} `json:"timestamp"`
	}
	if err := json.Unmarshal(blkBody, &blk); err != nil {
		return 0, used, err
	}
	ts := blk.Timestamp.From
	accBody, used2, err := c.httpGet(ctx, fmt.Sprintf("%s/accounts/%s?timestamp=lte:%s", mirror, url.PathEscape(account), url.QueryEscape(ts)))
	if err != nil {
		return 0, used2, err
	}
	var acc struct {
		Balance struct {
			Balance int64 `json:"balance"`
		} `json:"balance"`
	}
	if err := json.Unmarshal(accBody, &acc); err != nil {
		return 0, used2, err
	}
	c.cachePut(cacheChain, "balance", account, blockNum, []byte(strconv.FormatInt(acc.Balance.Balance, 10)))
	return acc.Balance.Balance, used2, nil
}

// HbarHTSBalanceAtBlock returns HTS token balance (token id e.g. 0.0.456858) at block.
func (c *LedgerClient) HbarHTSBalanceAtBlock(ctx context.Context, account, tokenID string, blockNum int64) (int64, string, error) {
	cacheChain := "hbar:hts:" + tokenID
	if raw, ok := c.cacheGet(cacheChain, "balance", account, blockNum); ok {
		v, _ := strconv.ParseInt(string(raw), 10, 64)
		return v, "cache", nil
	}
	mirror := "https://mainnet-public.mirrornode.hedera.com/api/v1"
	if u := os.Getenv("HBAR_MIRROR"); u != "" {
		mirror = strings.TrimRight(u, "/")
	}
	blkBody, used, err := c.httpGet(ctx, fmt.Sprintf("%s/blocks/%d", mirror, blockNum))
	if err != nil {
		return 0, used, err
	}
	var blk struct {
		Timestamp struct {
			From string `json:"from"`
		} `json:"timestamp"`
	}
	if err := json.Unmarshal(blkBody, &blk); err != nil {
		return 0, used, err
	}
	accBody, used2, err := c.httpGet(ctx, fmt.Sprintf("%s/accounts/%s/tokens?timestamp=lte:%s", mirror, url.PathEscape(account), url.QueryEscape(blk.Timestamp.From)))
	if err != nil {
		return 0, used2, err
	}
	var acc struct {
		Tokens []struct {
			TokenID string `json:"token_id"`
			Balance int64  `json:"balance"`
		} `json:"tokens"`
	}
	if err := json.Unmarshal(accBody, &acc); err != nil {
		return 0, used2, err
	}
	for _, t := range acc.Tokens {
		if t.TokenID == tokenID {
			c.cachePut(cacheChain, "balance", account, blockNum, []byte(strconv.FormatInt(t.Balance, 10)))
			return t.Balance, used2, nil
		}
	}
	c.cachePut(cacheChain, "balance", account, blockNum, []byte("0"))
	return 0, used2, nil
}

func (c *LedgerClient) SubstrateBalanceAtBlock(ctx context.Context, rpcURL, ss58Addr string, height int64) (*big.Int, string, error) {
	// custody = free + reserved (reserved holds nomination-pool / lock balances on Asset Hub).
	cacheChain := "substrate:custody:" + rpcURL
	if raw, ok := c.cacheGet(cacheChain, "balance", ss58Addr, height); ok {
		v, _ := new(big.Int).SetString(string(raw), 10)
		return v, "cache", nil
	}
	wsURL := strings.Replace(rpcURL, "https://", "wss://", 1)
	api, err := gsrpc.NewSubstrateAPI(wsURL)
	if err != nil {
		return nil, wsURL, fmt.Errorf("substrate connect: %w", err)
	}
	hash, err := api.RPC.Chain.GetBlockHash(uint64(height))
	if err != nil {
		return nil, wsURL, err
	}
	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, wsURL, err
	}
	_, pubKey, err := subkey.SS58Decode(ss58Addr)
	if err != nil {
		return nil, wsURL, fmt.Errorf("ss58 decode: %w", err)
	}
	var accountID types.AccountID
	if len(pubKey) != 32 {
		return nil, wsURL, fmt.Errorf("unexpected pubkey length %d", len(pubKey))
	}
	copy(accountID[:], pubKey)
	key, err := types.CreateStorageKey(meta, "System", "Account", accountID[:], nil)
	if err != nil {
		return nil, wsURL, err
	}
	var acct types.AccountInfo
	ok, err := api.RPC.State.GetStorage(key, &acct, hash)
	if err != nil {
		return nil, wsURL, err
	}
	if !ok {
		return big.NewInt(0), wsURL, nil
	}
	free, ok := new(big.Int).SetString(acct.Data.Free.String(), 10)
	if !ok {
		return nil, wsURL, fmt.Errorf("substrate free balance decode")
	}
	reserved, ok := new(big.Int).SetString(acct.Data.Reserved.String(), 10)
	if !ok {
		return nil, wsURL, fmt.Errorf("substrate reserved balance decode")
	}
	custody := new(big.Int).Add(free, reserved)
	c.cachePut(cacheChain, "balance", ss58Addr, height, []byte(custody.String()))
	return custody, wsURL, nil
}

func bytesTrimSpace(b []byte) []byte {
	return []byte(strings.TrimSpace(string(b)))
}
