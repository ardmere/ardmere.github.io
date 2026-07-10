package verifyrun

import (
	"context"
	"log"

	"github.com/ardmere/ardmere/internal/verifier"
	"github.com/ardmere/ardmere/internal/walletzip"
)

type runner func(context.Context, Input, string) verifier.Verification

var (
	byRef = map[string]runner{}
	byID  = map[string]runner{}
)

func init() {
	registerID("artifact-integrity", runArtifactIntegrity)
	registerID("solvency-claim", runSolvencyClaim)
	registerID("internal-consistency", runInternalConsistency)
	registerID("onchain-balance-hot", runOnchainBalanceHot)
	registerID("onchain-balance-token", runOnchainBalanceToken)
	registerID("onchain-balance-ledger", runOnchainBalanceLedger)
	registerID("onchain-balance-deposit", runOnchainBalanceDeposit)
	registerID("btc-anchor", runBTCAnchor)

	registerRef("address-ownership@okx-1", runAddressOwnershipOKX)
	registerRef("global-zk-proof@okx-1", runGlobalZKProofOKX)
	registerRef("global-zk-proof@htx-1", runGlobalZKProofHTX)
}

func registerID(id string, fn runner) { byID[id] = fn }
func registerRef(ref string, fn runner) { byRef[ref] = fn }

func runVerifier(ctx context.Context, in Input, ref string) verifier.Verification {
	if fn, ok := byRef[ref]; ok {
		return fn(ctx, in, ref)
	}
	id, _ := verifier.ParseVerifierRef(ref)
	if fn, ok := byID[id]; ok {
		return fn(ctx, in, ref)
	}
	return verifier.Stub(ref, in.Exchange, in.Snapshot.ID, in.SummarySha, "Verifier not implemented")
}

func logVerification(v verifier.Verification) {
	switch v.VerifierID {
	case "onchain-balance-hot", "onchain-balance-token", "onchain-balance-ledger", "onchain-balance-deposit", "address-ownership":
		log.Printf("  %s: %s coverage=%.4f (%d findings)", v.VerifierID, v.Verdict, v.Coverage, len(v.Findings))
	case "global-zk-proof":
		log.Printf("  %s: %s coverage=%.4f", v.VerifierID, v.Verdict, v.Coverage)
		if v.Reason != "" {
			log.Printf("    reason: %s", v.Reason)
		}
	default:
		if len(v.Findings) > 0 {
			log.Printf("  %s: %s (%d findings)", v.VerifierID, v.Verdict, len(v.Findings))
		} else {
			log.Printf("  %s: %s", v.VerifierID, v.Verdict)
		}
	}
	if v.Reason != "" && v.VerifierID != "global-zk-proof" {
		log.Printf("    reason: %s", v.Reason)
	}
}

func runArtifactIntegrity(_ context.Context, in Input, _ string) verifier.Verification {
	return verifier.ArtifactIntegrity{
		Artifacts:    artifactRefs(in.Artifacts),
		ArtifactsDir: in.ArtifactsDir,
		SnapshotID:   in.Snapshot.ID,
	}.Run()
}

func runSolvencyClaim(_ context.Context, in Input, _ string) verifier.Verification {
	return verifier.SolvencyClaim{
		SummarySha256: in.SummarySha,
		Snapshot:      in.Snapshot,
	}.Run()
}

func runInternalConsistency(_ context.Context, in Input, ref string) verifier.Verification {
	if in.WalletZipPath == "" || in.WalletAgg == nil {
		return verifier.Stub(ref, in.Exchange, in.Snapshot.ID, in.SummarySha, "")
	}
	return verifier.InternalConsistency{
		SummarySha256:   in.SummarySha,
		WalletZipSha256: in.WalletZipSha,
		Snapshot:        in.Snapshot,
		Aggregate:       in.WalletAgg,
	}.Run()
}

func runOnchainBalanceHot(ctx context.Context, in Input, ref string) verifier.Verification {
	if in.WalletZipPath == "" || in.WalletAgg == nil || in.SkipRPC {
		reason := verifier.StubReason(in.Exchange, "onchain-balance-hot")
		if in.SkipRPC && in.WalletZipPath != "" && in.WalletAgg != nil {
			reason = "RPC queries skipped (--skip-rpc)"
		}
		return verifier.Stub(ref, in.Exchange, in.Snapshot.ID, in.SummarySha, reason)
	}
	return verifier.OnchainBalanceHot{
		Exchange:        in.Exchange,
		WalletCSV:       walletzip.WalletFileForExchange(in.Exchange),
		WalletZipPath:   in.WalletZipPath,
		WalletZipSha256: in.WalletZipSha,
		BapiSha256:      in.SummarySha,
		SnapshotID:      in.Snapshot.ID,
		RPC:             in.RPC,
		ETHDeposits:     in.ETHDeposits,
		Concurrency:     2,
	}.Run(ctx)
}

func runOnchainBalanceToken(ctx context.Context, in Input, ref string) verifier.Verification {
	if in.WalletZipPath == "" || in.WalletAgg == nil || in.SkipRPC {
		reason := verifier.StubReason(in.Exchange, "onchain-balance-token")
		if in.SkipRPC && in.WalletZipPath != "" && in.WalletAgg != nil {
			reason = "RPC queries skipped (--skip-rpc)"
		}
		return verifier.Stub(ref, in.Exchange, in.Snapshot.ID, in.SummarySha, reason)
	}
	return verifier.OnchainBalanceToken{
		Exchange:        in.Exchange,
		WalletCSV:       walletzip.WalletFileForExchange(in.Exchange),
		WalletZipPath:   in.WalletZipPath,
		WalletZipSha256: in.WalletZipSha,
		BapiSha256:      in.SummarySha,
		SnapshotID:      in.Snapshot.ID,
		RPC:             in.RPC,
		Concurrency:     4,
	}.Run(ctx)
}

func runOnchainBalanceLedger(ctx context.Context, in Input, ref string) verifier.Verification {
	if in.WalletZipPath == "" || in.WalletAgg == nil || in.SkipRPC {
		reason := verifier.StubReason(in.Exchange, "onchain-balance-ledger")
		if in.SkipRPC && in.WalletZipPath != "" && in.WalletAgg != nil {
			reason = "RPC queries skipped (--skip-rpc)"
		}
		return verifier.Stub(ref, in.Exchange, in.Snapshot.ID, in.SummarySha, reason)
	}
	return verifier.OnchainBalanceLedger{
		Exchange:        in.Exchange,
		WalletCSV:       walletzip.WalletFileForExchange(in.Exchange),
		WalletZipPath:   in.WalletZipPath,
		WalletZipSha256: in.WalletZipSha,
		BapiSha256:      in.SummarySha,
		SnapshotID:      in.Snapshot.ID,
		Concurrency:     3,
	}.Run(ctx)
}

func runAddressOwnershipOKX(_ context.Context, in Input, ref string) verifier.Verification {
	if in.WalletZipPath == "" {
		return verifier.Stub(ref, in.Exchange, in.Snapshot.ID, in.SummarySha, "")
	}
	return verifier.AddressOwnershipOKX{
		WalletZipPath:   in.WalletZipPath,
		WalletZipSha256: in.WalletZipSha,
		SnapshotID:      in.Snapshot.ID,
	}.Run()
}

func runGlobalZKProofOKX(_ context.Context, in Input, _ string) verifier.Verification {
	path, sha := proofBundleFromArtifacts(in.Artifacts, in.ArtifactsDir)
	return verifier.GlobalZKProofOKX{
		SnapshotID:      in.Snapshot.ID,
		SummarySha256:   in.SummarySha,
		ProofBundlePath: path,
		ProofBundleSha:  sha,
		MerkleRoot:      in.Snapshot.MerkleRoot,
	}.Run()
}

func runOnchainBalanceDeposit(ctx context.Context, in Input, ref string) verifier.Verification {
	if in.WalletZipPath == "" || in.WalletAgg == nil || in.SkipRPC {
		reason := "Wallet address bundle not available for deposit audit"
		if in.SkipRPC && in.WalletZipPath != "" && in.WalletAgg != nil {
			reason = "RPC queries skipped (--skip-rpc)"
		}
		return verifier.Stub(ref, in.Exchange, in.Snapshot.ID, in.SummarySha, reason)
	}
	return verifier.OnchainBalanceDeposit{
		Exchange:        in.Exchange,
		WalletZipPath:   in.WalletZipPath,
		WalletZipSha256: in.WalletZipSha,
		BapiSha256:      in.SummarySha,
		SnapshotID:      in.Snapshot.ID,
		RPC:             in.RPC,
		Concurrency:     4,
		MaxSamples:      in.DepositMaxSamples,
		TopKPerCoin:     in.DepositTopKPerCoin,
	}.Run(ctx)
}

func runBTCAnchor(ctx context.Context, in Input, ref string) verifier.Verification {
	_, ver := verifier.ParseVerifierRef(ref)
	return verifier.BTCAnchor{
		SummarySha256: in.SummarySha,
		Snapshot:      in.Snapshot,
		Version:       ver,
	}.Run(ctx)
}
