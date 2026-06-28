package verifier_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ardmere/ardmere/internal/verifier"
)

func bybitFixture(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "fixtures", "bybit", name)
}

func TestUserMerkleBybitMissingProof(t *testing.T) {
	v := verifier.UserMerkleBybit{SnapshotID: "2025061709"}.Run()
	if v.Verdict != verifier.VerdictUnverifiable {
		t.Fatalf("got %s", v.Verdict)
	}
}

func TestUserMerkleBybitValidFixture(t *testing.T) {
	path := bybitFixture("mock_user_merkle_tree_path_40_v5.json")
	v := verifier.UserMerkleBybit{
		ProofPath:     path,
		SummaryMerkle: "c91d5de0554f97244e4d9f8056fad70fa0cb2cdb23c290ec597a042645dcbc03",
		AuditID:       "2025061709",
		SnapshotID:    "2025061709",
	}.Run()
	if v.Verdict != verifier.VerdictPass {
		t.Fatalf("got %s reason=%s", v.Verdict, v.Reason)
	}
}

func TestUserMerkleBybitRootMismatch(t *testing.T) {
	path := bybitFixture("mock_user_merkle_tree_path_40_v5.json")
	v := verifier.UserMerkleBybit{
		ProofPath:     path,
		SummaryMerkle: "deadbeef",
		SnapshotID:    "2025061709",
	}.Run()
	if v.Verdict != verifier.VerdictFail {
		t.Fatalf("got %s", v.Verdict)
	}
}

func TestUserMerkleBybitBadFile(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "bad.json")
	_ = os.WriteFile(tmp, []byte(`{"self":{},"path":[]}`), 0o644)
	v := verifier.UserMerkleBybit{ProofPath: tmp, SnapshotID: "x"}.Run()
	if v.Verdict != verifier.VerdictFail {
		t.Fatalf("got %s", v.Verdict)
	}
}
