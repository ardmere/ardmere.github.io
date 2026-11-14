package reportgen

import (
	"strings"
	"testing"

	"github.com/ardmere/ardmere/internal/bundle"
	"github.com/ardmere/ardmere/internal/verifier"
)

func TestLookupPublicPoRApi_MEXCNo(t *testing.T) {
	info := lookupPublicPoRApi("mexc")
	if info.Provided != "no" {
		t.Fatalf("mexc publicPoRApi.provided=%q want no", info.Provided)
	}
	if info.Notes == "" {
		t.Fatal("expected registry note for mexc")
	}
	rf, ok := noPublicPoRApiRiskFlag(info)
	if !ok || rf.ID != "NO_PUBLIC_POR_API" {
		t.Fatalf("risk flag: ok=%v id=%q", ok, rf.ID)
	}
	if !strings.Contains(strings.ToLower(rf.Message), "spa") &&
		!strings.Contains(strings.ToLower(PageOnlyPoRSourceMessage), "spa") {
		t.Fatal("page-only messaging should mention SPA risk")
	}
}

func TestLookupPublicPoRApi_GateYes(t *testing.T) {
	info := lookupPublicPoRApi("gateio")
	if info.Provided != "yes" {
		t.Fatalf("gateio publicPoRApi.provided=%q want yes", info.Provided)
	}
	if _, ok := noPublicPoRApiRiskFlag(info); ok {
		t.Fatal("yes should not emit NO_PUBLIC_POR_API")
	}
}

func TestLookupPublicPoRApi_UnknownDefault(t *testing.T) {
	info := lookupPublicPoRApi("does-not-exist-exchange")
	if info.Provided != "unknown" {
		t.Fatalf("missing exchange provided=%q want unknown", info.Provided)
	}
}

func TestEvaluateMEXCIncludesNoPublicPoRApiFlag(t *testing.T) {
	art := bundle.ArtifactBundle{
		SnapshotID: "202610",
		Artifacts:  []bundle.Artifact{{Kind: "summarySnapshot"}},
	}
	vers := []verifier.Verification{
		{VerifierID: "user-merkle-proof", Verdict: verifier.VerdictUnverifiable},
	}
	eval := evaluate("mexc", art, vers)
	found := false
	for _, rf := range eval.RiskFlags {
		if rf.ID == "NO_PUBLIC_POR_API" {
			found = true
			if rf.Severity != "medium" || rf.StageEffect != "none" {
				t.Fatalf("flag meta: %+v", rf)
			}
		}
	}
	if !found {
		t.Fatal("expected NO_PUBLIC_POR_API on mexc evaluation")
	}
}

func TestPublicPoRApiAssessmentMap(t *testing.T) {
	m := publicPoRApiAssessmentMap(publicPoRApiInfo{Provided: "no"})
	if m["provided"] != "no" {
		t.Fatalf("%v", m)
	}
	notes, _ := m["notes"].(string)
	if !strings.Contains(notes, "SPA") {
		t.Fatalf("notes=%q", notes)
	}
}
