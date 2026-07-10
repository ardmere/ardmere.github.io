package reportgen

import (
	"testing"

	"github.com/ardmere/ardmere/internal/bundle"
	"github.com/ardmere/ardmere/internal/verifier"
)

func TestEvaluateBinanceBlockedEvidenceNames(t *testing.T) {
	art := bundle.ArtifactBundle{
		SnapshotID: "PR01JUN26",
		Artifacts: []bundle.Artifact{
			{Kind: "bapiSnapshot"},
			{Kind: "walletZip"},
		},
	}
	vers := []verifier.Verification{
		{VerifierID: "global-zk-proof", Verdict: verifier.VerdictUnverifiable},
	}
	eval := evaluate("binance", art, vers)
	if len(eval.Blocked) != 3 {
		t.Fatalf("blocked: %d", len(eval.Blocked))
	}
	want := []string{"wallet_ownership_proof", "global_proof.csv, verifying_key", "trusted_setup_transcript"}
	for i, name := range want {
		if eval.Blocked[i].Evidence != name {
			t.Fatalf("blocked[%d].Evidence=%q want %q", i, eval.Blocked[i].Evidence, name)
		}
	}
}
