package verifier

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"time"
)

// ArtifactRef is a minimal artifact descriptor for integrity checks.
type ArtifactRef struct {
	Kind      string
	SHA256    string
	LocalPath string
}

// ArtifactIntegrity checks that each artifact's recorded sha256 matches the
// content on disk at verification time.
type ArtifactIntegrity struct {
	Artifacts    []ArtifactRef
	ArtifactsDir string
	SnapshotID   string
}

func (v ArtifactIntegrity) Run() Verification {
	out := Verification{
		VerifierID: "artifact-integrity",
		Version:    "1.0",
		SnapshotID: v.SnapshotID,
		VerifiedAt: time.Now().UTC(),
		Coverage:   1.0,
	}

	if len(v.Artifacts) == 0 {
		out.Verdict = VerdictFail
		out.Reason = "no artifacts in bundle"
		return out
	}

	allPass := true
	for _, art := range v.Artifacts {
		path := resolveArtifactPath(v.ArtifactsDir, art.LocalPath)
		sumHex, err := fileSHA256(path)
		st := VerdictPass
		note := ""
		if err != nil {
			st = VerdictFail
			allPass = false
			note = err.Error()
		} else if sumHex != art.SHA256 {
			st = VerdictFail
			allPass = false
			note = "recorded sha256 mismatch"
		}
		out.InputArtifacts = appendUnique(out.InputArtifacts, art.SHA256)
		out.Findings = append(out.Findings, Finding{
			Subject: art.Kind,
			Field:   "sha256",
			Claim:   art.SHA256,
			Actual:  sumHex,
			Status:  st,
			Note:    note,
		})
	}

	if allPass {
		out.Verdict = VerdictPass
	} else {
		out.Verdict = VerdictFail
	}
	return out
}

func resolveArtifactPath(artifactsDir, localPath string) string {
	if localPath == "" {
		return ""
	}
	if filepath.IsAbs(localPath) {
		return localPath
	}
	if p := filepath.Join(artifactsDir, localPath); fileExists(p) {
		return p
	}
	if _, err := os.Stat(localPath); err == nil {
		return localPath
	}
	return filepath.Join(artifactsDir, filepath.Base(localPath))
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func fileSHA256(path string) (string, error) {
	if path == "" {
		return "", os.ErrNotExist
	}
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func appendUnique(list []string, v string) []string {
	for _, x := range list {
		if x == v {
			return list
		}
	}
	return append(list, v)
}
