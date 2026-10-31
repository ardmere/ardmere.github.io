package merkle_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ardmere/ardmere/internal/exchanges/mexc/merkle"
)

func fixturePath(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "fixtures", "mexc", name)
}

func TestValidateProofFixture(t *testing.T) {
	raw, err := os.ReadFile(fixturePath("merkle_proof.json"))
	if err != nil {
		t.Fatal(err)
	}
	ok, err := merkle.ValidateProof(raw)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected official MEXC merkle_proof.json fixture to validate")
	}
	root, err := merkle.RootHash(raw)
	if err != nil {
		t.Fatal(err)
	}
	if root != "35d01ffef2beec3ff0166c03f3d9b0cc" {
		t.Fatalf("root: got %s", root)
	}
	proofID, err := merkle.ProofID(raw)
	if err != nil {
		t.Fatal(err)
	}
	if proofID != "PR-2025-01-22" {
		t.Fatalf("proofId: got %s", proofID)
	}
}
