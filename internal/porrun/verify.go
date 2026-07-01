package porrun

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ardmere/ardmere/internal/pipeline"
	"github.com/ardmere/ardmere/internal/rpc"
	"github.com/ardmere/ardmere/internal/verifier"

	"github.com/shopspring/decimal"
)

// RunVerify re-runs verifiers against a cached artifact bundle.
func RunVerify(args []string) int {
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	exchangeID := fs.String("exchange", "", "exchange adapter id (default: from bundle)")
	snapshotID := fs.String("snapshot", "PR01JUN26", "snapshot id")
	artifactsDir := fs.String("artifacts", "./artifacts", "artifact bundle directory")
	skipRPC := fs.Bool("skip-rpc", false, "skip onchain verifiers")
	depositMax := fs.Int("deposit-max-samples", 500, "max Deposit.csv rows to verify on-chain (0=env/default)")
	depositTopK := fs.Int("deposit-top-k-per-coin", 2000, "per-coin heap size when sampling Deposit.csv")
	depositOnly := fs.Bool("deposit-only", false, "run only internal-consistency + onchain-balance-deposit")
	writeBundle := fs.Bool("write-bundle", false, "write verification-bundle.json (default: on for full verify, off for --deposit-only)")
	stakeScan := fs.Bool("stake-scan", false, "BNB Stake Hub regression scan")
	_ = fs.Parse(args)

	doWriteBundle := *writeBundle || !*depositOnly

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
	defer cancel()

	loaded, err := pipeline.LoadFromBundle(*artifactsDir, *snapshotID, *exchangeID)
	if err != nil {
		log.Fatalf("load: %v", err)
	}
	snap := loaded.Snapshot
	log.Printf("exchange=%s snapshot=%s merkleRoot=%s tier=%d",
		snap.Exchange, snap.ID, snap.MerkleRoot, loaded.Adapter.Capabilities().Tier)

	rpcClient := rpc.New()
	result, err := pipeline.RunVerify(ctx, pipeline.VerifyOpts{
		Adapter:       loaded.Adapter,
		Snapshot:      snap,
		Artifacts:     loaded.Artifacts,
		ArtifactsDir:  loaded.ArtifactsDir,
		SummarySha:    loaded.Parsed.SummarySha,
		WalletZipPath: loaded.Parsed.WalletZipPath,
		WalletZipSha:  loaded.Parsed.WalletZipSha,
		UserProofPath: loaded.Parsed.UserProofPath,
		UserProofSha:  loaded.Parsed.UserProofSha,
		SkipRPC:       *skipRPC,
		RPC:           rpcClient,
		ETHDeposits:         verifier.NewEtherscanETHDepositIndexer(os.Getenv("ETHERSCAN_API_KEY")),
		DepositMaxSamples:   *depositMax,
		DepositTopKPerCoin:  *depositTopK,
		DepositOnly:         *depositOnly,
	})
	if err != nil {
		log.Fatalf("verify: %v", err)
	}

	for _, v := range result.Verifications {
		if v.VerifierID == "onchain-balance-hot" && len(v.Findings) > 0 {
			printKeyAddressResults(v.Findings)
		}
	}

	_, verBundle := pipeline.BuildBundles(snap, loaded.Artifacts, result.Verifications, result.VerdictSummary, result.CoverageBps)
	if doWriteBundle {
		outMain := filepath.Join(loaded.BundlesDir, snap.ID+".verification-bundle.json")
		outV2 := filepath.Join(loaded.BundlesDir, snap.ID+".verification-bundle.v2.json")
		mustWriteJSON(outMain, verBundle)
		mustWriteJSON(outV2, verBundle)
		log.Printf("wrote %s and %s", outMain, outV2)
	} else {
		log.Printf("skipped writing verification bundle (--deposit-only; pass --write-bundle to override)")
	}
	log.Printf("verificationBundleRoot=%s verdictSummary=0x%02x coverageBps=%d",
		verBundle.MerkleRoot, result.VerdictSummary, result.CoverageBps)

	if *stakeScan {
		runStakeScan(ctx, rpcClient)
	}
	return 0
}

func runStakeScan(ctx context.Context, rpcClient *rpc.Client) {
	log.Printf("=== BNB Stake Hub full scan (regression addresses) ===")
	for _, ka := range verifier.KeyAddresses {
		if ka.Coin != "BNB" {
			continue
		}
		liquid, used, err := rpcClient.GetBalance(ctx, rpc.NetBSC, ka.Address, ka.Height)
		if err != nil {
			log.Printf("%s liquid error: %v", ka.Label, err)
			continue
		}
		staked, err := verifier.StakeHubAccounted(ctx, rpcClient, ka.Address, ka.Height)
		if err != nil {
			log.Printf("%s stakehub error: %v", ka.Label, err)
			continue
		}
		liquidDec := decimal.NewFromBigInt(liquid, -18)
		total := liquidDec.Add(staked.Staked).Add(staked.Locked)
		log.Printf("%s %s", ka.Label, ka.Address)
		log.Printf("  block=%d provider=%s", ka.Height, used)
		log.Printf("  liquid=%s staked=%s unbonding=%s total=%s validators=%d",
			liquidDec, staked.Staked, staked.Locked, total, staked.Scanned)
		if staked.Note != "" {
			log.Printf("  note: %s", staked.Note)
		}
	}
}

func printKeyAddressResults(findings []verifier.Finding) {
	for _, ka := range verifier.KeyAddresses {
		for _, f := range findings {
			if f.Subject == ka.Address {
				log.Printf("  key %s: %s claim=%s actual=%s %v", ka.Label, f.Status, f.Claim, f.Actual, f.Components)
			}
		}
	}
}
