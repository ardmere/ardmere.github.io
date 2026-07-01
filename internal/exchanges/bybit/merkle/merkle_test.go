package merkle_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ardmere/ardmere/internal/exchanges/bybit/merkle"
)

func fixturePath(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "fixtures", "bybit", name)
}

func TestValidateProofV5Mock(t *testing.T) {
	raw, err := os.ReadFile(fixturePath("mock_user_merkle_tree_path_40_v5.json"))
	if err != nil {
		t.Fatal(err)
	}
	ok, err := merkle.ValidateProofV5(raw)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected mock v5 proof to validate")
	}
	root, err := merkle.RootHash(raw)
	if err != nil {
		t.Fatal(err)
	}
	want := "c91d5de0554f97244e4d9f8056fad70fa0cb2cdb23c290ec597a042645dcbc03"
	if root != want {
		t.Fatalf("root: got %s want %s", root, want)
	}
}
