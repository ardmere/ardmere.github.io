package merkle_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ardmere/ardmere/internal/exchanges/bitunix/merkle"
)

func fixture(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "fixtures", "bitunix", name)
}

func TestValidateOfficialSample(t *testing.T) {
	raw, err := os.ReadFile(fixture("merkle_tree_bitunix.json"))
	if err != nil {
		t.Fatal(err)
	}
	ok, err := merkle.ValidateProof(raw)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected official Bitunix merkle_tree_bitunix.json fixture to validate")
	}
	root, err := merkle.RootHash(raw)
	if err != nil {
		t.Fatal(err)
	}
	if root != "7d409a80ebbbf4cef17f3f9c5f56f6e26040f1ca285395897cf695b21391679d" {
		t.Fatalf("root=%q", root)
	}
}
