package rpc

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"strconv"
	"strings"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/vedhavyas/go-subkey/v2"
)

func (c *LedgerClient) algodRPC() string {
	if u := os.Getenv("ALGO_RPC"); u != "" {
		return strings.TrimRight(u, "/")
	}
	return "https://mainnet-api.algonode.cloud"
}

func (c *LedgerClient) horizonRPC() string {
	if u := os.Getenv("STELLAR_HORIZON"); u != "" {
		return strings.TrimRight(u, "/")
	}
	return "https://horizon.stellar.org"
}

func (c *LedgerClient) tzktRPC() string {
	if u := os.Getenv("TZKT_URL"); u != "" {
		return strings.TrimRight(u, "/")
	}
	return "https://api.tzkt.io/v1"
}

func (c *LedgerClient) chromiaNode() string {
	if u := os.Getenv("CHROMIA_NODE"); u != "" {
		return strings.TrimRight(u, "/")
	}
	return "https://bootstrap1.chromia.com:7740"
}

// AlgoASABalanceAtRound returns ASA balance (smallest units) at an algod round.
func (c *LedgerClient) AlgoASABalanceAtRound(ctx context.Context, addr string, assetID, round int64) (int64, string, error) {
	cacheChain := fmt.Sprintf("algo:asa:%d", assetID)
	if raw, ok := c.cacheGet(cacheChain, "balance", addr, round); ok {
		v, _ := strconv.ParseInt(string(raw), 10, 64)
		return v, "cache", nil
	}
	path := fmt.Sprintf("%s/v2/accounts/%s?round=%d", c.algodRPC(), url.PathEscape(addr), round)
	body, used, err := c.httpGet(ctx, path)
	if err != nil {
		return 0, used, err
	}
	var out struct {
		Assets []struct {
			AssetID int64  `json:"asset-id"`
			Amount  uint64 `json:"amount"`
		} `json:"assets"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return 0, used, err
	}
	for _, a := range out.Assets {
		if a.AssetID == assetID {
			v := int64(a.Amount)
			c.cachePut(cacheChain, "balance", addr, round, []byte(strconv.FormatInt(v, 10)))
			return v, used, nil
		}
	}
	c.cachePut(cacheChain, "balance", addr, round, []byte("0"))
	return 0, used, nil
}

// AptosCoinBalanceByTypeAtVersion returns CoinStore balance for a Move coin type at ledger version.
func (c *LedgerClient) AptosCoinBalanceByTypeAtVersion(ctx context.Context, addr, coinType string, version int64) (int64, string, error) {
	cacheChain := "aptos:coin:" + coinType
	if raw, ok := c.cacheGet(cacheChain, "balance", addr, version); ok {
		v, _ := strconv.ParseInt(string(raw), 10, 64)
		return v, "cache", nil
	}
	base := strings.TrimRight(c.aptosRPC(), "/")
	resourceType := fmt.Sprintf("0x1::coin::CoinStore<%s>", coinType)
	path := fmt.Sprintf("%s/accounts/%s/resource/%s?ledger_version=%d",
		base, url.PathEscape(addr), url.PathEscape(resourceType), version)
	body, used, err := c.httpGet(ctx, path)
	if err != nil {
		return 0, used, err
	}
	var out struct {
		Data struct {
			Coin struct {
				Value string `json:"value"`
			} `json:"coin"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return 0, used, err
	}
	v, err := strconv.ParseInt(out.Data.Coin.Value, 10, 64)
	if err != nil {
		return 0, used, fmt.Errorf("aptos coin value: %w", err)
	}
	c.cachePut(cacheChain, "balance", addr, version, []byte(strconv.FormatInt(v, 10)))
	return v, used, nil
}

type substrateAssetAccount struct {
	Balance types.U128
	Status  types.U8
}

// SubstrateAssetBalanceAtBlock returns Assets pallet balance at block height.
func (c *LedgerClient) SubstrateAssetBalanceAtBlock(ctx context.Context, rpcURL, ss58Addr string, assetID, height int64) (int64, string, error) {
	cacheChain := fmt.Sprintf("substrate:asset:%d:%s", assetID, rpcURL)
	if raw, ok := c.cacheGet(cacheChain, "balance", ss58Addr, height); ok {
		v, _ := strconv.ParseInt(string(raw), 10, 64)
		return v, "cache", nil
	}
	wsURL := strings.Replace(rpcURL, "https://", "wss://", 1)
	api, err := gsrpc.NewSubstrateAPI(wsURL)
	if err != nil {
		return 0, wsURL, fmt.Errorf("substrate connect: %w", err)
	}
	hash, err := api.RPC.Chain.GetBlockHash(uint64(height))
	if err != nil {
		return 0, wsURL, err
	}
	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return 0, wsURL, err
	}
	_, pubKey, err := subkey.SS58Decode(ss58Addr)
	if err != nil {
		return 0, wsURL, fmt.Errorf("ss58 decode: %w", err)
	}
	var accountID types.AccountID
	if len(pubKey) != 32 {
		return 0, wsURL, fmt.Errorf("unexpected pubkey length %d", len(pubKey))
	}
	copy(accountID[:], pubKey)
	aidBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(aidBytes, uint32(assetID))
	key, err := types.CreateStorageKey(meta, "Assets", "Account", aidBytes, accountID[:])
	if err != nil {
		return 0, wsURL, err
	}
	var acct substrateAssetAccount
	ok, err := api.RPC.State.GetStorage(key, &acct, hash)
	if err != nil {
		return 0, wsURL, err
	}
	if !ok {
		c.cachePut(cacheChain, "balance", ss58Addr, height, []byte("0"))
		return 0, wsURL, nil
	}
	balStr := acct.Balance.String()
	v, ok := new(big.Int).SetString(balStr, 10)
	if !ok {
		return 0, wsURL, fmt.Errorf("substrate asset balance decode")
	}
	if !v.IsInt64() {
		return 0, wsURL, fmt.Errorf("substrate asset balance overflow")
	}
	c.cachePut(cacheChain, "balance", ss58Addr, height, []byte(v.String()))
	return v.Int64(), wsURL, nil
}

// StellarAssetBalanceAtLedger returns asset balance (stroops) at horizon ledger sequence.
func (c *LedgerClient) StellarAssetBalanceAtLedger(ctx context.Context, account, assetCode, issuer string, ledger int64) (int64, string, error) {
	cacheChain := "stellar:" + assetCode + ":" + issuer
	if raw, ok := c.cacheGet(cacheChain, "balance", account, ledger); ok {
		v, _ := strconv.ParseInt(string(raw), 10, 64)
		return v, "cache", nil
	}
	path := fmt.Sprintf("%s/accounts/%s?at=%d", c.horizonRPC(), url.PathEscape(account), ledger)
	body, used, err := c.httpGet(ctx, path)
	if err != nil {
		return 0, used, err
	}
	var out struct {
		Balances []struct {
			Balance   string `json:"balance"`
			AssetCode string `json:"asset_code"`
			AssetType string `json:"asset_type"`
			Issuer    string `json:"asset_issuer"`
		} `json:"balances"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return 0, used, err
	}
	for _, b := range out.Balances {
		if b.AssetCode == assetCode && b.Issuer == issuer {
			v, err := parseDecimalToSmallest(b.Balance, 7)
			if err != nil {
				return 0, used, err
			}
			c.cachePut(cacheChain, "balance", account, ledger, []byte(strconv.FormatInt(v, 10)))
			return v, used, nil
		}
	}
	c.cachePut(cacheChain, "balance", account, ledger, []byte("0"))
	return 0, used, nil
}

// TezosFABalanceAtLevel returns FA token balance (smallest units) at tzkt level.
func (c *LedgerClient) TezosFABalanceAtLevel(ctx context.Context, account, contract string, level int64) (int64, string, error) {
	cacheChain := "tezos:fa:" + contract
	if raw, ok := c.cacheGet(cacheChain, "balance", account, level); ok {
		v, _ := strconv.ParseInt(string(raw), 10, 64)
		return v, "cache", nil
	}
	path := fmt.Sprintf("%s/tokens/balances?account=%s&token.contract=%s&level=%d",
		c.tzktRPC(), url.QueryEscape(account), url.QueryEscape(contract), level)
	body, used, err := c.httpGet(ctx, path)
	if err != nil {
		return 0, used, err
	}
	var rows []struct {
		Balance string `json:"balance"`
	}
	if err := json.Unmarshal(body, &rows); err != nil {
		return 0, used, err
	}
	if len(rows) == 0 {
		c.cachePut(cacheChain, "balance", account, level, []byte("0"))
		return 0, used, nil
	}
	v, err := strconv.ParseInt(rows[0].Balance, 10, 64)
	if err != nil {
		return 0, used, fmt.Errorf("tezos fa balance: %w", err)
	}
	c.cachePut(cacheChain, "balance", account, level, []byte(strconv.FormatInt(v, 10)))
	return v, used, nil
}

// XRPLEmittedBalanceAtLedger returns issued-token balance in smallest units via account_lines.
func (c *LedgerClient) XRPLEmittedBalanceAtLedger(ctx context.Context, account, issuer, currency string, ledgerIndex int64, decimals int) (int64, string, error) {
	cacheChain := "xrpl:issued:" + issuer + ":" + currency
	if raw, ok := c.cacheGet(cacheChain, "balance", account, ledgerIndex); ok {
		v, _ := strconv.ParseInt(string(raw), 10, 64)
		return v, "cache", nil
	}
	bodyRaw, used, err := c.httpPostJSON(ctx, c.xrplRPC, map[string]any{
		"method": "account_lines",
		"params": []any{map[string]any{
			"account":      account,
			"ledger_index": ledgerIndex,
		}},
	})
	if err != nil {
		return 0, used, err
	}
	var out struct {
		Result struct {
			Lines  []struct {
				Account  string `json:"account"`
				Balance  string `json:"balance"`
				Currency string `json:"currency"`
			} `json:"lines"`
			Status string `json:"status"`
		} `json:"result"`
	}
	if err := json.Unmarshal(bodyRaw, &out); err != nil {
		return 0, used, err
	}
	if out.Result.Status != "success" {
		return 0, used, fmt.Errorf("xrpl account_lines failed for %s at %d", account, ledgerIndex)
	}
	for _, line := range out.Result.Lines {
		if strings.EqualFold(line.Account, issuer) && strings.EqualFold(line.Currency, currency) {
			v, err := parseDecimalToSmallest(line.Balance, decimals)
			if err != nil {
				return 0, used, err
			}
			c.cachePut(cacheChain, "balance", account, ledgerIndex, []byte(strconv.FormatInt(v, 10)))
			return v, used, nil
		}
	}
	c.cachePut(cacheChain, "balance", account, ledgerIndex, []byte("0"))
	return 0, used, nil
}

// ChromiaCHRBalanceAtHeight queries FT4 balances on Chromia Economy chain at block height.
func (c *LedgerClient) ChromiaCHRBalanceAtHeight(ctx context.Context, nodeURL, blockchainRID, accountHex string, height int64) (int64, string, error) {
	cacheChain := "chromia:" + blockchainRID
	if raw, ok := c.cacheGet(cacheChain, "balance", accountHex, height); ok {
		v, _ := strconv.ParseInt(string(raw), 10, 64)
		return v, "cache", nil
	}
	node := nodeURL
	if node == "" {
		node = c.chromiaNode()
	}
	accountHex = strings.TrimPrefix(strings.ToLower(accountHex), "0x")
	acctBytes, err := hex.DecodeString(accountHex)
	if err != nil {
		return 0, node, fmt.Errorf("chromia account hex: %w", err)
	}
	queryPath := fmt.Sprintf("%s/query/%s", node, blockchainRID)
	if height > 0 {
		queryPath = fmt.Sprintf("%s/%d", queryPath, height)
	}
	payload := []any{
		"ft4.get_balances_by_account_id",
		map[string]any{"account_id": hex.EncodeToString(acctBytes)},
	}
	bodyRaw, used, err := c.httpPostJSON(ctx, queryPath, payload)
	if err != nil {
		return 0, used, err
	}
	var rows []struct {
		Asset struct {
			Name string `json:"name"`
		} `json:"asset"`
		Amount json.Number `json:"amount"`
	}
	if err := json.Unmarshal(bodyRaw, &rows); err != nil {
		// Some nodes wrap the result.
		var wrapped struct {
			Result json.RawMessage `json:"result"`
		}
		if err2 := json.Unmarshal(bodyRaw, &wrapped); err2 == nil && len(wrapped.Result) > 0 {
			if err3 := json.Unmarshal(wrapped.Result, &rows); err3 != nil {
				return 0, used, fmt.Errorf("chromia decode: %w", err)
			}
		} else {
			return 0, used, fmt.Errorf("chromia decode: %w", err)
		}
	}
	for _, row := range rows {
		if strings.EqualFold(row.Asset.Name, "CHR") {
			v, err := row.Amount.Int64()
			if err != nil {
				bi := new(big.Int)
				if _, ok := bi.SetString(row.Amount.String(), 10); !ok {
					return 0, used, fmt.Errorf("chromia amount decode")
				}
				if !bi.IsInt64() {
					return 0, used, fmt.Errorf("chromia amount overflow")
				}
				v = bi.Int64()
			}
			c.cachePut(cacheChain, "balance", accountHex, height, []byte(strconv.FormatInt(v, 10)))
			return v, used, nil
		}
	}
	c.cachePut(cacheChain, "balance", accountHex, height, []byte("0"))
	return 0, used, nil
}

func parseDecimalToSmallest(s string, decimals int) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	neg := false
	if strings.HasPrefix(s, "-") {
		neg = true
		s = s[1:]
	}
	parts := strings.SplitN(s, ".", 2)
	whole, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, err
	}
	fracStr := ""
	if len(parts) == 2 {
		fracStr = parts[1]
	}
	if len(fracStr) > decimals {
		fracStr = fracStr[:decimals]
	}
	for len(fracStr) < decimals {
		fracStr += "0"
	}
	var frac int64
	if fracStr != "" {
		frac, err = strconv.ParseInt(fracStr, 10, 64)
		if err != nil {
			return 0, err
		}
	}
	multiplier := int64(1)
	for i := 0; i < decimals; i++ {
		multiplier *= 10
	}
	out := whole*multiplier + frac
	if neg {
		out = -out
	}
	return out, nil
}
