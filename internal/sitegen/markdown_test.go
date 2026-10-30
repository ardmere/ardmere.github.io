package sitegen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ardmere/ardmere/internal/exchangeregistry"
)

func TestWriteSnapshotHistory(t *testing.T) {
	root := t.TempDir()
	reports := filepath.Join(root, "docs", "reports")
	regPath := filepath.Join(root, "config", "exchanges", "registry.yaml")
	outPath := filepath.Join(root, "docs", "reports", "snapshot-history.md")
	notesPath := filepath.Join(root, "config", "site", "snapshot-notes.yaml")

	setupRegistry(t, regPath)
	writeFullAssessment(t, reports, "okx", "508399035", assessmentFixture{
		snapshotTime: "2026-06-18T16:00:00Z",
		stage:        "Stage 1",
		effectivePoR: true,
		generatedAt:  "2026-07-08T12:00:00Z",
		kinds:        []string{"summarySnapshot", "walletZip", "globalProofBundle"},
	})
	writeFullAssessment(t, reports, "gateio", "20260618", assessmentFixture{
		snapshotTime: "2026-06-18T00:00:00Z",
		stage:        "Stage 0",
		generatedAt:  "2026-07-08T12:00:00Z",
		kinds:        []string{"summarySnapshot"},
	})

	if err := os.MkdirAll(filepath.Join(root, "config", "site"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(notesPath, []byte(`overrides:
  gateio/20260618: Page-scrape import (official merkle root)
footers:
  - |-
    HTX June 2026 pending import.
`), 0o644); err != nil {
		t.Fatal(err)
	}
	oldWD, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })

	reg, err := exchangeregistry.Load(regPath)
	if err != nil {
		t.Fatal(err)
	}
	rows, err := LoadAssessments(reports, reg)
	if err != nil {
		t.Fatal(err)
	}
	if err := WriteSnapshotHistory(outPath, rows, reg); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	s := string(content)
	for _, want := range []string{
		"508399035",
		"**Stage 1**",
		"Page-scrape import (official merkle root)",
		"HTX June 2026 pending import.",
		"go run ./cmd/por site",
	} {
		if !strings.Contains(s, want) {
			t.Fatalf("missing %q in:\n%s", want, s)
		}
	}
}

func TestWriteExchangeComparison(t *testing.T) {
	root := t.TempDir()
	reports := filepath.Join(root, "docs", "reports")
	regPath := filepath.Join(root, "config", "exchanges", "registry.yaml")
	outPath := filepath.Join(root, "docs", "reports", "exchange-comparison.md")

	setupRegistry(t, regPath)
	writeFullAssessment(t, reports, "binance", "PR01JUL26", assessmentFixture{
		snapshotTime:  "2026-07-01T00:00:00Z",
		stage:         "Stage 0",
		gen:           "Gen 2",
		evidenceLevel: "E2",
		confidence:    "medium",
		generatedAt:   "2026-07-08T12:00:00Z",
		headline:      "public summary and wallet list",
		periodSeq:     44,
		kinds:         []string{"bapiSnapshot", "walletZip"},
	})
	writeFullAssessment(t, reports, "okx", "508399035", assessmentFixture{
		snapshotTime:  "2026-06-18T16:00:00Z",
		stage:         "Stage 1",
		gen:           "Gen 2",
		evidenceLevel: "E2",
		confidence:    "medium",
		effectivePoR:  true,
		generatedAt:   "2026-07-08T12:00:00Z",
		headline:      "public reserves, ownership proof, global zk proof",
		kinds:         []string{"summarySnapshot", "walletZip", "globalProofBundle"},
	})

	reg, err := exchangeregistry.Load(regPath)
	if err != nil {
		t.Fatal(err)
	}
	rows, err := LoadAssessments(reports, reg)
	if err != nil {
		t.Fatal(err)
	}
	if err := WriteExchangeComparison(outPath, rows, reg); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	s := string(content)
	for _, want := range []string{
		"PR01JUL26",
		"508399035",
		"[v1 sigs](https://github.com/okx/proof-of-reserves)",
		"[Full report](./binance/PR01JUL26-transparency-report.md)",
		"public summary and wallet list",
	} {
		if !strings.Contains(s, want) {
			t.Fatalf("missing %q in:\n%s", want, s)
		}
	}
}

type assessmentFixture struct {
	snapshotTime  string
	stage         string
	gen           string
	evidenceLevel string
	confidence    string
	effectivePoR  bool
	generatedAt   string
	headline      string
	periodSeq     int
	kinds         []string
}

func setupRegistry(t *testing.T, regPath string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(regPath), 0o755); err != nil {
		t.Fatal(err)
	}
	yaml := `version: 1
updated: "2026-07-01"
exchanges:
  - id: binance
    name: Binance
    porPage: https://example.com/binance
    upstreamRepos:
      - url: https://github.com/binance/zkmerkle-proof-of-solvency
        role: global_zk_verifier_spec
  - id: okx
    name: OKX
    porPage: https://example.com/okx
    upstreamRepos:
      - url: https://github.com/okx/proof-of-reserves
        role: address_ownership_verifier
      - url: https://github.com/okx/proof-of-reserves-v2
        role: global_zk_verifier_spec
  - id: gateio
    name: Gate.io
    porPage: https://example.com/gateio
    upstreamRepos:
      - url: https://github.com/gateio/proof-of-reserves
        role: verifier_toolchain
`
	if err := os.WriteFile(regPath, []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeFullAssessment(t *testing.T, reportsDir, exchange string, snapshotID string, fx assessmentFixture) {
	t.Helper()
	dir := filepath.Join(reportsDir, exchange)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	md := filepath.Join(dir, snapshotID+"-transparency-report.md")
	if err := os.WriteFile(md, []byte("# test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	kindsJSON := ""
	for i, k := range fx.kinds {
		if i > 0 {
			kindsJSON += ","
		}
		kindsJSON += `{"kind":"` + k + `"}`
	}
	payload := `{
  "exchange": "` + exchange + `",
  "porStage": "` + fx.stage + `",
  "gen": "` + fx.gen + `",
  "evidenceLevel": "` + fx.evidenceLevel + `",
  "confidence": "` + fx.confidence + `",
  "effectivePoR": ` + boolJSON(fx.effectivePoR) + `,
  "generatedAt": "` + fx.generatedAt + `",
  "snapshot": {
    "id": "` + snapshotID + `",
    "snapshotTime": "` + fx.snapshotTime + `",
    "periodSeq": ` + intJSON(fx.periodSeq) + `
  },
  "summary": {
    "headline": "` + fx.headline + `"
  },
  "artifacts": {
    "present": [` + kindsJSON + `]
  }
}`
	jsonPath := filepath.Join(dir, snapshotID+"-assessment.json")
	if err := os.WriteFile(jsonPath, []byte(payload), 0o644); err != nil {
		t.Fatal(err)
	}
}

func boolJSON(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

func intJSON(v int) string {
	return fmt.Sprintf("%d", v)
}
