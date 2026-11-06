package merkle_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ardmere/ardmere/internal/exchanges/ourbit/merkle"
)

func fixturePath(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "fixtures", "ourbit", name)
}

func TestValidateProofREADMEFixture(t *testing.T) {
	raw, err := os.ReadFile(fixturePath("merkel_tree_ourbit.json"))
	if err != nil {
		t.Fatal(err)
	}
	ok, err := merkle.ValidateProof(raw)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected Ourbit README merkel_tree_ourbit.json fixture to validate")
	}
	root, err := merkle.RootHash(raw)
	if err != nil {
		t.Fatal(err)
	}
	if root != "df1d34512abd50712b9f5713b51fad8e" {
		t.Fatalf("root: got %s", root)
	}
	proofID, err := merkle.ProofID(raw)
	if err != nil {
		t.Fatal(err)
	}
	if proofID != "PR-2025-03-18" {
		t.Fatalf("proofId: got %s", proofID)
	}
}
