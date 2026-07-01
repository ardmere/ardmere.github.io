package porprobe

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/ardmere/ardmere/internal/rpc"
)

type chainlistChain struct {
	ChainID int             `json:"chainId"`
	Name    string          `json:"name"`
	RPC     json.RawMessage `json:"rpc"`
}

func RunRPC(args []string) {
	fs := flag.NewFlagSet("rpc", flag.ExitOnError)
	netFlag := fs.String("network", "BSC", "ETH or BSC")
	height := fs.Int64("height", 101590091, "historical block to test")
	address := fs.String("address", "0x86523c87c8ec98c7539e2c58cd813ee9d1a08d96", "address for eth_getBalance probe")
	chainlist := fs.Bool("chainlist", false, "also probe extra Chainlist RPC URLs")
	depth := fs.Bool("depth", false, "measure archive depth for working providers")
	_ = fs.Parse(args)

	net := rpc.Network(strings.ToUpper(*netFlag))
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	providers, err := rpc.LoadProviderConfig()
	if err != nil {
		log.Fatalf("load providers: %v", err)
	}
	configured := providers[net]
	if len(configured) == 0 {
		log.Fatalf("no providers for %s", net)
	}

	var candidates []rpc.Provider
	seen := map[string]bool{}
	for _, p := range configured {
		if seen[p.URL] {
			continue
		}
		seen[p.URL] = true
		candidates = append(candidates, p)
	}
	if *chainlist {
		extra, err := fetchChainlistRPCs(ctx, net)
		if err != nil {
			log.Printf("chainlist fetch: %v", err)
		} else {
			for _, url := range extra {
				if seen[url] {
					continue
				}
				seen[url] = true
				candidates = append(candidates, rpc.Provider{URL: url, Archive: false, Weight: 1})
			}
		}
	}

	fmt.Printf("probing %d %s providers at block %d\n", len(candidates), net, *height)
	blk := fmt.Sprintf("0x%x", *height)
	type row struct {
		URL     string
		Archive bool
		OK      bool
		Err     string
		Depth   int64
	}
	var rows []row
	for _, p := range candidates {
		_, err := singleRPC(ctx, p.URL, "eth_getBalance", []any{*address, blk})
		r := row{URL: p.URL, Archive: p.Archive}
		if err != nil {
			r.Err = err.Error()
		} else {
			r.OK = true
		}
		if r.OK && *depth {
			r.Depth = measureArchiveDepth(ctx, p.URL, *address)
		}
		rows = append(rows, r)
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].OK != rows[j].OK {
			return rows[i].OK
		}
		return rows[i].Depth > rows[j].Depth
	})
	for _, r := range rows {
		status := "FAIL"
		if r.OK {
			status = "OK"
		}
		depthNote := ""
		if r.Depth > 0 {
			depthNote = fmt.Sprintf(" depth~%d blocks", r.Depth)
		}
		archiveNote := ""
		if r.Archive {
			archiveNote = " [config:archive]"
		}
		fmt.Printf("%-8s %-48s%s%s", status, r.URL, archiveNote, depthNote)
		if r.Err != "" {
			msg := r.Err
			if len(msg) > 100 {
				msg = msg[:100] + "..."
			}
			fmt.Printf(" -> %s", msg)
		}
		fmt.Println()
	}
}

func chainIDForNetwork(net rpc.Network) (int, bool) {
	switch net {
	case rpc.NetEthereum:
		return 1, true
	case rpc.NetBSC:
		return 56, true
	case rpc.NetArbitrum:
		return 42161, true
	case rpc.NetOptimism:
		return 10, true
	case rpc.NetBase:
		return 8453, true
	case rpc.NetPolygon:
		return 137, true
	case rpc.NetAvalanche:
		return 43114, true
	case rpc.NetOpBNB:
		return 204, true
	case rpc.NetSonic:
		return 146, true
	case rpc.NetWorld:
		return 480, true
	default:
		return 0, false
	}
}

func fetchChainlistRPCs(ctx context.Context, net rpc.Network) ([]string, error) {
	chainID, ok := chainIDForNetwork(net)
	if !ok {
		return nil, fmt.Errorf("unsupported network %s", net)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://chainlist.org/rpcs.json", nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var chains []chainlistChain
	if err := json.Unmarshal(raw, &chains); err != nil {
		return nil, err
	}
	var out []string
	for _, chain := range chains {
		if chain.ChainID != chainID {
			continue
		}
		rpcs, err := decodeChainlistRPCs(chain.RPC)
		if err != nil {
			continue
		}
		for _, url := range rpcs {
			if url == "" || strings.Contains(url, "${") {
				continue
			}
			if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
				continue
			}
			out = append(out, url)
		}
	}
	return out, nil
}

func decodeChainlistRPCs(raw json.RawMessage) ([]string, error) {
	var stringsOnly []string
	if err := json.Unmarshal(raw, &stringsOnly); err == nil {
		for i := range stringsOnly {
			stringsOnly[i] = strings.TrimSpace(stringsOnly[i])
		}
		return stringsOnly, nil
	}
	var mixed []any
	if err := json.Unmarshal(raw, &mixed); err != nil {
		return nil, err
	}
	out := make([]string, 0, len(mixed))
	for _, item := range mixed {
		switch v := item.(type) {
		case string:
			out = append(out, strings.TrimSpace(v))
		case map[string]any:
			if url, ok := v["url"].(string); ok {
				out = append(out, strings.TrimSpace(url))
			}
		}
	}
	return out, nil
}

func measureArchiveDepth(ctx context.Context, url, address string) int64 {
	tip, err := singleRPC(ctx, url, "eth_blockNumber", []any{})
	if err != nil {
		return 0
	}
	var tipHex string
	if err := json.Unmarshal(tip, &tipHex); err != nil {
		return 0
	}
	tipInt := parseHexInt(tipHex)
	if tipInt <= 0 {
		return 0
	}
	offsets := []int64{0, 100, 500, 1000, 5000, 10000, 50000, 100000, 500000, 1000000}
	lastOK := int64(0)
	for _, off := range offsets {
		h := tipInt - off
		if h < 0 {
			break
		}
		blk := fmt.Sprintf("0x%x", h)
		_, err := singleRPC(ctx, url, "eth_getBalance", []any{address, blk})
		if err == nil {
			lastOK = off
			continue
		}
		break
	}
	if lastOK == 0 {
		return 0
	}
	lo, hi := lastOK, minInt64(lastOK*2, tipInt)
	for i := 0; i < 12; i++ {
		mid := (lo + hi) / 2
		h := tipInt - mid
		blk := fmt.Sprintf("0x%x", h)
		_, err := singleRPC(ctx, url, "eth_getBalance", []any{address, blk})
		if err == nil {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return lo
}

func singleRPC(ctx context.Context, url, method string, params []any) (json.RawMessage, error) {
	body, _ := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      1,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "ardmere-por-probe/0.1")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	var rr struct {
		Result json.RawMessage `json:"result"`
		Error  *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(raw, &rr); err != nil {
		return nil, err
	}
	if rr.Error != nil {
		return nil, fmt.Errorf("%s", rr.Error.Message)
	}
	return rr.Result, nil
}

func parseHexInt(hexStr string) int64 {
	hexStr = strings.TrimPrefix(hexStr, "0x")
	var n int64
	for _, ch := range hexStr {
		n <<= 4
		switch {
		case ch >= '0' && ch <= '9':
			n += int64(ch - '0')
		case ch >= 'a' && ch <= 'f':
			n += int64(ch-'a') + 10
		case ch >= 'A' && ch <= 'F':
			n += int64(ch-'A') + 10
		}
	}
	return n
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
