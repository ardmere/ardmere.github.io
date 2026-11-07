package merkle

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateDemoProofs(t *testing.T) {
	// Prefer in-repo fixture; fall back to relative path used by tests from repo root.
	candidates := []string{
		filepath.Join("..", "..", "..", "..", "fixtures", "toobit", "user_1_demo_proof.json"),
		"fixtures/toobit/user_1_demo_proof.json",
	}
	var raw []byte
	var err error
	for _, p := range candidates {
		raw, err = os.ReadFile(p)
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Fatalf("read demo proof: %v", err)
	}
	ok, err := ValidateProof(raw)
	if err != nil {
		t.Fatalf("ValidateProof: %v", err)
	}
	if !ok {
		t.Fatal("expected demo proof to validate")
	}
	root, err := RootHash(raw)
	if err != nil {
		t.Fatalf("RootHash: %v", err)
	}
	want := "3586fdd20930df70e214ca0157c68bb6226289824f05f92d5e0465eb5297f727"
	if root != want {
		t.Fatalf("root = %s, want %s", root, want)
	}
	audit, err := AuditID(raw)
	if err != nil {
		t.Fatalf("AuditID: %v", err)
	}
	if audit != "test_demo_v1" {
		t.Fatalf("audit = %q, want test_demo_v1", audit)
	}
}
