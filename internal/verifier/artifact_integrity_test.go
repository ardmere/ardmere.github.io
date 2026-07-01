package verifier

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func TestArtifactIntegrityPass(t *testing.T) {
	dir := t.TempDir()
	body := []byte(`{"code":0}`)
	sum := sha256.Sum256(body)
	sumHex := hex.EncodeToString(sum[:])
	path := filepath.Join(dir, sumHex+".json")
	if err := os.WriteFile(path, body, 0o644); err != nil {
		t.Fatal(err)
	}

	v := ArtifactIntegrity{
		Artifacts: []ArtifactRef{{
			Kind:      "summarySnapshot",
			SHA256:    sumHex,
			LocalPath: path,
		}},
		ArtifactsDir: dir,
		SnapshotID:   "test",
	}.Run()

	if v.Verdict != VerdictPass {
		t.Fatalf("verdict=%s findings=%+v", v.Verdict, v.Findings)
	}
}

func TestArtifactIntegrityFailMismatch(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	v := ArtifactIntegrity{
		Artifacts: []ArtifactRef{{
			Kind:      "summarySnapshot",
			SHA256:    "deadbeef",
			LocalPath: path,
		}},
		ArtifactsDir: dir,
		SnapshotID:   "test",
	}.Run()

	if v.Verdict != VerdictFail {
		t.Fatalf("verdict=%s", v.Verdict)
	}
}
