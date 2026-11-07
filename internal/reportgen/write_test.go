package reportgen

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// writeBundleFixture lays out artifacts/<exchange>/<snapshotID>/bundles with a
// minimal artifact + verification bundle so Write() can render a report.
func writeBundleFixture(t *testing.T, artifactsBase, exchangeID, snapshotID, snapshotTime string) {
	t.Helper()
	bundlesDir := filepath.Join(artifactsBase, exchangeID, snapshotID, "bundles")
	if err := os.MkdirAll(bundlesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	art := `{
  "exchange": "` + exchangeID + `",
  "snapshotId": "` + snapshotID + `",
  "snapshotTime": "` + snapshotTime + `",
  "builtAt": "2026-01-01T00:00:00Z",
  "kind": "artifact-bundle",
  "artifacts": [
    {"kind": "summarySnapshot", "sha256": "aa", "url": "https://example.com", "fetchedAt": "2026-01-01T00:00:00Z"}
  ],
  "merkleRoot": "0x01"
}
`
	ver := `{
  "exchange": "` + exchangeID + `",
  "snapshotId": "` + snapshotID + `",
  "snapshotTime": "` + snapshotTime + `",
  "builtAt": "2026-01-01T00:00:00Z",
  "kind": "verification-bundle",
  "verifications": [],
  "merkleRoot": "0x02"
}
`
	if err := os.WriteFile(filepath.Join(bundlesDir, snapshotID+".artifact-bundle.json"), []byte(art), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bundlesDir, snapshotID+".verification-bundle.v2.json"), []byte(ver), 0o644); err != nil {
		t.Fatal(err)
	}
}

var tableRowPattern = regexp.MustCompile(`\| ([^|]+?) \| ` + "`" + `?([^` + "`" + `|]*)` + "`" + `? \|`)

// frequencySection returns section 3 of a generated report as field -> value.
func frequencySection(t *testing.T, reportMD string) map[string]string {
	t.Helper()
	start := strings.Index(reportMD, "## 3. Frequency and Freshness")
	end := strings.Index(reportMD, "## 4.")
	if start < 0 || end < 0 || end < start {
		t.Fatalf("section 3 not found in report:\n%s", reportMD)
	}
	fields := map[string]string{}
	for _, m := range tableRowPattern.FindAllStringSubmatch(reportMD[start:end], -1) {
		fields[strings.TrimSpace(m[1])] = strings.TrimSpace(m[2])
	}
	return fields
}

// TestWriteReportIncludesItsOwnSnapshot is the regression test for the ordering
// fault: `por report` rendered the page before the assessment JSON existed, so
// the Frequency table omitted the snapshot the page was about.
func TestWriteReportIncludesItsOwnSnapshot(t *testing.T) {
	root := t.TempDir()
	artifactsBase := filepath.Join(root, "artifacts")
	reportsDir := filepath.Join(root, "reports")

	writeBundleFixture(t, artifactsBase, "binance", "PR01APR26", "2026-04-01T00:00:00Z")
	writeBundleFixture(t, artifactsBase, "binance", "PR01MAY26", "2026-05-01T00:00:00Z")
	writeBundleFixture(t, artifactsBase, "binance", "PR01JUN26", "2026-06-01T00:00:00Z")
	writeBundleFixture(t, artifactsBase, "binance", "PR01JUL26", "2026-07-01T00:00:00Z")

	for _, id := range []string{"PR01APR26", "PR01MAY26", "PR01JUN26"} {
		writeAssessment(t, reportsDir, "binance", id, map[string]string{
			"PR01APR26": "2026-04-01T00:00:00Z",
			"PR01MAY26": "2026-05-01T00:00:00Z",
			"PR01JUN26": "2026-06-01T00:00:00Z",
		}[id])
	}

	// PR01JUL26 is NOT in the reports dir: this is exactly the state Write sees
	// while rendering its own page.
	if err := Write(context.Background(), Options{
		Exchange:   "binance",
		SnapshotID: "PR01JUL26",
		Artifacts:  artifactsBase,
		ReportsDir: reportsDir,
	}); err != nil {
		t.Fatalf("Write: %v", err)
	}

	raw, err := os.ReadFile(filepath.Join(reportsDir, "binance", "PR01JUL26-transparency-report.md"))
	if err != nil {
		t.Fatal(err)
	}
	fields := frequencySection(t, string(raw))

	if got := fields["Latest snapshot (evaluation set)"]; got != "2026-07-01T00:00:00Z" {
		t.Errorf("latest: got %q want 2026-07-01T00:00:00Z (own snapshot missing from Frequency table)", got)
	}
	if got := fields["Previous snapshot"]; got != "2026-06-01T00:00:00Z" {
		t.Errorf("previous: got %q want 2026-06-01T00:00:00Z", got)
	}
	if got := fields["History available"]; got != "4 snapshot(s) in public evaluation set" {
		t.Errorf("history: got %q want 4 snapshots", got)
	}
	if got := fields["Observed cadence"]; got != "~monthly" {
		t.Errorf("cadence: got %q want ~monthly", got)
	}

	// Regenerating must not change the numbers: the page is a fixed point.
	before := fields
	if err := Write(context.Background(), Options{
		Exchange:   "binance",
		SnapshotID: "PR01JUL26",
		Artifacts:  artifactsBase,
		ReportsDir: reportsDir,
	}); err != nil {
		t.Fatalf("Write (rerun): %v", err)
	}
	raw, err = os.ReadFile(filepath.Join(reportsDir, "binance", "PR01JUL26-transparency-report.md"))
	if err != nil {
		t.Fatal(err)
	}
	after := frequencySection(t, string(raw))
	for k, want := range before {
		if after[k] != want {
			t.Errorf("rerun changed %q: %q -> %q", k, want, after[k])
		}
	}
}
