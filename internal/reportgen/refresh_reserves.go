package reportgen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ardmere/ardmere/internal/pipeline"
)

// RefreshReserveDisplayResult is one report whose reserve section was examined.
type RefreshReserveDisplayResult struct {
	Exchange   string
	SnapshotID string
	Changed    bool
}

// RefreshReserveDisplay rewrites the self-reported reserve markdown section from
// assessment JSON on disk (null → Unavailable). It does not re-price amounts or
// call CoinGecko — use PatchReserveFigures when totals themselves must change.
func RefreshReserveDisplay(reportsDir string) ([]RefreshReserveDisplayResult, error) {
	matches, err := filepath.Glob(filepath.Join(reportsDir, "*", "*-transparency-report.md"))
	if err != nil {
		return nil, err
	}
	sort.Strings(matches)

	var out []RefreshReserveDisplayResult
	for _, path := range matches {
		exchange := filepath.Base(filepath.Dir(path))
		snapshotID := strings.TrimSuffix(filepath.Base(path), "-transparency-report.md")
		assessmentPath := filepath.Join(reportsDir, exchange, snapshotID+"-assessment.json")
		fig, err := figuresFromAssessmentFile(assessmentPath, exchange, snapshotID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", assessmentPath, err)
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		updated, err := spliceReportMarkdown(string(raw), fig)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		changed := updated != string(raw)
		if changed {
			if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
				return nil, err
			}
		}
		out = append(out, RefreshReserveDisplayResult{
			Exchange:   exchange,
			SnapshotID: snapshotID,
			Changed:    changed,
		})
	}
	return out, nil
}

func figuresFromAssessmentFile(path, exchange, snapshotID string) (pipeline.ReserveFigures, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return pipeline.ReserveFigures{}, err
	}
	var doc struct {
		TotalLiabilityUsd     *float64                        `json:"totalLiabilityUsd"`
		TotalReserveUsd       *float64                        `json:"totalReserveUsd"`
		UnpricedCoins         []string                        `json:"unpricedCoins"`
		MajorCoinReserveRates []pipeline.MajorCoinReserveRate `json:"majorCoinReserveRates"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return pipeline.ReserveFigures{}, err
	}
	fig := pipeline.ReserveFigures{
		TotalLiabilityUsd:     doc.TotalLiabilityUsd,
		TotalReserveUsd:       doc.TotalReserveUsd,
		UnpricedCoins:         doc.UnpricedCoins,
		MajorCoinReserveRates: doc.MajorCoinReserveRates,
	}
	if fig.UnpricedCoins == nil {
		fig.UnpricedCoins = []string{}
	}
	if fig.MajorCoinReserveRates == nil {
		fig.MajorCoinReserveRates = []pipeline.MajorCoinReserveRate{}
	}
	if fig.TotalLiabilityUsd != nil || fig.TotalReserveUsd != nil || len(fig.UnpricedCoins) > 0 {
		fig.AmountsPublished = true
	}
	if fig.TotalLiabilityUsd == nil && fig.TotalReserveUsd == nil && len(fig.UnpricedCoins) == 0 {
		if _, err := FindSummaryPath("fixtures", exchange, snapshotID); err != nil {
			fig.SourceMissing = true
		} else {
			fig.AmountsPublished = false
		}
	}
	return fig, nil
}
