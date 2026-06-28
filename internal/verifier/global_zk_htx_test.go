package verifier

import (
	"os"
	"testing"
)

func TestGlobalZKHTXStructurePartial(t *testing.T) {
	path := os.Getenv("HTX_PUBLIC_DATA_ZIP")
	if path == "" {
		path = "/tmp/htx-pub/public-data.zip"
	}
	if _, err := os.Stat(path); err != nil {
		t.Skipf("public-data.zip not available: %v", err)
	}
	root := "2c51a40e991be73032f034c61b5d6a620c6da1fc9a7e0fdf104ca396e45d9ee3"
	v := GlobalZKProofHTX{
		SnapshotID:      "20230910",
		SummarySha256:   "abc",
		ProofBundlePath: path,
		ProofBundleSha:  "def",
		MerkleRoot:      root,
	}.Run()
	if v.Verdict != VerdictPartial {
		t.Fatalf("expected PARTIAL without HTX_ZK_VERIFIER, got %s reason=%s", v.Verdict, v.Reason)
	}
}
