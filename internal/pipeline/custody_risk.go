package pipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

// CustodyRiskUnavailable is the on-chain sentinel (uint8) when no conclusive score exists.
const CustodyRiskUnavailable = 255

// AssessmentPath returns the canonical public assessment JSON path for a snapshot.
func AssessmentPath(repoRoot, exchange, auditID string) string {
	return filepath.Join(repoRoot, "docs", "reports", exchange, auditID+"-assessment.json")
}

// FindRepoRoot walks upward from cwd for go.mod, or uses ARDMERE_ROOT.
func FindRepoRoot() string {
	if r := os.Getenv("ARDMERE_ROOT"); r != "" {
		return r
	}
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	for dir := cwd; dir != "/" && dir != "."; dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
	}
	return ""
}

// CustodyRiskScoreAtAnchor runs skills/por-custody-risk with --no-anchor so the score
// uses the assessment disclosure and bundle roots in the file (the same roots being
// anchored), without requiring a prior on-chain record (avoids ANCHOR_MISMATCH chicken-and-egg).
func CustodyRiskScoreAtAnchor(assessmentPath string) (uint8, error) {
	if assessmentPath == "" {
		return CustodyRiskUnavailable, nil
	}
	if _, err := os.Stat(assessmentPath); err != nil {
		return CustodyRiskUnavailable, nil
	}
	repo := FindRepoRoot()
	if repo == "" {
		return CustodyRiskUnavailable, fmt.Errorf("repo root not found (set ARDMERE_ROOT)")
	}
	script := filepath.Join(repo, "skills", "por-custody-risk", "scripts", "risk_score.py")
	cmd := exec.Command("python3", script, "--no-anchor", assessmentPath)
	cmd.Dir = repo
	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok && len(ee.Stderr) > 0 {
			return CustodyRiskUnavailable, fmt.Errorf("risk_score.py: %s", string(ee.Stderr))
		}
		return CustodyRiskUnavailable, fmt.Errorf("risk_score.py: %w", err)
	}
	var doc struct {
		Conclusive bool `json:"conclusive"`
		RiskScore  *int `json:"riskScore"`
	}
	if err := json.Unmarshal(out, &doc); err != nil {
		return CustodyRiskUnavailable, err
	}
	if !doc.Conclusive || doc.RiskScore == nil {
		return CustodyRiskUnavailable, nil
	}
	score := *doc.RiskScore
	if score < 0 {
		return CustodyRiskUnavailable, nil
	}
	if score > 100 {
		score = 100
	}
	return uint8(score), nil
}

// FormatCustodyRiskScore prints a human-readable custody risk line for anchor calldata.
func FormatCustodyRiskScore(score uint8) string {
	if score == CustodyRiskUnavailable {
		return fmt.Sprintf("%d (unavailable)", score)
	}
	return strconv.FormatUint(uint64(score), 10)
}
