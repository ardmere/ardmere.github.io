package verifyrun

import (
	"context"
	"log"

	"github.com/ardmere/ardmere/internal/verifier"
)

// Profile executes shared verifiers from the profile and emits stubs for the rest.
func Profile(ctx context.Context, in Input) []verifier.Verification {
	var out []verifier.Verification

	shared := in.Profile.Shared
	if in.DepositOnly {
		shared = filterDepositProfile(shared)
	}

	for _, ref := range shared {
		v := runVerifier(ctx, in, ref)
		logVerification(v)
		out = append(out, v)
	}

	for _, ref := range in.Profile.Stubs {
		v := verifier.Stub(ref, in.Exchange, in.Snapshot.ID, in.SummarySha, "")
		log.Printf("  %s: %s", v.VerifierID, v.Verdict)
		out = append(out, v)
	}
	return out
}

func filterDepositProfile(refs []string) []string {
	keep := map[string]bool{
		"artifact-integrity":       true,
		"internal-consistency":     true,
		"solvency-claim":           true,
		"onchain-balance-deposit": true,
	}
	var out []string
	for _, ref := range refs {
		id, _ := verifier.ParseVerifierRef(ref)
		if keep[id] {
			out = append(out, ref)
		}
	}
	return out
}
