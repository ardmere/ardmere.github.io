package reportgen

import (
	"strings"
	"testing"

	"github.com/ardmere/ardmere/internal/verifier"
)

func TestWriteSection6FindingDetails(t *testing.T) {
	var b strings.Builder
	vers := []verifier.Verification{
		{
			VerifierID: "address-ownership",
			Verdict:    verifier.VerdictUnverifiable,
			Reason:     "No public ownership proof.",
		},
		{
			VerifierID: "onchain-balance-hot",
			Verdict:    verifier.VerdictFail,
			Findings: []verifier.Finding{
				{Subject: "0xabc", Field: "BNB@BSC#1", Claim: "10", Actual: "9", Status: verifier.VerdictFail, Note: "delta 1"},
				{Subject: "BTC|BSC", Field: "row_count", Actual: "6", Status: verifier.VerdictUnverifiable, Note: "no native verifier"},
			},
		},
	}
	writeSection6FindingDetails(&b, vers, "PR01JUN26")
	out := b.String()
	if !strings.Contains(out, "Capability and artifact gaps") {
		t.Fatal("missing stub section")
	}
	if !strings.Contains(out, "No public ownership proof") {
		t.Fatal("missing stub reason")
	}
	if !strings.Contains(out, "BNB@BSC#1") {
		t.Fatal("missing FAIL row")
	}
	if !strings.Contains(out, "no native verifier") {
		t.Fatal("missing row_count aggregate")
	}
}
