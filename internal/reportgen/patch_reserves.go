package reportgen

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ardmere/ardmere/internal/coingecko"
	"github.com/ardmere/ardmere/internal/exchanges/binance"
	"github.com/ardmere/ardmere/internal/exchanges/binance/bapi"
	"github.com/ardmere/ardmere/internal/exchanges/okx"
	"github.com/ardmere/ardmere/internal/exchanges/okx/okxapi"
	"github.com/ardmere/ardmere/internal/pipeline"
	"github.com/ardmere/ardmere/internal/por"
)

// PatchResult is one public assessment after the reserve-figure splice.
type PatchResult struct {
	Exchange          string
	SnapshotID        string
	Source            string
	TotalLiabilityUsd *float64
	TotalReserveUsd   *float64
	MajorCoins        int
	Unpriced          int
	AmountsPublished  bool
	AssessmentChanged bool
	ReportChanged     bool
}

// PatchOptions selects the public report tree and where summaries are read.
type PatchOptions struct {
	ReportsDir   string
	FixturesDir  string
	ArtifactsDir string
	// Only limits the walk to exchange/snapshotID pairs when non-empty.
	Only map[string]struct{}
	// AllowNetwork fetches the current Binance or OKX summary when the audit id matches.
	AllowNetwork bool
}

// PatchReserveFigures writes totalLiabilityUsd, totalReserveUsd, unpricedCoins,
// and majorCoinReserveRates into each assessment and inserts the report section.
// It does not re-run verifiers. Summaries come from fixtures, cached artifacts,
// or (when AllowNetwork) the current Binance/OKX page when that page is this snapshot.
func PatchReserveFigures(ctx context.Context, opt PatchOptions) ([]PatchResult, error) {
	if opt.ReportsDir == "" {
		opt.ReportsDir = "docs/reports"
	}
	if opt.FixturesDir == "" {
		opt.FixturesDir = "fixtures"
	}
	matches, err := filepath.Glob(filepath.Join(opt.ReportsDir, "*", "*-assessment.json"))
	if err != nil {
		return nil, err
	}
	prices := coingecko.New()
	var out []PatchResult
	for _, path := range matches {
		exchange := filepath.Base(filepath.Dir(path))
		snapshotID := strings.TrimSuffix(filepath.Base(path), "-assessment.json")
		if len(opt.Only) > 0 {
			if _, ok := opt.Only[exchange+"/"+snapshotID]; !ok {
				continue
			}
		}
		fig, source := figuresForPatch(ctx, opt, prices, exchange, snapshotID)
		changedA, err := spliceAssessmentFile(path, fig)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		mdPath := filepath.Join(opt.ReportsDir, exchange, snapshotID+"-transparency-report.md")
		changedM := false
		if raw, err := os.ReadFile(mdPath); err == nil {
			updated, err := spliceReportMarkdown(string(raw), fig)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", mdPath, err)
			}
			changedM = updated != string(raw)
			if changedM {
				if err := os.WriteFile(mdPath, []byte(updated), 0o644); err != nil {
					return nil, err
				}
			}
		}
		out = append(out, PatchResult{
			Exchange:          exchange,
			SnapshotID:        snapshotID,
			Source:            source,
			TotalLiabilityUsd: fig.TotalLiabilityUsd,
			TotalReserveUsd:   fig.TotalReserveUsd,
			MajorCoins:        len(fig.MajorCoinReserveRates),
			Unpriced:          len(fig.UnpricedCoins),
			AmountsPublished:  fig.AmountsPublished,
			AssessmentChanged: changedA,
			ReportChanged:     changedM,
		})
	}
	return out, nil
}

func figuresForPatch(ctx context.Context, opt PatchOptions, prices pipeline.PriceOracle, exchange, snapshotID string) (pipeline.ReserveFigures, string) {
	if opt.ArtifactsDir != "" {
		dir := filepath.Join(opt.ArtifactsDir, exchange, snapshotID)
		if loaded, err := pipeline.LoadFromBundle(dir, snapshotID, exchange); err == nil {
			return pipeline.ComputeReserveFigures(ctx, loaded.Snapshot, prices), "artifacts"
		}
	}
	if path, err := FindSummaryPath(opt.FixturesDir, exchange, snapshotID); err == nil {
		snap, err := LoadSummaryFile(exchange, path)
		if err == nil {
			if snap.SnapshotTime.IsZero() {
				snap.SnapshotTime = assessmentSnapshotTimeValue(filepath.Join(opt.ReportsDir, exchange, snapshotID+"-assessment.json"))
			}
			return pipeline.ComputeReserveFigures(ctx, snap, prices), "fixture:" + path
		}
	}
	if opt.AllowNetwork {
		nctx, cancel := context.WithTimeout(ctx, 20*time.Second)
		snap, ok := fetchCurrentIfMatch(nctx, exchange, snapshotID)
		cancel()
		if ok {
			return pipeline.ComputeReserveFigures(ctx, snap, prices), "network-current"
		}
	}
	return pipeline.ReserveFiguresSourceMissing(), "missing-summary"
}

func fetchCurrentIfMatch(ctx context.Context, exchange, snapshotID string) (por.Snapshot, bool) {
	switch exchange {
	case "binance":
		_, parsed, err := bapi.FetchSnapshot(ctx)
		if err != nil || parsed.AuditID != snapshotID {
			return por.Snapshot{}, false
		}
		return binance.Normalize(parsed, binance.Meta{Exchange: exchange}), true
	case "okx":
		bundle, _, err := okxapi.FetchCurrentSummary(ctx)
		if err != nil || bundle.Audit.AuditID != snapshotID {
			return por.Snapshot{}, false
		}
		return okx.Normalize(bundle.Audit, okx.Meta{Exchange: exchange}), true
	default:
		return por.Snapshot{}, false
	}
}

func assessmentSnapshotTimeValue(path string) time.Time {
	raw, err := os.ReadFile(path)
	if err != nil {
		return time.Time{}
	}
	var doc struct {
		Snapshot struct {
			SnapshotTime string `json:"snapshotTime"`
		} `json:"snapshot"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, doc.Snapshot.SnapshotTime)
	if err != nil {
		return time.Time{}
	}
	return t
}

func spliceAssessmentFile(path string, fig pipeline.ReserveFigures) (bool, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	updated, err := spliceAssessmentJSON(raw, fig)
	if err != nil {
		return false, err
	}
	if string(updated) == string(raw) {
		return false, nil
	}
	return true, os.WriteFile(path, updated, 0o644)
}

func spliceAssessmentJSON(raw []byte, fig pipeline.ReserveFigures) ([]byte, error) {
	text := string(raw)
	block, err := figureJSONBlock(fig)
	if err != nil {
		return nil, err
	}
	if start := strings.Index(text, "\n  \"totalLiabilityUsd\":"); start >= 0 {
		rest := text[start+1:]
		ratesAt := strings.Index(rest, "\"majorCoinReserveRates\":")
		if ratesAt < 0 {
			return nil, fmt.Errorf("totalLiabilityUsd present without majorCoinReserveRates")
		}
		value := strings.TrimLeft(rest[ratesAt+len("\"majorCoinReserveRates\":"):], " \t")
		if !strings.HasPrefix(value, "[") {
			return nil, fmt.Errorf("majorCoinReserveRates is not an array")
		}
		depth := 0
		consumed := -1
		for i, r := range value {
			switch r {
			case '[':
				depth++
			case ']':
				depth--
				if depth == 0 {
					consumed = i + 1
				}
			}
			if consumed >= 0 {
				break
			}
		}
		if consumed < 0 {
			return nil, fmt.Errorf("could not find end of majorCoinReserveRates")
		}
		after := value[consumed:]
		if strings.HasPrefix(after, ",") {
			after = after[1:]
		}
		if strings.HasPrefix(after, "\n") {
			after = after[1:]
		}
		return []byte(text[:start+1] + block + after), nil
	}
	needle := "\n  \"effectivePoR\":"
	i := strings.Index(text, needle)
	if i < 0 {
		return nil, fmt.Errorf("effectivePoR field not found")
	}
	nl := strings.Index(text[i+1:], "\n")
	if nl < 0 {
		return nil, fmt.Errorf("effectivePoR line has no newline")
	}
	insertAt := i + 1 + nl + 1
	return []byte(text[:insertAt] + block + text[insertAt:]), nil
}

func figureJSONBlock(fig pipeline.ReserveFigures) (string, error) {
	m := fig.AssessmentMap()
	liab := "null"
	if n, ok := m["totalLiabilityUsd"].(json.Number); ok {
		liab = n.String()
	}
	res := "null"
	if n, ok := m["totalReserveUsd"].(json.Number); ok {
		res = n.String()
	}
	unpriced, err := json.Marshal(m["unpricedCoins"])
	if err != nil {
		return "", err
	}
	rates, err := json.MarshalIndent(m["majorCoinReserveRates"], "  ", "  ")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("  \"totalLiabilityUsd\": %s,\n  \"totalReserveUsd\": %s,\n  \"unpricedCoins\": %s,\n  \"majorCoinReserveRates\": %s,\n",
		liab, res, unpriced, rates), nil
}

func spliceReportMarkdown(md string, fig pipeline.ReserveFigures) (string, error) {
	var b strings.Builder
	writeReserveSection(&b, fig)
	section := b.String()
	marker := "## 3. Frequency and Freshness"
	if i := strings.Index(md, reserveSectionHeading); i >= 0 {
		j := strings.Index(md[i:], marker)
		if j < 0 {
			return "", fmt.Errorf("section 3 missing after reserve section")
		}
		return md[:i] + section + md[i+j:], nil
	}
	i := strings.Index(md, marker)
	if i < 0 {
		return "", fmt.Errorf("section 3 heading not found")
	}
	return md[:i] + section + md[i:], nil
}
