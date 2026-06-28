package verifier_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ardmere/ardmere/internal/verifier"
)

func bitgetFixture(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "fixtures", "bitget", name)
}

func TestUserMerkleBitgetMissingProof(t *testing.T) {
	v := verifier.UserMerkleBitget{SnapshotID: "202605"}.Run()
	if v.Verdict != verifier.VerdictUnverifiable {
		t.Fatalf("got %s", v.Verdict)
	}
}

func TestUserMerkleBitgetValidFixture(t *testing.T) {
	v := verifier.UserMerkleBitget{
		ProofPath:     bitgetFixture("merkel_tree_bg.json"),
		SummaryMerkle: "ca89456bb711c913",
		AuditID:       "202605",
		SnapshotID:    "202605",
	}.Run()
	if v.Verdict != verifier.VerdictPass {
		t.Fatalf("got %s reason=%s", v.Verdict, v.Reason)
	}
	if len(v.Findings) < 3 {
		t.Fatalf("expected validation, hash-security, and root-bind findings: %+v", v.Findings)
	}
	if v.Findings[1].Status != verifier.VerdictWarn {
		t.Fatalf("expected hash-security WARN, got %+v", v.Findings[1])
	}
}

func TestUserMerkleBitgetRootMismatch(t *testing.T) {
	v := verifier.UserMerkleBitget{
		ProofPath:     bitgetFixture("merkel_tree_bg.json"),
		SummaryMerkle: "deadbeef",
		SnapshotID:    "202605",
	}.Run()
	if v.Verdict != verifier.VerdictFail {
		t.Fatalf("got %s", v.Verdict)
	}
}
