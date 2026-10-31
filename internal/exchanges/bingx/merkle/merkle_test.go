package merkle_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ardmere/ardmere/internal/exchanges/bingx/merkle"
)

func fixturePath(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "fixtures", "bingx", name)
}

func TestHashAndGenerateRoot(t *testing.T) {
	leaf := merkle.HashString("abcdefg,BTC:1.132142,ETH:2.54469871,USDT:20384,USDC:50.3,Total:49875.12599")
	// Build a tiny 1-sibling path manually: left sibling of zeros-like hash
	sib := merkle.HashString("sibling")
	root := merkle.HashString(sib + leaf)
	got, err := merkle.GenerateRoot(leaf, "l:"+sib)
	if err != nil {
		t.Fatal(err)
	}
	if got != root {
		t.Fatalf("root mismatch: got %s want %s", got, root)
	}
	// alternate hash:side encoding
	got2, err := merkle.GenerateRoot(leaf, sib+":l")
	if err != nil {
		t.Fatal(err)
	}
	if got2 != root {
		t.Fatalf("alt encoding mismatch: got %s", got2)
	}
}

func TestValidateProofFixture(t *testing.T) {
	raw, err := os.ReadFile(fixturePath("mock_user_proof.json"))
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
}
