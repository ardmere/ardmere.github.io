package reportgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

// TestCadenceIgnoresLegacySnapshot is the fix for the reported bug: one ancient
// snapshot next to a monthly series must not drag the observed cadence.
func TestCadenceIgnoresLegacySnapshot(t *testing.T) {
	reportsDir := t.TempDir()
	writeAssessment(t, reportsDir, "htx", "20230910", "2023-09-10T18:10:33Z")
	writeAssessment(t, reportsDir, "htx", "20260801", "2026-08-01T00:00:00Z")
	writeAssessment(t, reportsDir, "htx", "20260901", "2026-09-01T00:00:00Z")
	writeAssessment(t, reportsDir, "htx", "20261001", "2026-10-01T00:00:00Z")

	freq := collectFrequency("htx", "20261101", "2026-11-01T00:00:00Z", reportsDir)

	if freq.ObservedCadence != "~monthly" {
		t.Fatalf("cadence: got %q want ~monthly", freq.ObservedCadence)
	}
	// The discarded span must be disclosed, not silently dropped.
	for _, want := range []string{"20230910", "20260801", "1055 days"} {
		if !strings.Contains(freq.Notes, want) {
			t.Errorf("notes do not disclose %q: %s", want, freq.Notes)
		}
	}
}

// TestCadenceRetainsLatestSilence keeps the honest failure mode: a venue that
// has stopped publishing must not read as healthy just because older snapshots
// were frequent.
func TestCadenceRetainsLatestSilence(t *testing.T) {
	reportsDir := t.TempDir()
	writeAssessment(t, reportsDir, "dormant", "20250101", "2025-01-01T00:00:00Z")
	writeAssessment(t, reportsDir, "dormant", "20250201", "2025-02-01T00:00:00Z")
	writeAssessment(t, reportsDir, "dormant", "20250301", "2025-03-01T00:00:00Z")

	// The newest gap (2026-11-01 minus 2025-03-01 = 610 days) is the silence.
	freq := collectFrequency("dormant", "20261101", "2026-11-01T00:00:00Z", reportsDir)

	if freq.ObservedCadence == "~monthly" {
		t.Fatalf("dormant venue read as healthy: %q", freq.ObservedCadence)
	}
	if !strings.Contains(freq.ObservedCadence, "between snapshots") {
		t.Fatalf("dormant venue cadence: %q", freq.ObservedCadence)
	}
}

// TestCadenceHTX61101Regression pins the reported HTX page inputs.
func TestCadenceHTX61101Regression(t *testing.T) {
	reportsDir := filepath.Join("..", "..", "docs", "reports")
	freq := collectFrequency("htx", "20261101", "2026-11-01T00:00:00Z", reportsDir)

	if freq.LatestSnapshot != "2026-11-01T00:00:00Z" {
		t.Errorf("latest: %q", freq.LatestSnapshot)
	}
	if freq.PreviousSnapshot != "2026-10-01T00:00:00Z" {
		t.Errorf("previous: %q", freq.PreviousSnapshot)
	}
	if freq.ObservedCadence != "~monthly" {
		t.Errorf("cadence: got %q want ~monthly", freq.ObservedCadence)
	}
	if freq.HistoryAvailable != "5 snapshot(s) in public evaluation set" {
		t.Errorf("history: %q", freq.HistoryAvailable)
	}
}

func cadenceRec(id string, y int, m, d int) snapshotRecord {
	return snapshotRecord{ID: id, SnapshotTime: time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)}
}

// TestEstimateCadence_TrimAndDormancyTable pins the estimator on the shapes the
// rule has to distinguish, including the two-gap case where an averaged median
// would let the legacy gap set the reference interval.
func TestEstimateCadence_TrimAndDormancyTable(t *testing.T) {
	cases := []struct {
		name     string
		records  []snapshotRecord
		want     string
		wantKept int
		wantExcl int
	}{
		{
			name:     "legacy plus two monthly gaps",
			records:  []snapshotRecord{cadenceRec("legacy", 2023, 9, 10), cadenceRec("aug", 2026, 8, 1), cadenceRec("sep", 2026, 9, 1)},
			want:     "~monthly",
			wantKept: 1,
			wantExcl: 1,
		},
		{
			name: "monthly series with one legacy file at the front",
			records: []snapshotRecord{
				cadenceRec("legacy", 2023, 9, 10), cadenceRec("aug", 2026, 8, 1),
				cadenceRec("sep", 2026, 9, 1), cadenceRec("oct", 2026, 10, 1),
				cadenceRec("nov", 2026, 11, 1),
			},
			want:     "~monthly",
			wantKept: 3,
			wantExcl: 1,
		},
		{
			name:    "dormant after a quarterly run keeps the silence",
			records: []snapshotRecord{cadenceRec("q1", 2025, 1, 1), cadenceRec("q2", 2025, 4, 1), cadenceRec("q3", 2025, 7, 1), cadenceRec("q4", 2025, 10, 1), cadenceRec("silence", 2026, 11, 1)},
			// Gaps 90/91/92/396 days. The 396-day newest gap exceeds the trim
			// limit but no newer gap re-establishes a rhythm, so it stays in the
			// mean and the venue reads slow instead of ~quarterly.
			want:     "~167 days between snapshots",
			wantKept: 4,
			wantExcl: 0,
		},
		{
			name:    "coarse but steady venue is not trimmed",
			records: []snapshotRecord{cadenceRec("a", 2025, 1, 1), cadenceRec("b", 2025, 4, 1), cadenceRec("c", 2025, 7, 1), cadenceRec("d", 2025, 10, 1)},
			want:     "~91 days between snapshots",
			wantKept: 3,
			wantExcl: 0,
		},
		{
			name:    "only two snapshots ever",
			records: []snapshotRecord{cadenceRec("a", 2024, 1, 1), cadenceRec("b", 2026, 1, 1)},
			want:     "~731 days between snapshots",
			wantKept: 1,
			wantExcl: 0,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			est := estimateCadence(tc.records)
			if est.Label != tc.want {
				t.Errorf("label: got %q want %q", est.Label, tc.want)
			}
			if est.Kept != tc.wantKept || len(est.Excluded) != tc.wantExcl {
				t.Errorf("kept=%d excluded=%d, want kept=%d excluded=%d", est.Kept, len(est.Excluded), tc.wantKept, tc.wantExcl)
			}
		})
	}
}
