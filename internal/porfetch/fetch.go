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
	"github.com/ardmere/ardmere/internal/exchanges/bitget/bitgetapi"
	"github.com/ardmere/ardmere/internal/exchanges/bybit/bybitapi"
	"github.com/ardmere/ardmere/internal/exchanges/gateio/gateapi"
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
  okx        fetch or import OKX public PoR artifacts

Examples:
  por fetch gateio
  por fetch gateio -info-file ./info.json -coins-file ./coinList.json
  por fetch htx -zk-bundle ./public-data.zip
  por fetch bybit -summary-path ./summary.json -user-proof ./myProof.json
  por fetch bitget -summary-path ./summary.json -user-proof ./merkel_tree_bg.json
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
