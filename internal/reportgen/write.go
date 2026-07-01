package reportgen

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ardmere/ardmere/internal/bundle"
	"github.com/ardmere/ardmere/internal/verifier"
)

// Options configures report generation for one snapshot.
type Options struct {
	Exchange   string
	SnapshotID string
	Artifacts  string
	ReportsDir string
}

// Write emits assessment JSON and transparency report markdown.
func Write(ctx context.Context, opt Options) error {
	_ = ctx
	if opt.ReportsDir == "" {
		opt.ReportsDir = "./docs/reports"
	}
	if opt.Artifacts == "" {
		opt.Artifacts = "./artifacts"
	}

	dir := snapshotDir(opt.Artifacts, opt.Exchange, opt.SnapshotID)
	art, err := loadArtifactBundle(dir, opt.SnapshotID)
	if err != nil {
		return fmt.Errorf("artifact bundle: %w", err)
	}
	ver, err := loadVerificationBundle(dir, opt.SnapshotID)
	if err != nil {
		return fmt.Errorf("verification bundle: %w", err)
	}

	eval := evaluate(opt.Exchange, art, ver.Verifications)
	assessment := buildAssessment(opt.Exchange, art, ver, eval)
	reportMD := buildReportMarkdown(opt.Exchange, art, ver, eval)

	outDir := filepath.Join(opt.ReportsDir, opt.Exchange)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	assessPath := filepath.Join(outDir, opt.SnapshotID+"-assessment.json")
	assessBytes, err := json.MarshalIndent(assessment, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(assessPath, append(assessBytes, '\n'), 0o644); err != nil {
		return err
	}

	reportPath := filepath.Join(outDir, opt.SnapshotID+"-transparency-report.md")
	if err := os.WriteFile(reportPath, []byte(reportMD), 0o644); err != nil {
		return err
	}
	return nil
}

func buildAssessment(exchangeID string, art bundle.ArtifactBundle, ver bundle.VerificationBundle, eval stageEval) map[string]any {
	vr := make([]map[string]any, 0, len(ver.Verifications))
	for _, v := range ver.Verifications {
		vr = append(vr, map[string]any{
			"verifierId":     v.VerifierID,
			"version":        v.Version,
			"verdict":        strings.ToUpper(string(v.Verdict)),
			"coverage":       v.Coverage,
			"reason":         v.Reason,
			"inputArtifacts": v.InputArtifacts,
		})
	}
	present := make([]map[string]any, 0, len(art.Artifacts))
	for _, a := range art.Artifacts {
		present = append(present, map[string]any{
			"kind":         a.Kind,
			"sha256":       a.SHA256,
			"url":          a.URL,
			"availability": "public",
		})
	}
	return map[string]any{
		"schema":               "ardmere/exchange-assessment@1",
		"generatedAt":          nowRFC3339(),
		"methodologyVersion":   "por-transparency-framework@2026-07-01",
		"exchange":             exchangeID,
		"snapshot": map[string]any{
			"id":                 art.SnapshotID,
			"periodSeq":          art.PeriodSeq,
			"snapshotTime":       art.SnapshotTime,
			"btcBlockHeight":     art.BTCBlockHeight,
			"exchangeMerkleRoot": art.ExchangeMerkleRoot,
		},
		"porStage":        eval.Stage,
		"gen":             eval.Gen,
		"evidenceLevel":   eval.EvidenceLevel,
		"confidence":      eval.Confidence,
		"effectivePoR":    eval.EffectivePoR,
		"summary": map[string]any{
			"headline":  eval.Headline,
			"rationale": eval.Rationale,
			"scope":     fmt.Sprintf("%s snapshot %s", exchangeID, art.SnapshotID),
		},
		"stageBlockedReasons": eval.Blocked,
		"riskFlags":           eval.RiskFlags,
		"artifacts": map[string]any{
			"present":                  present,
			"artifactBundleRoot":       art.MerkleRoot,
			"verificationBundleRoot":   ver.MerkleRoot,
		},
		"verifierResults": vr,
	}
}

func buildReportMarkdown(exchangeID string, art bundle.ArtifactBundle, ver bundle.VerificationBundle, eval stageEval) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s %s PoR Transparency Report\n\n", strings.ToUpper(exchangeID), art.SnapshotID)
	b.WriteString("> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  \n")
	fmt.Fprintf(&b, "> Assessment JSON: [%s-assessment.json](./%s-assessment.json)\n\n", art.SnapshotID, art.SnapshotID)

	b.WriteString("## 1. Summary\n\n| Field | Value |\n| --- | --- |\n")
	fmt.Fprintf(&b, "| Exchange | `%s` |\n", exchangeID)
	fmt.Fprintf(&b, "| Snapshot | `%s` |\n", art.SnapshotID)
	fmt.Fprintf(&b, "| Snapshot time | `%s` |\n", art.SnapshotTime)
	fmt.Fprintf(&b, "| PoR Stage | `%s — %s` |\n", eval.Stage, eval.StageName)
	fmt.Fprintf(&b, "| Gen / Evidence | `%s / %s` |\n", eval.Gen, eval.EvidenceLevel)
	fmt.Fprintf(&b, "| Confidence | `%s` |\n", eval.Confidence)
	fmt.Fprintf(&b, "| Effective PoR | `%v` |\n\n", eval.EffectivePoR)
	b.WriteString(eval.Headline + "\n\n")
	b.WriteString(eval.Rationale + "\n\n")

	b.WriteString("## 2. Stage Decision\n\n")
	if eval.EffectivePoR {
		b.WriteString("### Stage 1, Effective PoR\n\n")
	} else {
		b.WriteString("### Stage 0, not effective PoR\n\n")
	}
	if len(eval.Blocked) > 0 {
		b.WriteString("| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |\n| --- | --- | --- | --- |\n")
		for _, row := range eval.Blocked {
			fmt.Fprintf(&b, "| — | `%s` | %s | %s |\n", row.RiskFlag, row.MaxStage, row.Reason)
		}
		b.WriteString("\n")
	}

	b.WriteString("## 5. Verifier Evidence\n\n")
	b.WriteString("| Verifier | Version | Verdict | Coverage | Summary |\n| --- | --- | --- | --- | --- |\n")
	for _, v := range ver.Verifications {
		fmt.Fprintf(&b, "| `%s` | `%s` | `%s` | `%.4f` | %s |\n",
			v.VerifierID, v.Version, strings.ToUpper(string(v.Verdict)), v.Coverage, escapeCell(v.Reason))
	}
	b.WriteString("\n")

	b.WriteString("## 8. Boundary\n\n")
	b.WriteString("This report evaluates public PoR artifacts and reproducibility. It does not evaluate the exchange's overall financial health, corporate governance, internal controls, off-chain assets, legal compliance, or complete off-balance-sheet liabilities.\n\n")
	b.WriteString("Missing data must be marked as `UNVERIFIABLE`, not treated as `PASS`.\n")
	return b.String()
}

func escapeCell(s string) string {
	return strings.ReplaceAll(s, "|", "\\|")
}

// Ensure verifier types referenced.
var _ verifier.Verdict
