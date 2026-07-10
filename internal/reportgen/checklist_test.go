package reportgen

import (
	"testing"

	"github.com/ardmere/ardmere/internal/bundle"
	"github.com/ardmere/ardmere/internal/verifier"
)

func TestCollectChecklist_BinanceStage0(t *testing.T) {
	ctx := reportContext{
		Exchange: "binance",
		Art: bundle.ArtifactBundle{
			SnapshotID:   "PR01JUN26",
			SnapshotTime: "2026-06-01T00:00:00Z",
			Artifacts: []bundle.Artifact{
				{Kind: "bapiSnapshot"},
				{Kind: "walletZip"},
			},
		},
		Ver: bundle.VerificationBundle{
			Verifications: []verifier.Verification{
				{VerifierID: "address-ownership", Verdict: "UNVERIFIABLE"},
				{VerifierID: "global-zk-proof", Verdict: "UNVERIFIABLE"},
				{VerifierID: "third-party-attestation", Verdict: "UNVERIFIABLE"},
				{VerifierID: "onchain-balance-hot", Verdict: "FAIL"},
			},
		},
		Eval: stageEval{
			Gen:          "Gen 2",
			EffectivePoR: false,
			RiskFlags: []riskFlag{
				{ID: "NO_WALLET_OWNERSHIP_PROOF"},
				{ID: "UNVERIFIABLE", Message: "Public global proof/vk not available."},
				{ID: "OPAQUE_TRUSTED_SETUP", Message: "Trusted setup transcript is not public."},
			},
		},
	}
	items := collectChecklist(ctx)
	if len(items) != 10 {
		t.Fatalf("items: %d", len(items))
	}
	if items[4].Status != "unverifiable" {
		t.Fatalf("S1-2: %s", items[4].Status)
	}
	if items[8].Status != "not_applicable" {
		t.Fatalf("S2-1: %s", items[8].Status)
	}

	recs := collectRecommendations(ctx)
	if len(recs) == 0 {
		t.Fatal("expected recommendations")
	}
	foundTrustedSetup := false
	for _, rec := range recs {
		if rec.RelatedRiskFlags != nil && rec.RelatedRiskFlags[0] == "OPAQUE_TRUSTED_SETUP" {
			foundTrustedSetup = true
		}
	}
	if !foundTrustedSetup {
		t.Fatal("expected OPAQUE_TRUSTED_SETUP recommendation")
	}
}

func TestCollectChecklist_OKXTransparentSetup(t *testing.T) {
	ctx := reportContext{
		Exchange: "okx",
		Art: bundle.ArtifactBundle{
			SnapshotID: "508399035",
			Artifacts: []bundle.Artifact{
				{Kind: "summarySnapshot"},
				{Kind: "walletZip"},
				{Kind: "globalProofBundle"},
			},
		},
		Ver: bundle.VerificationBundle{
			Verifications: []verifier.Verification{
				{VerifierID: "address-ownership", Verdict: "PASS"},
				{VerifierID: "global-zk-proof", Verdict: "PARTIAL"},
			},
		},
		Eval: stageEval{
			Gen:          "Gen 2",
			EffectivePoR: true,
			RiskFlags: []riskFlag{
				{ID: "NO_CANONICAL_ANCHOR"},
				{ID: "HIGH_FREQUENCY_GAP"},
			},
		},
	}
	items := collectChecklist(ctx)
	var s15 checklistItem
	for _, item := range items {
		if item.ID == "S1-5" {
			s15 = item
		}
	}
	if s15.Status != "pass" {
		t.Fatalf("S1-5: want pass, got %s", s15.Status)
	}

	recs := collectRecommendations(ctx)
	for _, rec := range recs {
		for _, f := range rec.RelatedRiskFlags {
			if f == "OPAQUE_TRUSTED_SETUP" {
				t.Fatal("OKX should not get OPAQUE_TRUSTED_SETUP recommendation")
			}
		}
	}
}
