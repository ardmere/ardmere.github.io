package reportgen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type snapshotRecord struct {
	ID           string
	SnapshotTime time.Time
}

type frequencyInfo struct {
	LatestSnapshot   string
	PreviousSnapshot string
	ObservedCadence  string
	HistoryAvailable string
	EventTriggered   string
	DailyAnchor      string
	StageImpact      string
	Notes            string
}

func collectFrequency(exchangeID, snapshotID, snapshotTimeRFC, reportsDir string) frequencyInfo {
	info := frequencyInfo{
		LatestSnapshot:   snapshotTimeRFC,
		PreviousSnapshot: "UNVERIFIABLE",
		EventTriggered:   "UNVERIFIABLE",
		DailyAnchor:      "UNVERIFIABLE",
		StageImpact:      "Monthly or slower cadence; does not meet Stage 2 frequency expectations.",
	}
	if strings.TrimSpace(snapshotTimeRFC) != "" {
		if t, err := time.Parse(time.RFC3339, snapshotTimeRFC); err == nil {
			info.LatestSnapshot = t.UTC().Format(time.RFC3339)
		}
	}

	records := includeCurrentSnapshot(loadExchangeSnapshots(exchangeID, reportsDir), snapshotID, snapshotTimeRFC)
	if len(records) == 0 {
		info.ObservedCadence = "UNVERIFIABLE"
		info.HistoryAvailable = "UNVERIFIABLE"
		info.Notes = "No sibling assessments found in the ardmere public evaluation set."
		return info
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].SnapshotTime.After(records[j].SnapshotTime)
	})

	info.LatestSnapshot = records[0].SnapshotTime.UTC().Format(time.RFC3339)
	info.HistoryAvailable = fmt.Sprintf("%d snapshot(s) in public evaluation set", len(records))

	idx := -1
	for i, r := range records {
		if r.ID == snapshotID {
			idx = i
			break
		}
	}
	if idx >= 0 && idx+1 < len(records) {
		info.PreviousSnapshot = records[idx+1].SnapshotTime.UTC().Format(time.RFC3339)
	}

	est := estimateCadence(records)
	info.ObservedCadence = est.Label
	info.Notes = joinSentences(frequencyNotes(exchangeID, records, idx), cadenceExclusionNote(est))
	return info
}

func joinSentences(parts ...string) string {
	var kept []string
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			kept = append(kept, p)
		}
	}
	return strings.Join(kept, " ")
}

// cadenceExclusionNote discloses the gaps the cadence average left out, so a
// reader dividing the history by hand can reproduce the number and see the
// legacy span that was excluded rather than mistaking it for a missing period.
func cadenceExclusionNote(est cadenceEstimate) string {
	if len(est.Excluded) == 0 {
		return ""
	}
	parts := make([]string, 0, len(est.Excluded))
	for _, g := range est.Excluded {
		parts = append(parts, fmt.Sprintf("`%s` → `%s` (%d days)", g.From.ID, g.To.ID, int(g.Duration.Hours()/24)))
	}
	return fmt.Sprintf(
		"Observed cadence averages %d of %d adjacent gaps; excluded as an older span superseded by the current rhythm: %s. "+
			"Retained gaps are each at most %dx the %d-day median gap.",
		est.Kept, est.Gaps, strings.Join(parts, ", "), cadenceOutlierFactor, est.medianDays())
}

// includeCurrentSnapshot injects the snapshot being written into the record set
// scanned from disk. Write() renders the report before the assessment JSON
// exists, so without this the Frequency section describes the previous run and
// omits the page's own snapshot -- the source of the one-run-stale numbers.
func includeCurrentSnapshot(records []snapshotRecord, snapshotID, snapshotTimeRFC string) []snapshotRecord {
	if strings.TrimSpace(snapshotID) == "" {
		return records
	}
	t, err := time.Parse(time.RFC3339, snapshotTimeRFC)
	if err != nil {
		return records
	}
	current := snapshotRecord{ID: snapshotID, SnapshotTime: t.UTC()}

	out := make([]snapshotRecord, 0, len(records)+1)
	replaced := false
	for _, r := range records {
		if r.ID == snapshotID {
			if !replaced {
				out = append(out, current)
				replaced = true
			}
			continue
		}
		out = append(out, r)
	}
	if !replaced {
		out = append(out, current)
	}
	return out
}

func loadExchangeSnapshots(exchangeID, reportsDir string) []snapshotRecord {
	dir := filepath.Join(reportsDir, exchangeID)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []snapshotRecord
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), "-assessment.json") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		var doc struct {
			Snapshot struct {
				ID           string `json:"id"`
				SnapshotTime string `json:"snapshotTime"`
			} `json:"snapshot"`
		}
		if err := json.Unmarshal(raw, &doc); err != nil || doc.Snapshot.ID == "" {
			continue
		}
		t, err := time.Parse(time.RFC3339, doc.Snapshot.SnapshotTime)
		if err != nil {
			continue
		}
		out = append(out, snapshotRecord{
			ID:           doc.Snapshot.ID,
			SnapshotTime: t,
		})
	}
	return out
}

// cadenceOutlierFactor is how many times the median adjacent gap a single gap
// may span before it is treated as evidence of an older publishing regime
// rather than of the venue's current cadence.
const cadenceOutlierFactor = 3

type snapshotGap struct {
	From     snapshotRecord
	To       snapshotRecord
	Duration time.Duration
}

// cadenceEstimate is the observed-cadence number plus the gaps it discarded, so
// Section 3 can disclose what was left out of the arithmetic.
type cadenceEstimate struct {
	Label    string
	Mean     time.Duration
	Median   time.Duration
	Gaps     int
	Kept     int
	Excluded []snapshotGap
}

func (e cadenceEstimate) medianDays() int {
	return int(e.Median.Hours() / 24)
}

// estimateCadence reports how often a venue publishes today.
//
// Rule: build every adjacent gap in the public history, take their median as
// the reference interval, and average only the gaps that are at most
// cadenceOutlierFactor times that median -- and only when a newer gap already
// re-establishes the faster rhythm. HTX's 2023-09-10 file next to a monthly
// 2026 series produces one gap 34x the median, so it is retired and the number
// follows the current regime. A venue whose newest behaviour is a long silence
// keeps that silence, and a venue that is slow throughout has no gap far above
// its own median, so both still read slow.
func estimateCadence(records []snapshotRecord) cadenceEstimate {
	if len(records) < 2 {
		return cadenceEstimate{Label: "UNVERIFIABLE"}
	}
	sorted := append([]snapshotRecord(nil), records...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].SnapshotTime.Before(sorted[j].SnapshotTime)
	})

	gaps := make([]snapshotGap, 0, len(sorted)-1)
	durations := make([]time.Duration, 0, len(sorted)-1)
	for i := 1; i < len(sorted); i++ {
		g := snapshotGap{
			From:     sorted[i-1],
			To:       sorted[i],
			Duration: sorted[i].SnapshotTime.Sub(sorted[i-1].SnapshotTime),
		}
		gaps = append(gaps, g)
		durations = append(durations, g.Duration)
	}

	median := lowerMedianDuration(durations)
	limit := time.Duration(float64(median) * cadenceOutlierFactor)
	withinLimit := make([]bool, len(gaps))
	for i, g := range gaps {
		withinLimit[i] = g.Duration <= limit
	}
	// newerIsCurrent[i] means a gap after i already shows the current rhythm,
	// which is what makes gap i safe to retire as an older regime.
	newerIsCurrent := make([]bool, len(gaps))
	sawCurrent := false
	for i := len(gaps) - 1; i >= 0; i-- {
		newerIsCurrent[i] = sawCurrent
		if withinLimit[i] {
			sawCurrent = true
		}
	}

	est := cadenceEstimate{Gaps: len(gaps), Median: median}
	kept := make([]time.Duration, 0, len(gaps))
	for i, g := range gaps {
		if !withinLimit[i] && newerIsCurrent[i] {
			est.Excluded = append(est.Excluded, g)
			continue
		}
		kept = append(kept, g.Duration)
	}
	est.Kept = len(kept)
	est.Mean = averageDuration(kept)
	est.Label = cadenceLabel(est.Mean)
	return est
}

func cadenceLabel(avg time.Duration) string {
	switch {
	case avg >= 25*24*time.Hour && avg <= 35*24*time.Hour:
		return "~monthly"
	case avg >= 6*24*time.Hour && avg <= 8*24*time.Hour:
		return "~weekly"
	case avg < 24*time.Hour:
		return "sub-daily"
	default:
		return fmt.Sprintf("~%.0f days between snapshots", avg.Hours()/24)
	}
}

func averageDuration(ds []time.Duration) time.Duration {
	if len(ds) == 0 {
		return 0
	}
	var sum time.Duration
	for _, d := range ds {
		sum += d
	}
	return sum / time.Duration(len(ds))
}

// lowerMedianDuration returns the lower median when the sample size is even.
// Averaging the two middle gaps would let an old, very long gap pull the
// reference interval up, which is the failure this estimator exists to avoid:
// with gaps [30d, 1055d] the average of the middle pair is 543 days, and only
// the 30-day gap survives the trim. The lower median keeps the reference at the
// faster of the two middle observations, so a legacy snapshot never sets it.
func lowerMedianDuration(ds []time.Duration) time.Duration {
	if len(ds) == 0 {
		return 0
	}
	sorted := append([]time.Duration(nil), ds...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	return sorted[(len(sorted)-1)/2]
}

func frequencyNotes(exchangeID string, records []snapshotRecord, idx int) string {
	if len(records) == 1 {
		return fmt.Sprintf("Only one %s snapshot is in the ardmere public evaluation set.", exchangeID)
	}
	if idx == 0 {
		return "This is the newest snapshot in the ardmere public evaluation set for " + exchangeID + "."
	}
	if idx > 0 {
		return fmt.Sprintf("Newer snapshot %s is also in the public evaluation set.", records[idx-1].ID)
	}
	return ""
}
