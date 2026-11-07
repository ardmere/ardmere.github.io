package reportgen

import (
	"os"
	"path/filepath"
	"testing"
)

// writeAssessment puts a minimal assessment JSON into reports/<exchange>/ so the
// frequency scan can see it.
func writeAssessment(t *testing.T, reportsDir, exchangeID, snapshotID, snapshotTime string) {
	t.Helper()
	dir := filepath.Join(reportsDir, exchangeID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := `{
  "schema": "ardmere/exchange-assessment@1",
  "exchange": "` + exchangeID + `",
  "snapshot": {"id": "` + snapshotID + `", "snapshotTime": "` + snapshotTime + `"}
}
`
	if err := os.WriteFile(filepath.Join(dir, snapshotID+"-assessment.json"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestCollectFrequency_IncludesSnapshotBeingWritten(t *testing.T) {
	reportsDir := t.TempDir()
	writeAssessment(t, reportsDir, "binance", "PR01APR26", "2026-04-01T00:00:00Z")
	writeAssessment(t, reportsDir, "binance", "PR01MAY26", "2026-05-01T00:00:00Z")
	writeAssessment(t, reportsDir, "binance", "PR01JUN26", "2026-06-01T00:00:00Z")
	// PR01JUL26 is deliberately absent from disk: this is the state Write() sees
	// when it renders the report for a snapshot it has not written yet.
	freq := collectFrequency("binance", "PR01JUL26", "2026-07-01T00:00:00Z", reportsDir)

	if freq.LatestSnapshot != "2026-07-01T00:00:00Z" {
		t.Fatalf("latest: got %q want the snapshot being written", freq.LatestSnapshot)
	}
	if freq.PreviousSnapshot != "2026-06-01T00:00:00Z" {
		t.Fatalf("previous: got %q want 2026-06-01T00:00:00Z", freq.PreviousSnapshot)
	}
	if freq.ObservedCadence != "~monthly" {
		t.Fatalf("cadence: got %q want ~monthly", freq.ObservedCadence)
	}
	if freq.HistoryAvailable != "4 snapshot(s) in public evaluation set" {
		t.Fatalf("history: got %q want 4 snapshots", freq.HistoryAvailable)
	}
}

func TestCollectFrequency_SingleSnapshotCountsItself(t *testing.T) {
	reportsDir := t.TempDir()
	freq := collectFrequency("htx", "20261101", "2026-11-01T00:00:00Z", reportsDir)

	if freq.LatestSnapshot != "2026-11-01T00:00:00Z" {
		t.Fatalf("latest: %q", freq.LatestSnapshot)
	}
	if freq.HistoryAvailable != "1 snapshot(s) in public evaluation set" {
		t.Fatalf("history: %q", freq.HistoryAvailable)
	}
	if freq.ObservedCadence != "UNVERIFIABLE" {
		t.Fatalf("cadence: got %q want UNVERIFIABLE for one snapshot", freq.ObservedCadence)
	}
}

func TestCollectFrequency_DoesNotDoubleCountExistingAssessment(t *testing.T) {
	reportsDir := t.TempDir()
	writeAssessment(t, reportsDir, "gateio", "20260316", "2026-03-16T00:00:00Z")
	writeAssessment(t, reportsDir, "gateio", "20260618", "2026-06-18T00:00:00Z")

	// Re-running the report for a snapshot already on disk must not count it twice.
	freq := collectFrequency("gateio", "20260618", "2026-06-18T00:00:00Z", reportsDir)

	if freq.HistoryAvailable != "2 snapshot(s) in public evaluation set" {
		t.Fatalf("history: %q", freq.HistoryAvailable)
	}
	if freq.PreviousSnapshot != "2026-03-16T00:00:00Z" {
		t.Fatalf("previous: %q", freq.PreviousSnapshot)
	}
}

func TestCollectFrequency_CadenceIsSaneForMonthlySeries(t *testing.T) {
	reportsDir := t.TempDir()
	writeAssessment(t, reportsDir, "binance", "PR01APR26", "2026-04-01T00:00:00Z")
	writeAssessment(t, reportsDir, "binance", "PR01MAY26", "2026-05-01T00:00:00Z")

	freq := collectFrequency("binance", "PR01JUN26", "2026-06-01T00:00:00Z", reportsDir)
	if freq.ObservedCadence != "~monthly" {
		t.Fatalf("cadence: got %q want ~monthly", freq.ObservedCadence)
	}
}

func TestCollectFrequency_OlderSnapshotGetsLaterSiblingAsLatest(t *testing.T) {
	reportsDir := t.TempDir()
	writeAssessment(t, reportsDir, "htx", "20230910", "2023-09-10T18:10:33Z")
	writeAssessment(t, reportsDir, "htx", "20260901", "2026-09-01T00:00:00Z")

	freq := collectFrequency("htx", "20260801", "2026-08-01T00:00:00Z", reportsDir)

	if freq.LatestSnapshot != "2026-09-01T00:00:00Z" {
		t.Fatalf("latest: %q", freq.LatestSnapshot)
	}
	if freq.PreviousSnapshot != "2023-09-10T18:10:33Z" {
		t.Fatalf("previous: %q", freq.PreviousSnapshot)
	}
	if freq.HistoryAvailable != "3 snapshot(s) in public evaluation set" {
		t.Fatalf("history: %q", freq.HistoryAvailable)
	}
}
