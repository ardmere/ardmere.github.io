package porrun

import (
	"context"
	"flag"
	"log"
	"time"
)

// RunAnchor fetches, verifies, writes bundles, and prints anchor calldata.
func RunAnchor(args []string) int {
	fs := flag.NewFlagSet("anchor", flag.ExitOnError)
	exchangeID := fs.String("exchange", "binance", "exchange adapter id")
	auditID := fs.String("audit-id", "", "specific snapshot audit id (okx, binance)")
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	if _, err := RunAnchorPipeline(ctx, AnchorOpts{
		ExchangeID:         *exchangeID,
		AuditID:            *auditID,
		OutDir:             *outDir,
		SkipZip:            *skipZip,
		SkipRPC:            *skipRPC,
		SkipLiability:      *skipLiab,
		DepositMaxSamples:  *depositMax,
		DepositTopKPerCoin: *depositTopK,
		DepositOnly:        *depositOnly,
		SummaryPath:        *summaryPath,
		UserProofPath:      *userProofPath,
		ZkBundlePath:       *zkBundlePath,
	}); err != nil {
		log.Fatalf("%v", err)
	}
	return 0
}
