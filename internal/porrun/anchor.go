package porrun

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ardmere/ardmere/internal/exchange"
	"github.com/ardmere/ardmere/internal/exchangereg"
	"github.com/ardmere/ardmere/internal/pipeline"
	"github.com/ardmere/ardmere/internal/rpc"
	"github.com/ardmere/ardmere/internal/verifier"
)

// RunAnchor fetches, verifies, writes bundles, and prints anchor calldata.
func RunAnchor(args []string) int {
	fs := flag.NewFlagSet("anchor", flag.ExitOnError)
	exchangeID := fs.String("exchange", "binance", "exchange adapter id")
	skipZip := fs.Bool("skip-zip", false, "do not download wallet ZIP")
	skipRPC := fs.Bool("skip-rpc", false, "skip onchain verifiers")
	depositMax := fs.Int("deposit-max-samples", 500, "max Deposit.csv rows to verify on-chain")
	depositTopK := fs.Int("deposit-top-k-per-coin", 2000, "per-coin heap when sampling Deposit.csv")
	depositOnly := fs.Bool("deposit-only", false, "run only internal-consistency + onchain-balance-deposit")
	summaryPath := fs.String("summary-path", "", "import local summary JSON")
	userProofPath := fs.String("user-proof", "", "import login-gated user Merkle proof JSON")
	zkBundlePath := fs.String("zk-bundle", "", "import local Gate zk bundle tar.gz")
	skipLiab := fs.Bool("skip-liability", false, "skip OKX liability zip")
	outDir := fs.String("out", "./artifacts", "artifacts root")
	_ = fs.Parse(args)

	adapter, err := exchangereg.Get(*exchangeID)
	if err != nil {
		log.Fatalf("exchange: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		log.Fatalf("mkdir %s: %v", *outDir, err)
	}

	log.Printf("[1/5] fetching snapshot via %s adapter (tier %d)...", adapter.ID(), adapter.Capabilities().Tier)
	stored, err := adapter.FetchAndStore(ctx, *outDir, exchange.FetchOpts{
		SkipWalletZip:    *skipZip,
		SummaryPath:      *summaryPath,
		UserProofPath:    *userProofPath,
		ZkBundlePath:     *zkBundlePath,
		SkipLiabilityZip: *skipLiab,
	})
	if err != nil {
		log.Fatalf("fetch snapshot: %v", err)
	}
	fetch := stored.Fetch
	snap := fetch.Snapshot
	log.Printf("      auditId=%s snapshotTime=%s merkleRoot=%s", snap.ID, snap.SnapshotTimeRaw, snap.MerkleRoot)
	log.Printf("      raw artifacts: %s", stored.RawDir)

	if fetch.WalletZipPath != "" {
		log.Printf("[2/5] wallet zip sha256=%s", fetch.WalletZipSha)
	} else if *skipZip {
		log.Printf("[2/5] --skip-zip set; zip-dependent verifiers will stub")
	} else {
		log.Printf("[2/5] no wallet bundle; zip-dependent verifiers will stub")
	}

	log.Printf("[3/5] running verifiers...")
	rpcClient := rpc.New()
	result, err := pipeline.RunVerify(ctx, pipeline.VerifyOpts{
		Adapter:            adapter,
		Snapshot:           snap,
		Artifacts:          fetch.Artifacts,
		ArtifactsDir:       stored.SnapshotDir,
		SummarySha:         fetch.SummarySha,
		WalletZipPath:      fetch.WalletZipPath,
		WalletZipSha:       fetch.WalletZipSha,
		SkipRPC:            *skipRPC,
		RPC:                rpcClient,
		ETHDeposits:        verifier.NewEtherscanETHDepositIndexer(os.Getenv("ETHERSCAN_API_KEY")),
		DepositMaxSamples:  *depositMax,
		DepositTopKPerCoin: *depositTopK,
		DepositOnly:        *depositOnly,
	})
	if err != nil {
		log.Fatalf("verify: %v", err)
	}

	log.Printf("[5/5] computing bundle roots...")
	artBundle, verBundle := pipeline.BuildBundles(snap, fetch.Artifacts, result.Verifications, result.VerdictSummary, result.CoverageBps)

	mustWriteJSON(filepath.Join(stored.BundlesDir, snap.ID+".artifact-bundle.json"), artBundle)
	mustWriteJSON(filepath.Join(stored.BundlesDir, snap.ID+".verification-bundle.json"), verBundle)
	mustWriteJSON(filepath.Join(stored.BundlesDir, snap.ID+".anchor.json"), map[string]any{
		"exchange":               snap.Exchange,
		"snapshotId":             snap.ID,
		"periodSeq":              snap.PeriodSeq,
		"snapshotTime":           snap.SnapshotTime.UTC().Format(time.RFC3339),
		"btcBlockHeight":         pipeline.BTCBlockHeight(snap),
		"exchangeMerkleRoot":     snap.MerkleRoot,
		"artifactBundleRoot":     artBundle.MerkleRoot,
		"verificationBundleRoot": verBundle.MerkleRoot,
		"verdictSummary":         result.VerdictSummary,
		"coverageBps":            result.CoverageBps,
	})
	log.Printf("      bundles written under %s/", stored.BundlesDir)

	pipeline.PrintAnchorCalldata(pipeline.AnchorParamsFrom(
		snap,
		pipeline.ParseHex32(snap.MerkleRoot),
		artBundle.MerkleRoot,
		verBundle.MerkleRoot,
		result.VerdictSummary,
		result.CoverageBps,
	))
	return 0
}

func mustWriteJSON(path string, v any) {
	if err := pipeline.WriteJSON(path, v); err != nil {
		log.Fatalf("%v", err)
	}
}

func validateAnchorSelector() {
	pipeline.ValidateAnchorSelector()
}

// ValidateAnchorSelector checks anchor calldata selector at startup.
func ValidateAnchorSelector() { validateAnchorSelector() }
