package pipeline

import (
	"context"
	"log"

	"github.com/ardmere/ardmere/internal/bundle"
	"github.com/ardmere/ardmere/internal/exchange"
	"github.com/ardmere/ardmere/internal/por"
	"github.com/ardmere/ardmere/internal/rpc"
	"github.com/ardmere/ardmere/internal/verifier"
	"github.com/ardmere/ardmere/internal/verifyrun"
	"github.com/ardmere/ardmere/internal/walletzip"
)

// VerifyOpts configures a verifier run.
type VerifyOpts struct {
	Adapter       exchange.Adapter
	Snapshot      por.Snapshot
	Artifacts     []bundle.Artifact
	ArtifactsDir  string
	SummarySha    string
	WalletZipPath string
	WalletZipSha  string
	SkipRPC       bool
	RPC           *rpc.Client
	ETHDeposits         verifier.ETHDepositIndexer
	DepositMaxSamples  int
	DepositTopKPerCoin int
	DepositOnly        bool
	UserProofPath      string
	UserProofSha       string
}

// VerifyResult holds verifier output and optional wallet aggregate.
type VerifyResult struct {
	Verifications []verifier.Verification
	WalletAgg     *walletzip.Aggregate
	VerdictSummary uint8
	CoverageBps    uint16
}

// AggregateWallet parses the wallet zip when present.
func AggregateWallet(adapter exchange.Adapter, walletPath string) (*walletzip.Aggregate, error) {
	if walletPath == "" {
		return nil, nil
	}
	return adapter.AggregateWalletZip(walletPath)
}

// RunVerify aggregates wallet data (if any) and executes the adapter profile.
func RunVerify(ctx context.Context, opts VerifyOpts) (VerifyResult, error) {
	var walletAgg *walletzip.Aggregate
	if opts.WalletZipPath != "" {
		log.Printf("aggregating wallet zip...")
		var err error
		walletAgg, err = AggregateWallet(opts.Adapter, opts.WalletZipPath)
		if err != nil {
			return VerifyResult{}, err
		}
		log.Printf("aggregated %d rows", walletAgg.TotalRows)
	}

	profile := opts.Adapter.VerifierProfile()
	log.Printf("running verifiers from %s profile (%d shared, %d stubs)...",
		opts.Adapter.ID(), len(profile.Shared), len(profile.Stubs))

	verifications := verifyrun.Profile(ctx, verifyrun.Input{
		Exchange:      opts.Adapter.ID(),
		Profile:       profile,
		Snapshot:      opts.Snapshot,
		Artifacts:     opts.Artifacts,
		ArtifactsDir:  opts.ArtifactsDir,
		SummarySha:    opts.SummarySha,
		WalletZipPath: opts.WalletZipPath,
		WalletZipSha:  opts.WalletZipSha,
		WalletAgg:     walletAgg,
		SkipRPC:       opts.SkipRPC,
		RPC:           opts.RPC,
		ETHDeposits:         opts.ETHDeposits,
		DepositMaxSamples:   opts.DepositMaxSamples,
		DepositTopKPerCoin:  opts.DepositTopKPerCoin,
		DepositOnly:         opts.DepositOnly,
		UserProofPath:       opts.UserProofPath,
		UserProofSha:        opts.UserProofSha,
	})

	verdictSummary, coverageBps := SummarizeVerifications(verifications)
	return VerifyResult{
		Verifications:  verifications,
		WalletAgg:      walletAgg,
		VerdictSummary: verdictSummary,
		CoverageBps:    coverageBps,
	}, nil
}
