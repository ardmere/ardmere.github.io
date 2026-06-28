package verifier

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// ETHDepositIndexer maps an execution-layer address to ETH sent into the
// canonical Eth2 DepositContract before the snapshot block.
type ETHDepositIndexer interface {
	DepositedETH(ctx context.Context, address string, endBlock int64) (ETHDepositAccounted, error)
}

type ETHDepositAccounted struct {
	Deposited decimal.Decimal
	TxCount   int
	Source    string
	Note      string
}

type EtherscanETHDepositIndexer struct {
	APIKey  string
	BaseURL string
	httpc   *http.Client
}

func NewEtherscanETHDepositIndexer(apiKey string) *EtherscanETHDepositIndexer {
	return &EtherscanETHDepositIndexer{
		APIKey:  apiKey,
		BaseURL: "https://api.etherscan.io/v2/api",
		httpc:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (e *EtherscanETHDepositIndexer) DepositedETH(ctx context.Context, address string, endBlock int64) (ETHDepositAccounted, error) {
	if e == nil || strings.TrimSpace(e.APIKey) == "" {
		return ETHDepositAccounted{}, fmt.Errorf("ETHERSCAN_API_KEY not configured")
	}
	base := e.BaseURL
	if base == "" {
		base = "https://api.etherscan.io/v2/api"
	}
	httpc := e.httpc
	if httpc == nil {
		httpc = http.DefaultClient
	}

	var wei big.Int
	txCount := 0
	const pageSize = 10000
	for page := 1; ; page++ {
		txs, err := e.etherscanTxPage(ctx, httpc, base, address, endBlock, page, pageSize)
		if err != nil {
			return ETHDepositAccounted{}, err
		}
		for _, tx := range txs {
			if tx.IsError != "0" {
				continue
			}
			if !strings.EqualFold(tx.To, ethDepositContract) {
				continue
			}
			value, ok := new(big.Int).SetString(tx.Value, 10)
			if !ok {
				return ETHDepositAccounted{}, fmt.Errorf("bad tx value for %s: %q", tx.Hash, tx.Value)
			}
			wei.Add(&wei, value)
			txCount++
		}
		if len(txs) < pageSize {
			break
		}
	}

	return ETHDepositAccounted{
		Deposited: decimal.NewFromBigInt(&wei, -evmNativeDecimals),
		TxCount:   txCount,
		Source:    "etherscan:account.txlist",
		Note:      fmt.Sprintf("summed tx.value sent to %s up to block %d", ethDepositContract, endBlock),
	}, nil
}

type etherscanTx struct {
	BlockNumber string `json:"blockNumber"`
	Hash        string `json:"hash"`
	To          string `json:"to"`
	Value       string `json:"value"`
	IsError     string `json:"isError"`
}

func (e *EtherscanETHDepositIndexer) etherscanTxPage(ctx context.Context, httpc *http.Client, base, address string, endBlock int64, page, pageSize int) ([]etherscanTx, error) {
	q := url.Values{}
	q.Set("chainid", "1")
	q.Set("module", "account")
	q.Set("action", "txlist")
	q.Set("address", address)
	q.Set("startblock", "0")
	q.Set("endblock", fmt.Sprintf("%d", endBlock))
	q.Set("page", fmt.Sprintf("%d", page))
	q.Set("offset", fmt.Sprintf("%d", pageSize))
	q.Set("sort", "asc")
	q.Set("apikey", e.APIKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+"?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "ardmere/0.1")
	resp, err := httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("etherscan txlist HTTP %d", resp.StatusCode)
	}

	var body struct {
		Status  string          `json:"status"`
		Message string          `json:"message"`
		Result  json.RawMessage `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	if body.Status == "0" && strings.EqualFold(body.Message, "No transactions found") {
		return nil, nil
	}
	if body.Status != "1" {
		return nil, fmt.Errorf("etherscan txlist: status=%s message=%s", body.Status, body.Message)
	}
	var txs []etherscanTx
	if err := json.Unmarshal(body.Result, &txs); err != nil {
		return nil, fmt.Errorf("decode etherscan txlist result: %w", err)
	}
	return txs, nil
}
