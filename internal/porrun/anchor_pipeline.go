package porrun

import (
	"context"
	"fmt"
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

// AnchorOpts configures fetch → verify → bundle for one snapshot.
type AnchorOpts struct {
	ExchangeID         string
	AuditID            string
	OutDir             string
	SkipZip            bool
	SkipRPC            bool
	SkipLiability      bool
	DepositMaxSamples  int
	DepositTopKPerCoin int
	DepositOnly        bool
	SummaryPath        string
	UserProofPath      string
	ZkBundlePath       string
}

// AnchorResult is the outcome of a successful anchor pipeline.
type AnchorResult struct {
	SnapshotID  string
	SnapshotDir string
	BundlesDir  string
}

// RunAnchorPipeline fetches, verifies, writes bundles, and prints anchor calldata.
func RunAnchorPipeline(ctx context.Context, opt AnchorOpts) (AnchorResult, error) {
	empty := AnchorResult{}
	if opt.OutDir == "" {
		opt.OutDir = "./artifacts"
	}
	if opt.DepositMaxSamples == 0 {
		opt.DepositMaxSamples = 500
	}
	if opt.DepositTopKPerCoin == 0 {
		opt.DepositTopKPerCoin = 2000
	}

	adapter, err := exchangereg.Get(opt.ExchangeID)
	if err != nil {
		return empty, err
	}
	if err := os.MkdirAll(opt.OutDir, 0o755); err != nil {
		return empty, err
	}

	log.Printf("[1/5] fetching snapshot via %s adapter (tier %d)...", adapter.ID(), adapter.Capabilities().Tier)
	stored, err := adapter.FetchAndStore(ctx, opt.OutDir, exchange.FetchOpts{
		AuditID:          opt.AuditID,
		SkipWalletZip:    opt.SkipZip,
		SummaryPath:      opt.SummaryPath,
		UserProofPath:    opt.UserProofPath,
		ZkBundlePath:     opt.ZkBundlePath,
		SkipLiabilityZip: opt.SkipLiability,
	})
	if err != nil {
		return empty, fmt.Errorf("fetch snapshot: %w", err)
	}
	fetch := stored.Fetch
	snap := fetch.Snapshot
	log.Printf("      auditId=%s snapshotTime=%s merkleRoot=%s", snap.ID, snap.SnapshotTimeRaw, snap.MerkleRoot)
	log.Printf("      raw artifacts: %s", stored.RawDir)

	if fetch.WalletZipPath != "" {
		log.Printf("[2/5] wallet zip sha256=%s", fetch.WalletZipSha)
	} else if opt.SkipZip {
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
		SkipRPC:            opt.SkipRPC,
		RPC:                rpcClient,
		ETHDeposits:        verifier.NewEtherscanETHDepositIndexer(os.Getenv("ETHERSCAN_API_KEY")),
		DepositMaxSamples:  opt.DepositMaxSamples,
		DepositTopKPerCoin: opt.DepositTopKPerCoin,
		DepositOnly:        opt.DepositOnly,
	})
	if err != nil {
		return empty, fmt.Errorf("verify: %w", err)
	}

	log.Printf("[5/5] computing bundle roots...")
	artBundle, verBundle := pipeline.BuildBundles(snap, fetch.Artifacts, result.Verifications, result.VerdictSummary, result.CoverageBps)

	mustWriteJSON(filepath.Join(stored.BundlesDir, snap.ID+".artifact-bundle.json"), artBundle)
	mustWriteJSON(filepath.Join(stored.BundlesDir, snap.ID+".verification-bundle.json"), verBundle)
	mustWriteJSON(filepath.Join(stored.BundlesDir, snap.ID+".verification-bundle.v2.json"), verBundle)
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
	return AnchorResult{
		SnapshotID:  snap.ID,
		SnapshotDir: stored.SnapshotDir,
		BundlesDir:  stored.BundlesDir,
	}, nil
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
