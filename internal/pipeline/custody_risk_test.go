package pipeline

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCustodyRiskScoreAtAnchor_okx508756245(t *testing.T) {
	root := FindRepoRoot()
	if root == "" {
		t.Skip("repo root not found")
	}
	path := AssessmentPath(root, "okx", "508756245")
	if _, err := os.Stat(path); err != nil {
		t.Skip("assessment not in checkout")
	}
	score, err := CustodyRiskScoreAtAnchor(path)
	if err != nil {
		t.Fatal(err)
	}
	if score == CustodyRiskUnavailable {
		t.Fatalf("expected conclusive file-only score, got %d", score)
	}
	if score < 1 || score > 100 {
		t.Fatalf("score=%d out of range", score)
	}
}

func TestCustodyRiskScoreAtAnchor_missingFile(t *testing.T) {
	score, err := CustodyRiskScoreAtAnchor(filepath.Join(t.TempDir(), "nope.json"))
	if err != nil {
		t.Fatal(err)
	}
	if score != CustodyRiskUnavailable {
		t.Fatalf("want %d, got %d", CustodyRiskUnavailable, score)
	}
}
