package verifyrun

import (
	"context"

	"github.com/ardmere/ardmere/internal/verifier"
)

func runGlobalZKProofHTX(_ context.Context, in Input, _ string) verifier.Verification {
	path, sha := proofBundleFromArtifacts(in.Artifacts, in.ArtifactsDir)
	return verifier.GlobalZKProofHTX{
		SnapshotID:      in.Snapshot.ID,
		SummarySha256:   in.SummarySha,
		ProofBundlePath: path,
		ProofBundleSha:  sha,
		MerkleRoot:      in.Snapshot.MerkleRoot,
	}.Run()
}
