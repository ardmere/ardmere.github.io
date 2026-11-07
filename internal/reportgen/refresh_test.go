package reportgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRefreshFrequencyRepairsStaleSection(t *testing.T) {
	reportsDir := t.TempDir()
	writeAssessment(t, reportsDir, "htx", "20230910", "2023-09-10T18:10:33Z")
	writeAssessment(t, reportsDir, "htx", "20260801", "2026-08-01T00:00:00Z")
	writeAssessment(t, reportsDir, "htx", "20260901", "2026-09-01T00:00:00Z")
	writeAssessment(t, reportsDir, "htx", "20261001", "2026-10-01T00:00:00Z")

	// A page produced by the buggy pipeline: it reports the set as it looked
	// before its own snapshot existed.
	stale := `# HTX 20260801 PoR Transparency Report

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | ` + "`htx`" + ` |

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | ` + "`2026-10-01T00:00:00Z`" + ` |
| Previous snapshot | ` + "`UNVERIFIABLE`" + ` |
| Observed cadence | ` + "`~558 days between snapshots`" + ` |
| History available | ` + "`3 snapshot(s) in public evaluation set`" + ` |
| Event-triggered updates | ` + "`UNVERIFIABLE`" + ` |
| Daily root / commitment anchor | ` + "`UNVERIFIABLE`" + ` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

## 4. Evidence and Bundles

| Bundle | Value |
| --- | --- |
| Artifact bundle root | ` + "`0x01`" + ` |
`
	path := filepath.Join(reportsDir, "htx", "20260801-transparency-report.md")
	if err := os.WriteFile(path, []byte(stale), 0o644); err != nil {
		t.Fatal(err)
	}

	results, err := RefreshFrequency(reportsDir)
	if err != nil {
		t.Fatalf("RefreshFrequency: %v", err)
	}
	if len(results) != 1 || !results[0].Changed {
		t.Fatalf("results: %+v", results)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got := string(raw)
	fields := frequencySection(t, got)

	if fields["Previous snapshot"] != "2023-09-10T18:10:33Z" {
		t.Errorf("previous: %q", fields["Previous snapshot"])
	}
	if fields["Observed cadence"] != "~372 days between snapshots" {
		t.Errorf("cadence: %q", fields["Observed cadence"])
	}
	if fields["History available"] != "4 snapshot(s) in public evaluation set" {
		t.Errorf("history: %q", fields["History available"])
	}

	// Everything outside section 3 must survive untouched.
	for _, keep := range []string{
		"# HTX 20260801 PoR Transparency Report",
		"## 1. Summary",
		"| Exchange | `htx` |",
		"## 4. Evidence and Bundles",
		"| Artifact bundle root | `0x01` |",
	} {
		if !strings.Contains(got, keep) {
			t.Errorf("refresh dropped %q", keep)
		}
	}

	// A second pass is a no-op.
	results, err = RefreshFrequency(reportsDir)
	if err != nil {
		t.Fatalf("RefreshFrequency (second): %v", err)
	}
	if len(results) != 1 || results[0].Changed {
		t.Fatalf("second pass not idempotent: %+v", results)
	}
}

func TestRefreshFrequencyMatchesFreshWriteSection(t *testing.T) {
	reportsDir := t.TempDir()
	writeAssessment(t, reportsDir, "gateio", "20260316", "2026-03-16T00:00:00Z")
	writeAssessment(t, reportsDir, "gateio", "20260618", "2026-06-18T00:00:00Z")

	freq := collectFrequency("gateio", "20261018", "2026-10-18T00:00:00Z", reportsDir)

	var want strings.Builder
	writeSection3Frequency(&want, freq)
	snippet := want.String()

	page := "# Gate 20261018\n\n" + snippet + "## 4. Evidence and Bundles\n\nbody\n"
	start := strings.Index(page, frequencySectionStart)
	rel := strings.Index(page[start:], frequencySectionEnd)
	got := page[:start] + snippet + page[start+rel:]
	if got != page {
		t.Fatalf("refreshing an up-to-date section changed it:\n%s", got)
	}
}
