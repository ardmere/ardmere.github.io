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
	LatestSnapshot       string
	PreviousSnapshot     string
	ObservedCadence      string
	HistoryAvailable     string
	EventTriggered       string
	DailyAnchor          string
	StageImpact          string
	Notes                string
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

	records := loadExchangeSnapshots(exchangeID, reportsDir)
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

	var idx = -1
	for i, r := range records {
		if r.ID == snapshotID {
			idx = i
			break
		}
	}
	if idx >= 0 && idx+1 < len(records) {
		info.PreviousSnapshot = records[idx+1].SnapshotTime.UTC().Format(time.RFC3339)
	}

	info.ObservedCadence = describeCadence(records)
	info.Notes = frequencyNotes(exchangeID, records, idx)
	return info
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

func describeCadence(records []snapshotRecord) string {
	if len(records) < 2 {
		return "UNVERIFIABLE"
	}
	sorted := append([]snapshotRecord(nil), records...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].SnapshotTime.Before(sorted[j].SnapshotTime)
	})

	var gaps []time.Duration
	for i := 1; i < len(sorted); i++ {
		gaps = append(gaps, sorted[i].SnapshotTime.Sub(sorted[i-1].SnapshotTime))
	}
	avg := averageDuration(gaps)
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

func frequencyNotes(exchangeID string, records []snapshotRecord, idx int) string {
	if len(records) == 1 {
		return fmt.Sprintf("Only one %s snapshot is in the ardmere public evaluation set.", exchangeID)
	}
	if idx == 0 {
		return "This is the newest snapshot in the ardmere public evaluation set for " + exchangeID + "."
	}
	if idx > 0 {
		return fmt.Sprintf("Older snapshot %s is also in the public evaluation set.", records[0].ID)
	}
	return ""
}
