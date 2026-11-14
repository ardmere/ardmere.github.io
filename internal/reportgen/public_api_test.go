package reportgen

import (
	"strings"
	"testing"

	"github.com/ardmere/ardmere/internal/bundle"
	"github.com/ardmere/ardmere/internal/verifier"
)

func TestLookupPublicPoRApi_MEXCNull(t *testing.T) {
	info := lookupPublicPoRApi("mexc")
	if info.HasAPI() {
		t.Fatalf("mexc publicPoRApi=%q want null/empty", info.URL)
	}
	if info.Notes == "" {
		t.Fatal("expected registry note for mexc")
	}
	rf, ok := noPublicPoRApiRiskFlag(info)
	if !ok || rf.ID != "NO_PUBLIC_POR_API" {
		t.Fatalf("risk flag: ok=%v id=%q", ok, rf.ID)
	}
	if publicPoRApiAssessmentValue(info) != nil {
		t.Fatalf("assessment value=%v want null", publicPoRApiAssessmentValue(info))
	}
	if publicPoRApiDisplay(info) != "null" {
		t.Fatalf("display=%q", publicPoRApiDisplay(info))
	}
}

func TestLookupPublicPoRApi_GateURL(t *testing.T) {
	info := lookupPublicPoRApi("gateio")
	if !info.HasAPI() {
		t.Fatal("gateio expected documented PoR API URL")
	}
	if !strings.Contains(info.URL, "getProofOfReservesInfo") {
		t.Fatalf("gateio URL=%q", info.URL)
	}
	if _, ok := noPublicPoRApiRiskFlag(info); ok {
		t.Fatal("URL should not emit NO_PUBLIC_POR_API")
	}
	if publicPoRApiAssessmentValue(info) != info.URL {
		t.Fatalf("assessment value=%v", publicPoRApiAssessmentValue(info))
	}
}

func TestLookupPublicPoRApi_MissingIsNull(t *testing.T) {
	info := lookupPublicPoRApi("does-not-exist-exchange")
	if info.HasAPI() {
		t.Fatalf("missing exchange URL=%q want empty", info.URL)
	}
}

func TestNormalizePublicPoRApiURL(t *testing.T) {
	if got := normalizePublicPoRApiURL("https://example.com/por"); got != "https://example.com/por" {
		t.Fatalf("got %q", got)
	}
	for _, bad := range []string{"", "no", "null", "unknown", "not-a-url", "ftp://x"} {
		if got := normalizePublicPoRApiURL(bad); got != "" {
			t.Fatalf("%q -> %q want empty", bad, got)
		}
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
