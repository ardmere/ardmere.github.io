package merkle_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ardmere/ardmere/internal/exchanges/bitget/merkle"
)

func fixturePath(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "fixtures", "bitget", name)
}

func TestValidateProofFixture(t *testing.T) {
	raw, err := os.ReadFile(fixturePath("merkel_tree_bg.json"))
	if err != nil {
		t.Fatal(err)
	}
	ok, err := merkle.ValidateProof(raw)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected fixture proof to validate")
	}
	root, err := merkle.RootHash(raw)
	if err != nil {
		t.Fatal(err)
	}
	if root != "ca89456bb711c913" {
		t.Fatalf("root: got %s", root)
	}
}
