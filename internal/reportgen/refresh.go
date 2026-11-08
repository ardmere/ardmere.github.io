package reportgen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	frequencySectionStart = "## 3. Frequency and Freshness"
	frequencySectionEnd   = "## 4. Evidence and Bundles"
)

// RefreshFrequencyResult reports one report page whose markdown was examined.
type RefreshFrequencyResult struct {
	Exchange   string
	SnapshotID string
	Changed    bool
}

// RefreshFrequency recomputes Section 3 of every existing transparency report
// markdown file from the assessment JSONs on disk and rewrites the section when
// its numbers differ.
//
// This is the repair path for pages written by the pipeline version that
// rendered the report before its own assessment JSON existed: it produces the
// same Section 3 that a fresh Write() would, without re-running fetch/verify.
// HTML pages are derived from the markdown by scripts/build-public-html.sh.
func RefreshFrequency(reportsDir string) ([]RefreshFrequencyResult, error) {
	matches, err := filepath.Glob(filepath.Join(reportsDir, "*", "*-transparency-report.md"))
	if err != nil {
		return nil, err
	}
	sort.Strings(matches)

	results := make([]RefreshFrequencyResult, 0, len(matches))
	for _, path := range matches {
		exchangeID := filepath.Base(filepath.Dir(path))
		snapshotID := strings.TrimSuffix(filepath.Base(path), "-transparency-report.md")

		snapshotTime, err := assessmentSnapshotTime(reportsDir, exchangeID, snapshotID)
		if err != nil {
			return nil, err
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		freq := collectFrequency(exchangeID, snapshotID, snapshotTime, reportsDir)

		var section strings.Builder
		writeSection3Frequency(&section, freq)

		page := string(raw)
		start := strings.Index(page, frequencySectionStart)
		if start < 0 {
			return nil, fmt.Errorf("%s: section 3 heading not found", path)
		}
		rel := strings.Index(page[start:], frequencySectionEnd)
		if rel < 0 {
			return nil, fmt.Errorf("%s: section 4 heading not found after section 3", path)
		}
		updated := page[:start] + section.String() + page[start+rel:]

		changed := updated != page
		if changed {
			if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
				return nil, err
			}
		}
		results = append(results, RefreshFrequencyResult{
			Exchange:   exchangeID,
			SnapshotID: snapshotID,
			Changed:    changed,
		})
	}
	return results, nil
}

func assessmentSnapshotTime(reportsDir, exchangeID, snapshotID string) (string, error) {
	path := filepath.Join(reportsDir, exchangeID, snapshotID+"-assessment.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read assessment for %s %s: %w", exchangeID, snapshotID, err)
	}
	var doc struct {
		Snapshot struct {
			SnapshotTime string `json:"snapshotTime"`
		} `json:"snapshot"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return "", fmt.Errorf("parse assessment for %s %s: %w", exchangeID, snapshotID, err)
	}
	return doc.Snapshot.SnapshotTime, nil
}
