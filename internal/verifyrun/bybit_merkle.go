package verifyrun

import (
	"context"

	"github.com/ardmere/ardmere/internal/verifier"
)

func init() {
	registerRef("user-merkle-proof@bybit-1", runUserMerkleProofBybit)
}

func runUserMerkleProofBybit(_ context.Context, in Input, _ string) verifier.Verification {
	path, sha := userProofFromArtifacts(in.Artifacts, in.ArtifactsDir)
	if path == "" && in.UserProofPath != "" {
		path = in.UserProofPath
		sha = in.UserProofSha
	}
	return verifier.UserMerkleBybit{
		ProofPath:     path,
		ProofSha256:   sha,
		SummaryMerkle: in.Snapshot.MerkleRoot,
		AuditID:       in.Snapshot.ID,
		SnapshotID:    in.Snapshot.ID,
	}.Run()
}
