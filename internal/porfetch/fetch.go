package porfetch

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ardmere/ardmere/internal/exchange"
	"github.com/ardmere/ardmere/internal/exchangereg"
	"github.com/ardmere/ardmere/internal/exchanges/bingx/bingxapi"
	"github.com/ardmere/ardmere/internal/exchanges/bitget/bitgetapi"
	"github.com/ardmere/ardmere/internal/exchanges/bitunix/bitunixapi"
	"github.com/ardmere/ardmere/internal/exchanges/btcc/btccapi"
	"github.com/ardmere/ardmere/internal/exchanges/bybit/bybitapi"
	"github.com/ardmere/ardmere/internal/exchanges/coinex/coinexapi"
	"github.com/ardmere/ardmere/internal/exchanges/gateio/gateapi"
	"github.com/ardmere/ardmere/internal/exchanges/kucoin/kucoinapi"
	"github.com/ardmere/ardmere/internal/exchanges/mexc/mexcapi"
	"github.com/ardmere/ardmere/internal/exchanges/ourbit/ourbitapi"
	"github.com/ardmere/ardmere/internal/exchanges/phemex/phemexapi"
	"github.com/ardmere/ardmere/internal/exchanges/poloniex/poloniexapi"
	"github.com/ardmere/ardmere/internal/exchanges/toobit/toobitapi"
	"github.com/ardmere/ardmere/internal/exchanges/woox/wooxapi"
)

// Main runs fetch subcommands; args[0] is gateio or okx.
func Main(args []string) int {
	if len(args) < 1 {
		usage()
		return 2
	}
	switch args[0] {
	case "gateio":
		RunGateio(args[1:])
		return 0
	case "htx":
		RunHTX(args[1:])
		return 0
	case "bybit":
		RunBybit(args[1:])
		return 0
	case "bitget":
		RunBitget(args[1:])
		return 0
	case "mexc":
		RunMEXC(args[1:])
		return 0
	case "kucoin":
		RunKuCoin(args[1:])
		return 0
	case "bingx":
		RunBingX(args[1:])
		return 0
	case "poloniex":
		RunPoloniex(args[1:])
		return 0
	case "ourbit":
		RunOurbit(args[1:])
		return 0
	case "coinex":
		RunCoinEx(args[1:])
		return 0
	case "bitunix":
		RunBitunix(args[1:])
		return 0
	case "woox":
		RunWOOX(args[1:])
		return 0
	case "toobit":
		RunToobit(args[1:])
		return 0
	case "phemex":
		RunPhemex(args[1:])
		return 0
	case "btcc":
		RunBTCC(args[1:])
		return 0
	case "okx":
		RunOKX(args[1:])
		return 0
	case "help", "-h", "--help":
		usage()
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown exchange %q\n\n", args[0])
		usage()
		return 2
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `usage: por fetch <exchange> [flags]

Exchanges:
  gateio     fetch or import Gate.com public PoR summary
  htx        fetch or import HTX zk PoR bundle (public-data.zip)
  bybit      fetch or import Bybit PoR summary (+ optional user proof)
  bitget     fetch or import Bitget PoR summary (+ optional merkel_tree_bg.json)
  mexc       fetch or import MEXC PoR summary (+ optional myProof.json)
  kucoin     fetch or import KuCoin PoR summary (public SSR scrape)
  bingx      fetch or import BingX PoR summary (+ optional user Merkle JSON)
  poloniex   fetch or import Poloniex PoR summary (public proxy/sapi)
  ourbit     fetch or import Ourbit PoR summary (public stock_info + wallets)
  coinex     fetch or import CoinEx PoR summary (public merkle/reserve + optional user proof)
  bitunix    fetch or import Bitunix PoR summary (+ optional merkle_tree_bitunix.json)
  woox       fetch or import WOO X PoR summary (live dashboard + merkle audit)
  toobit     fetch or import Toobit PoR summary (+ optional SMT proof JSON)
  phemex     fetch or import Phemex PoR summary (porVersions balances)
  btcc       fetch or import BTCC PoR summary (reserveInfo ratios + truncated root)
  okx        fetch or import OKX public PoR artifacts

Examples:
  por fetch gateio
  por fetch gateio -info-file ./info.json -coins-file ./coinList.json
  por fetch htx -zk-bundle ./public-data.zip
  por fetch bybit -summary-path ./summary.json -user-proof ./myProof.json
  por fetch bitget -summary-path ./summary.json -user-proof ./merkel_tree_bg.json
  por fetch mexc -summary-path ./fixtures/mexc/202607-summary.json -user-proof ./fixtures/mexc/merkle_proof.json
  por fetch kucoin
  por fetch bingx -summary-path ./fixtures/bingx/20261015-summary.json
  por fetch poloniex
  por fetch poloniex -audit-id 20260901
  por fetch ourbit
  por fetch ourbit -summary-path ./fixtures/ourbit/20261101-summary.json
  por fetch coinex
  por fetch coinex -summary-path ./fixtures/coinex/20260914-summary.json
  por fetch bitunix
  por fetch bitunix -summary-path ./fixtures/bitunix/20261103-summary.json
  por fetch woox
  por fetch woox -summary-path ./fixtures/woox/20261015-summary.json
  por fetch toobit
  por fetch toobit -summary-path ./fixtures/toobit/20261001-summary.json
  por fetch phemex
  por fetch phemex -summary-path ./fixtures/phemex/20261101-summary.json
  por fetch btcc
  por fetch btcc -summary-path ./fixtures/btcc/20261015-summary.json
  por fetch okx -skip-wallet
`)
}

// RunGateio archives Gate public PoR data.
func RunGateio(args []string) {
	fs := flag.NewFlagSet("gateio", flag.ExitOnError)
	artifactsBase := fs.String("artifacts", "./artifacts", "artifacts root")
	summaryPath := fs.String("summary-path", "", "import merged summary bundle JSON")
	infoFile := fs.String("info-file", "", "browser-captured getProofOfReservesInfo response JSON")
	coinsFile := fs.String("coins-file", "", "browser-captured getProofOfReservesCoinList response JSON")
	listFile := fs.String("list-file", "", "optional browser-captured getProofOfReservesList response JSON")
	zkBundlePath := fs.String("zk-bundle", "", "import local zkmerkle_cex tar.gz into raw/")
	_ = fs.Parse(args)

	opts := exchange.FetchOpts{
		SummaryPath:  *summaryPath,
		ZkBundlePath: *zkBundlePath,
	}
	if *infoFile != "" {
		merged, err := gateapi.BuildSummaryBundle(*infoFile, *coinsFile, *listFile)
		if err != nil {
			log.Fatalf("merge browser capture: %v", err)
		}
		tmp, err := os.CreateTemp("", "gate-summary-*.json")
		if err != nil {
			log.Fatalf("temp file: %v", err)
		}
		tmpPath := tmp.Name()
		if _, err := tmp.Write(merged); err != nil {
			log.Fatalf("write temp: %v", err)
		}
		tmp.Close()
		defer os.Remove(tmpPath)
		opts.SummaryPath = tmpPath
		opts.ImportSource = "browser-api"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	adapter, err := exchangereg.Get("gateio")
	if err != nil {
		log.Fatalf("exchange: %v", err)
	}
	stored, err := adapter.FetchAndStore(ctx, *artifactsBase, opts)
	if err != nil {
		log.Fatalf("fetch gate snapshot: %v\n\nTip: open https://www.gate.com/proof-of-reserves in a browser, copy API responses from DevTools, then:\n  por fetch gateio -info-file ./info.json -coins-file ./coinList.json", err)
	}
	logResult(stored)
}

// RunHTX archives HTX zk PoR artifacts (public-data.zip).
func RunHTX(args []string) {
	fs := flag.NewFlagSet("htx", flag.ExitOnError)
	artifactsBase := fs.String("artifacts", "./artifacts", "artifacts root")
	summaryPath := fs.String("summary-path", "", "import browser-captured reserve ratio JSON")
	zkBundlePath := fs.String("zk-bundle", "", "import local public-data.zip (GitHub release or PoR Reports)")
	_ = fs.Parse(args)

	if *zkBundlePath == "" && *summaryPath == "" {
		log.Fatalf("htx fetch requires -zk-bundle or -summary-path\n\nTip: download sample from https://github.com/huobiapi/Tool-Go-MerkleVerify/releases/download/2.0.0/public-data.zip")
	}

	opts := exchange.FetchOpts{
		SummaryPath:  *summaryPath,
		ZkBundlePath: *zkBundlePath,
	}
	if *summaryPath != "" {
		opts.ImportSource = "import"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	adapter, err := exchangereg.Get("htx")
	if err != nil {
		log.Fatalf("exchange: %v", err)
	}
	stored, err := adapter.FetchAndStore(ctx, *artifactsBase, opts)
	if err != nil {
		log.Fatalf("fetch htx snapshot: %v", err)
	}
	logResult(stored)
}

// RunBybit archives Bybit PoR summary (and optional user Merkle proof).
func RunBybit(args []string) {
	fs := flag.NewFlagSet("bybit", flag.ExitOnError)
	artifactsBase := fs.String("artifacts", "./artifacts", "artifacts root")
	summaryPath := fs.String("summary-path", "", "import merged summary bundle JSON")
	infoFile := fs.String("info-file", "", "browser-captured reserve-ratio API response JSON")
	coinsFile := fs.String("coins-file", "", "browser-captured per-coin reserve list JSON")
	userProofPath := fs.String("user-proof", "", "import login-gated myProof.json")
	_ = fs.Parse(args)

	opts := exchange.FetchOpts{
		SummaryPath:   *summaryPath,
		UserProofPath: *userProofPath,
	}
	if *infoFile != "" {
		merged, err := bybitapi.BuildSummaryBundle(*infoFile, *coinsFile)
		if err != nil {
			log.Fatalf("merge browser capture: %v", err)
		}
		tmp, err := os.CreateTemp("", "bybit-summary-*.json")
		if err != nil {
			log.Fatalf("temp file: %v", err)
		}
		tmpPath := tmp.Name()
		if _, err := tmp.Write(merged); err != nil {
			log.Fatalf("write temp: %v", err)
		}
		tmp.Close()
		defer os.Remove(tmpPath)
		opts.SummaryPath = tmpPath
		opts.ImportSource = "browser-api"
	} else if *summaryPath != "" {
		opts.ImportSource = "import"
	}

	if opts.SummaryPath == "" {
		opts.ImportSource = ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	adapter, err := exchangereg.Get("bybit")
	if err != nil {
		log.Fatalf("exchange: %v", err)
	}
	stored, err := adapter.FetchAndStore(ctx, *artifactsBase, opts)
	if err != nil {
		log.Fatalf("fetch bybit snapshot: %v\n\nTip: open https://www.bybit.com/en/proof-of-reserves in a browser, copy API JSON from DevTools, then:\n  por fetch bybit -info-file ./info.json -coins-file ./coins.json", err)
	}
	logResult(stored)
}

// RunBitget archives Bitget PoR summary (and optional user Merkle proof).
func RunBitget(args []string) {
	fs := flag.NewFlagSet("bitget", flag.ExitOnError)
	artifactsBase := fs.String("artifacts", "./artifacts", "artifacts root")
	summaryPath := fs.String("summary-path", "", "import merged summary bundle JSON")
	infoFile := fs.String("info-file", "", "browser-captured reserve-ratio API response JSON")
	coinsFile := fs.String("coins-file", "", "browser-captured per-coin reserve list JSON")
	userProofPath := fs.String("user-proof", "", "import login-gated merkel_tree_bg.json")
	_ = fs.Parse(args)

	opts := exchange.FetchOpts{
		SummaryPath:   *summaryPath,
		UserProofPath: *userProofPath,
	}
	if *infoFile != "" {
		merged, err := bitgetapi.BuildSummaryBundle(*infoFile, *coinsFile)
		if err != nil {
			log.Fatalf("merge browser capture: %v", err)
		}
		tmp, err := os.CreateTemp("", "bitget-summary-*.json")
		if err != nil {
			log.Fatalf("temp file: %v", err)
		}
		tmpPath := tmp.Name()
		if _, err := tmp.Write(merged); err != nil {
			log.Fatalf("write temp: %v", err)
		}
		tmp.Close()
		defer os.Remove(tmpPath)
		opts.SummaryPath = tmpPath
		opts.ImportSource = "browser-api"
	} else if *summaryPath != "" {
		opts.ImportSource = "import"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	adapter, err := exchangereg.Get("bitget")
	if err != nil {
		log.Fatalf("exchange: %v", err)
	}
	stored, err := adapter.FetchAndStore(ctx, *artifactsBase, opts)
	if err != nil {
		log.Fatalf("fetch bitget snapshot: %v\n\nTip: open https://www.bitget.com/proof-of-reserves in a browser, copy API JSON from DevTools, then:\n  por fetch bitget -info-file ./info.json -coins-file ./coins.json", err)
	}
	logResult(stored)
}

// RunMEXC archives MEXC PoR summary (and optional user Merkle proof).
func RunMEXC(args []string) {
	fs := flag.NewFlagSet("mexc", flag.ExitOnError)
	artifactsBase := fs.String("artifacts", "./artifacts", "artifacts root")
	summaryPath := fs.String("summary-path", "", "import merged summary bundle JSON")
	infoFile := fs.String("info-file", "", "browser-captured reserve-ratio API response JSON")
	coinsFile := fs.String("coins-file", "", "browser-captured per-coin reserve list JSON")
	userProofPath := fs.String("user-proof", "", "import login-gated myProof.json")
	_ = fs.Parse(args)

	opts := exchange.FetchOpts{
		SummaryPath:   *summaryPath,
		UserProofPath: *userProofPath,
	}
	if *infoFile != "" {
		merged, err := mexcapi.BuildSummaryBundle(*infoFile, *coinsFile)
		if err != nil {
			log.Fatalf("merge browser capture: %v", err)
		}
		tmp, err := os.CreateTemp("", "mexc-summary-*.json")
		if err != nil {
			log.Fatalf("temp file: %v", err)
		}
		tmpPath := tmp.Name()
		if _, err := tmp.Write(merged); err != nil {
			log.Fatalf("write temp: %v", err)
		}
		tmp.Close()
		defer os.Remove(tmpPath)
		opts.SummaryPath = tmpPath
		opts.ImportSource = "browser-api"
	} else if *summaryPath != "" {
		opts.ImportSource = "import"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	adapter, err := exchangereg.Get("mexc")
	if err != nil {
		log.Fatalf("exchange: %v", err)
	}
	stored, err := adapter.FetchAndStore(ctx, *artifactsBase, opts)
	if err != nil {
		log.Fatalf("fetch mexc snapshot: %v\n\nTip: MEXC PoR JSON is WAF/login gated. Import a press-release or browser capture:\n  por fetch mexc -summary-path ./fixtures/mexc/202607-summary.json\n  por fetch mexc -summary-path ./summary.json -user-proof ./myProof.json", err)
	}
	logResult(stored)
}

// RunKuCoin archives KuCoin PoR summary from the public SSR page (or import).
func RunKuCoin(args []string) {
	fs := flag.NewFlagSet("kucoin", flag.ExitOnError)
	artifactsBase := fs.String("artifacts", "./artifacts", "artifacts root")
	summaryPath := fs.String("summary-path", "", "import merged summary bundle JSON")
	infoFile := fs.String("info-file", "", "browser-captured initialReserve / reserve API JSON")
	coinsFile := fs.String("coins-file", "", "optional coin list JSON (defaults to info-file)")
	userProofPath := fs.String("user-proof", "", "import login-gated user Merkle export (stored only; no OSS verifier)")
	_ = fs.Parse(args)

	opts := exchange.FetchOpts{
		SummaryPath:   *summaryPath,
		UserProofPath: *userProofPath,
	}
	if *infoFile != "" {
		merged, err := kucoinapi.BuildSummaryBundle(*infoFile, *coinsFile)
		if err != nil {
			log.Fatalf("merge browser capture: %v", err)
		}
		tmp, err := os.CreateTemp("", "kucoin-summary-*.json")
		if err != nil {
			log.Fatalf("temp file: %v", err)
		}
		tmpPath := tmp.Name()
		if _, err := tmp.Write(merged); err != nil {
			log.Fatalf("write temp: %v", err)
		}
		tmp.Close()
		defer os.Remove(tmpPath)
		opts.SummaryPath = tmpPath
		opts.ImportSource = "browser-api"
	} else if *summaryPath != "" {
		opts.ImportSource = "import"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	adapter, err := exchangereg.Get("kucoin")
	if err != nil {
		log.Fatalf("exchange: %v", err)
	}
	stored, err := adapter.FetchAndStore(ctx, *artifactsBase, opts)
	if err != nil {
		log.Fatalf("fetch kucoin snapshot: %v\n\nTip: open https://www.kucoin.com/proof-of-reserves or import:\n  por fetch kucoin -summary-path ./fixtures/kucoin/20260930-summary.json", err)
	}
	logResult(stored)
}

// RunBingX archives BingX PoR summary (and optional user Merkle proof).
func RunBingX(args []string) {
	fs := flag.NewFlagSet("bingx", flag.ExitOnError)
	artifactsBase := fs.String("artifacts", "./artifacts", "artifacts root")
	summaryPath := fs.String("summary-path", "", "import merged summary bundle JSON")
	infoFile := fs.String("info-file", "", "browser-captured reserve-ratio API response JSON")
	coinsFile := fs.String("coins-file", "", "browser-captured per-coin reserve list JSON")
	userProofPath := fs.String("user-proof", "", "import login-gated leaf/path JSON")
	_ = fs.Parse(args)

	opts := exchange.FetchOpts{
		SummaryPath:   *summaryPath,
		UserProofPath: *userProofPath,
	}
	if *infoFile != "" {
		merged, err := bingxapi.BuildSummaryBundle(*infoFile, *coinsFile)
		if err != nil {
			log.Fatalf("merge browser capture: %v", err)
		}
		tmp, err := os.CreateTemp("", "bingx-summary-*.json")
		if err != nil {
			log.Fatalf("temp file: %v", err)
		}
		tmpPath := tmp.Name()
		if _, err := tmp.Write(merged); err != nil {
			log.Fatalf("write temp: %v", err)
		}
		tmp.Close()
		defer os.Remove(tmpPath)
		opts.SummaryPath = tmpPath
		opts.ImportSource = "browser-api"
	} else if *summaryPath != "" {
		opts.ImportSource = "import"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	adapter, err := exchangereg.Get("bingx")
	if err != nil {
		log.Fatalf("exchange: %v", err)
	}
	stored, err := adapter.FetchAndStore(ctx, *artifactsBase, opts)
	if err != nil {
		log.Fatalf("fetch bingx snapshot: %v\n\nTip: BingX reserve APIs are app/security gated. Import a support-article fixture:\n  por fetch bingx -summary-path ./fixtures/bingx/20261015-summary.json", err)
	}
	logResult(stored)
}

// RunPoloniex archives Poloniex PoR summary from public proxy/sapi endpoints (or import).
func RunPoloniex(args []string) {
	fs := flag.NewFlagSet("poloniex", flag.ExitOnError)
	artifactsBase := fs.String("artifacts", "./artifacts", "artifacts root")
	summaryPath := fs.String("summary-path", "", "import merged summary bundle JSON")
	infoFile := fs.String("info-file", "", "browser-captured versionInfo / info JSON")
	coinsFile := fs.String("coins-file", "", "browser-captured snapshotBalance JSON")
	auditID := fs.String("audit-id", "", "snapshot version id (YYYYMMDD); default currentVersion")
	userProofPath := fs.String("user-proof", "", "import login-gated merkleProof.json (stored only; no OSS verifier)")
	_ = fs.Parse(args)

	opts := exchange.FetchOpts{
		SummaryPath:   *summaryPath,
		UserProofPath: *userProofPath,
		AuditID:       *auditID,
	}
	if *infoFile != "" {
		merged, err := poloniexapi.BuildSummaryBundle(*infoFile, *coinsFile, *auditID)
		if err != nil {
			log.Fatalf("merge browser capture: %v", err)
		}
		tmp, err := os.CreateTemp("", "poloniex-summary-*.json")
		if err != nil {
			log.Fatalf("temp file: %v", err)
		}
		tmpPath := tmp.Name()
		if _, err := tmp.Write(merged); err != nil {
			log.Fatalf("write temp: %v", err)
		}
		tmp.Close()
		defer os.Remove(tmpPath)
		opts.SummaryPath = tmpPath
		opts.ImportSource = "browser-api"
	} else if *summaryPath != "" {
		opts.ImportSource = "import"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	adapter, err := exchangereg.Get("poloniex")
	if err != nil {
		log.Fatalf("exchange: %v", err)
	}
	stored, err := adapter.FetchAndStore(ctx, *artifactsBase, opts)
	if err != nil {
		log.Fatalf("fetch poloniex snapshot: %v\n\nTip: open https://www.poloniex.com/topic/security/proof-of-reserves or import:\n  por fetch poloniex -summary-path ./fixtures/poloniex/20261001-summary.json", err)
	}
	logResult(stored)
}

// RunOKX archives OKX public PoR artifacts.
func RunOurbit(args []string) {
	fs := flag.NewFlagSet("ourbit", flag.ExitOnError)
	artifactsBase := fs.String("artifacts", "./artifacts", "artifacts root")
	summaryPath := fs.String("summary-path", "", "import merged summary bundle JSON")
	infoFile := fs.String("info-file", "", "browser-captured stock_info JSON")
	walletsFile := fs.String("wallets-file", "", "browser-captured wallet_detail JSON")
	currenciesFile := fs.String("currencies-file", "", "browser-captured currencies JSON")
	userProofPath := fs.String("user-proof", "", "import login-gated merkel_tree_ourbit.json")
	_ = fs.Parse(args)

	opts := exchange.FetchOpts{
		SummaryPath:   *summaryPath,
		UserProofPath: *userProofPath,
	}
	if *infoFile != "" {
		merged, err := ourbitapi.BuildSummaryBundle(*infoFile, *walletsFile, *currenciesFile)
		if err != nil {
			log.Fatalf("merge browser capture: %v", err)
		}
		tmp, err := os.CreateTemp("", "ourbit-summary-*.json")
		if err != nil {
			log.Fatalf("temp file: %v", err)
		}
		tmpPath := tmp.Name()
		if _, err := tmp.Write(merged); err != nil {
			log.Fatalf("write temp: %v", err)
		}
		tmp.Close()
		defer os.Remove(tmpPath)
		opts.SummaryPath = tmpPath
		opts.ImportSource = "browser-api"
	} else if *summaryPath != "" {
		opts.ImportSource = "import"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	adapter, err := exchangereg.Get("ourbit")
	if err != nil {
		log.Fatalf("exchange: %v", err)
	}
	stored, err := adapter.FetchAndStore(ctx, *artifactsBase, opts)
	if err != nil {
		log.Fatalf("fetch ourbit snapshot: %v\n\nTip: open https://www.ourbit.com/assets/reserve-proof or import:\n  por fetch ourbit -summary-path ./fixtures/ourbit/20261101-summary.json", err)
	}
	logResult(stored)
}

// RunOKX archives OKX public PoR artifacts.

func RunCoinEx(args []string) {
	fs := flag.NewFlagSet("coinex", flag.ExitOnError)
	artifactsBase := fs.String("artifacts", "./artifacts", "artifacts root")
	summaryPath := fs.String("summary-path", "", "import archived summary bundle JSON")
	periodFile := fs.String("period-file", "", "browser-captured period object or full merkle/reserve API JSON")
	userProofPath := fs.String("user-proof", "", "import login-gated Copy-my-audition JSON")
	_ = fs.Parse(args)

	opts := exchange.FetchOpts{
		SummaryPath:   *summaryPath,
		UserProofPath: *userProofPath,
	}
	if *periodFile != "" {
		merged, err := coinexapi.BuildSummaryBundle(*periodFile)
		if err != nil {
			log.Fatalf("build period capture: %v", err)
		}
		tmp, err := os.CreateTemp("", "coinex-summary-*.json")
		if err != nil {
			log.Fatalf("temp file: %v", err)
		}
		tmpPath := tmp.Name()
		if _, err := tmp.Write(merged); err != nil {
			log.Fatalf("write temp: %v", err)
		}
		tmp.Close()
		defer os.Remove(tmpPath)
		opts.SummaryPath = tmpPath
		opts.ImportSource = "browser-api"
	} else if *summaryPath != "" {
		opts.ImportSource = "import"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	adapter, err := exchangereg.Get("coinex")
	if err != nil {
		log.Fatalf("exchange: %v", err)
	}
	stored, err := adapter.FetchAndStore(ctx, *artifactsBase, opts)
	if err != nil {
		log.Fatalf("fetch coinex snapshot: %v\n\nTip: open https://www.coinex.com/en/reserve-proof or import:\n  por fetch coinex -summary-path ./fixtures/coinex/20260914-summary.json", err)
	}
	logResult(stored)
}

// RunOKX archives OKX public PoR artifacts.

func RunBitunix(args []string) {
	fs := flag.NewFlagSet("bitunix", flag.ExitOnError)
	artifactsBase := fs.String("artifacts", "./artifacts", "artifacts root")
	summaryPath := fs.String("summary-path", "", "import archived summary bundle JSON")
	infoFile := fs.String("info-file", "", "browser-captured queryAuditInfoByTime JSON")
	auditTime := fs.String("audit-time", "", "auditTime ISO when using -info-file")
	userProofPath := fs.String("user-proof", "", "import login-gated merkle_tree_bitunix.json")
	_ = fs.Parse(args)

	opts := exchange.FetchOpts{
		SummaryPath:   *summaryPath,
		UserProofPath: *userProofPath,
	}
	if *infoFile != "" {
		merged, err := bitunixapi.BuildSummaryBundle(*auditTime, *infoFile)
		if err != nil {
			log.Fatalf("build audit info capture: %v", err)
		}
		tmp, err := os.CreateTemp("", "bitunix-summary-*.json")
		if err != nil {
			log.Fatalf("temp file: %v", err)
		}
		tmpPath := tmp.Name()
		if _, err := tmp.Write(merged); err != nil {
			log.Fatalf("write temp: %v", err)
		}
		tmp.Close()
		defer os.Remove(tmpPath)
		opts.SummaryPath = tmpPath
		opts.ImportSource = "browser-api"
	} else if *summaryPath != "" {
		opts.ImportSource = "import"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	adapter, err := exchangereg.Get("bitunix")
	if err != nil {
		log.Fatalf("exchange: %v", err)
	}
	stored, err := adapter.FetchAndStore(ctx, *artifactsBase, opts)
	if err != nil {
		log.Fatalf("fetch bitunix snapshot: %v\n\nTip: public API is usually reachable; or import:\n  por fetch bitunix -summary-path ./fixtures/bitunix/20261103-summary.json\n  por fetch bitunix -summary-path ./summary.json -user-proof ./merkle_tree_bitunix.json", err)
	}
	logResult(stored)
}

// RunOKX archives OKX public PoR artifacts.

func RunWOOX(args []string) {
	fs := flag.NewFlagSet("woox", flag.ExitOnError)
	artifactsBase := fs.String("artifacts", "./artifacts", "artifacts root")
	summaryPath := fs.String("summary-path", "", "import archived summary bundle JSON")
	userProofPath := fs.String("user-proof", "", "optional import of login-gated user materials (stored only; verifier stubbed)")
	_ = fs.Parse(args)

	opts := exchange.FetchOpts{
		SummaryPath:   *summaryPath,
		UserProofPath: *userProofPath,
	}
	if *summaryPath != "" {
		opts.ImportSource = "import"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	adapter, err := exchangereg.Get("woox")
	if err != nil {
		log.Fatalf("exchange: %v", err)
	}
	stored, err := adapter.FetchAndStore(ctx, *artifactsBase, opts)
	if err != nil {
		log.Fatalf("fetch woox snapshot: %v\n\nTip: public API is usually reachable; or import:\n  por fetch woox -summary-path ./fixtures/woox/20261015-summary.json\nMerkle tree binaries are ~200MB and are not auto-downloaded; see %s", err, wooxapi.URLOSSVerify)
	}
	logResult(stored)
}

// RunOKX archives OKX public PoR artifacts.

func RunToobit(args []string) {
	fs := flag.NewFlagSet("toobit", flag.ExitOnError)
	artifactsBase := fs.String("artifacts", "./artifacts", "artifacts root")
	summaryPath := fs.String("summary-path", "", "import archived summary bundle JSON")
	infoFile := fs.String("info-file", "", "browser-captured latest-audit JSON")
	userProofPath := fs.String("user-proof", "", "import login-gated proof_of_reserve_{auditId}.json")
	_ = fs.Parse(args)

	opts := exchange.FetchOpts{
		SummaryPath:   *summaryPath,
		UserProofPath: *userProofPath,
	}
	if *infoFile != "" {
		merged, err := toobitapi.BuildSummaryBundle(*infoFile)
		if err != nil {
			log.Fatalf("build latest-audit capture: %v", err)
		}
		tmp, err := os.CreateTemp("", "toobit-summary-*.json")
		if err != nil {
			log.Fatalf("temp file: %v", err)
		}
		tmpPath := tmp.Name()
		if _, err := tmp.Write(merged); err != nil {
			log.Fatalf("write temp: %v", err)
		}
		tmp.Close()
		defer os.Remove(tmpPath)
		opts.SummaryPath = tmpPath
		opts.ImportSource = "browser-api"
	} else if *summaryPath != "" {
		opts.ImportSource = "import"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	adapter, err := exchangereg.Get("toobit")
	if err != nil {
		log.Fatalf("exchange: %v", err)
	}
	stored, err := adapter.FetchAndStore(ctx, *artifactsBase, opts)
	if err != nil {
		log.Fatalf("fetch toobit snapshot: %v\n\nTip: por fetch toobit\n  or import:\n  por fetch toobit -summary-path ./fixtures/toobit/20261001-summary.json", err)
	}
	logResult(stored)
}

// RunOKX archives OKX public PoR artifacts.

func RunPhemex(args []string) {
	fs := flag.NewFlagSet("phemex", flag.ExitOnError)
	artifactsBase := fs.String("artifacts", "./artifacts", "artifacts root")
	summaryPath := fs.String("summary-path", "", "import archived summary bundle JSON")
	infoFile := fs.String("info-file", "", "browser-captured porVersions entry or full list JSON")
	auditID := fs.String("audit-id", "", "select versionID (e.g. 2026-11-01 or 20261101)")
	userProofPath := fs.String("user-proof", "", "import login-gated user Merkle proof JSON if captured")
	_ = fs.Parse(args)

	opts := exchange.FetchOpts{
		SummaryPath:   *summaryPath,
		UserProofPath: *userProofPath,
		AuditID:       *auditID,
	}
	if *infoFile != "" {
		merged, err := phemexapi.BuildSummaryBundle(*infoFile)
		if err != nil {
			log.Fatalf("build porVersions capture: %v", err)
		}
		tmp, err := os.CreateTemp("", "phemex-summary-*.json")
		if err != nil {
			log.Fatalf("temp file: %v", err)
		}
		tmpPath := tmp.Name()
		if _, err := tmp.Write(merged); err != nil {
			log.Fatalf("write temp: %v", err)
		}
		tmp.Close()
		defer os.Remove(tmpPath)
		opts.SummaryPath = tmpPath
		opts.ImportSource = "browser-api"
	} else if *summaryPath != "" {
		opts.ImportSource = "import"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	adapter, err := exchangereg.Get("phemex")
	if err != nil {
		log.Fatalf("exchange: %v", err)
	}
	stored, err := adapter.FetchAndStore(ctx, *artifactsBase, opts)
	if err != nil {
		log.Fatalf("fetch phemex snapshot: %v\n\nTip: por fetch phemex\n  or import:\n  por fetch phemex -summary-path ./fixtures/phemex/20261101-summary.json", err)
	}
	logResult(stored)
}

// RunOKX archives OKX public PoR artifacts.

func RunBTCC(args []string) {
	fs := flag.NewFlagSet("btcc", flag.ExitOnError)
	artifactsBase := fs.String("artifacts", "./artifacts", "artifacts root")
	summaryPath := fs.String("summary-path", "", "import archived summary bundle JSON")
	infoFile := fs.String("info-file", "", "browser-captured reserveInfo data or full API envelope JSON")
	auditID := fs.String("audit-id", "", "select reportDate (e.g. 2026-10-15 or 20261015)")
	userProofPath := fs.String("user-proof", "", "import login-gated user Merkle leaf/proof JSON if captured")
	_ = fs.Parse(args)

	opts := exchange.FetchOpts{
		SummaryPath:   *summaryPath,
		UserProofPath: *userProofPath,
		AuditID:       *auditID,
	}
	if *infoFile != "" {
		merged, err := btccapi.BuildSummaryBundle(*infoFile)
		if err != nil {
			log.Fatalf("build reserveInfo capture: %v", err)
		}
		tmp, err := os.CreateTemp("", "btcc-summary-*.json")
		if err != nil {
			log.Fatalf("temp file: %v", err)
		}
		tmpPath := tmp.Name()
		if _, err := tmp.Write(merged); err != nil {
			log.Fatalf("write temp: %v", err)
		}
		tmp.Close()
		defer os.Remove(tmpPath)
		opts.SummaryPath = tmpPath
		opts.ImportSource = "browser-api"
	} else if *summaryPath != "" {
		opts.ImportSource = "import"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	adapter, err := exchangereg.Get("btcc")
	if err != nil {
		log.Fatalf("exchange: %v", err)
	}
	stored, err := adapter.FetchAndStore(ctx, *artifactsBase, opts)
	if err != nil {
		log.Fatalf("fetch btcc snapshot: %v\n\nTip: por fetch btcc\n  or import:\n  por fetch btcc -summary-path ./fixtures/btcc/20261015-summary.json", err)
	}
	logResult(stored)
}

// RunOKX archives OKX public PoR artifacts.

func RunOKX(args []string) {
	fs := flag.NewFlagSet("okx", flag.ExitOnError)
	artifactsBase := fs.String("artifacts", "./artifacts", "artifacts root")
	summaryPath := fs.String("summary-path", "", "import local OKX summary bundle JSON")
	walletZipPath := fs.String("wallet-zip", "", "import local OKX wallet zip")
	liabilityZipPath := fs.String("liability-zip", "", "import local OKX liability zk-STARK zip")
	skipWallet := fs.Bool("skip-wallet", false, "do not download wallet zip (~48 MB)")
	skipLiability := fs.Bool("skip-liability", false, "do not download liability zk-STARK zip")
	_ = fs.Parse(args)

	opts := exchange.FetchOpts{
		SummaryPath:      *summaryPath,
		WalletZipPath:    *walletZipPath,
		LiabilityZipPath: *liabilityZipPath,
		SkipWalletZip:    *skipWallet,
		SkipLiabilityZip: *skipLiability,
	}
	if *summaryPath != "" {
		opts.ImportSource = "import"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Minute)
	defer cancel()

	adapter, err := exchangereg.Get("okx")
	if err != nil {
		log.Fatalf("exchange: %v", err)
	}
	stored, err := adapter.FetchAndStore(ctx, *artifactsBase, opts)
	if err != nil {
		log.Fatalf("fetch okx snapshot: %v", err)
	}
	logResult(stored)
}

func logResult(stored exchange.StoreResult) {
	fetch := stored.Fetch
	log.Printf("exchange=%s auditId=%s merkleRoot=%s", fetch.Snapshot.Exchange, fetch.Snapshot.ID, fetch.Snapshot.MerkleRoot)
	log.Printf("saved raw artifacts:")
	for _, art := range fetch.Artifacts {
		log.Printf("  [%s] sha256=%s path=%s", art.Kind, art.SHA256, filepath.Join(stored.SnapshotDir, art.LocalPath))
	}
	log.Printf("snapshot dir : %s", stored.SnapshotDir)
	log.Printf("raw dir      : %s", stored.RawDir)
	log.Printf("bundles dir  : %s (run por anchor to populate)", stored.BundlesDir)
	log.Printf("metadata     : %s", filepath.Join(stored.SnapshotDir, "fetch.json"))
}
