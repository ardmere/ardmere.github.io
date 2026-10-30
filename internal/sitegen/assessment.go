package sitegen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ardmere/ardmere/internal/exchangeregistry"
)

// Assessment is a parsed public evaluation assessment with its report artifact.
type Assessment struct {
	Exchange      string
	ExchangeName  string
	SnapshotID    string
	SnapshotTime  time.Time
	PeriodSeq     int
	Stage         string
	Gen           string
	EvidenceLevel string
	Confidence    string
	EffectivePoR  bool
	GeneratedAt   time.Time
	Headline      string
	ArtifactKinds []string
	ReportRel     string // e.g. ./binance/PR01JUL26-transparency-report.md
}

// SiteOptions configures static site page generation.
type SiteOptions struct {
	IndexPath              string
	ReportsDir             string
	RegistryPath           string
	SnapshotHistoryPath    string
	ExchangeComparisonPath string
	HomepageLimit          int
}

func defaultSiteOptions() SiteOptions {
	return SiteOptions{
		IndexPath:              "index.html",
		ReportsDir:             "docs/reports",
		RegistryPath:           exchangeregistry.DefaultPath(),
		SnapshotHistoryPath:    "docs/reports/snapshot-history.md",
		ExchangeComparisonPath: "docs/reports/exchange-comparison.md",
		HomepageLimit:          defaultHomepageLimit,
	}
}

type assessmentJSON struct {
	Exchange      string `json:"exchange"`
	PorStage      string `json:"porStage"`
	Gen           string `json:"gen"`
	EvidenceLevel string `json:"evidenceLevel"`
	Confidence    string `json:"confidence"`
	EffectivePoR  bool   `json:"effectivePoR"`
	GeneratedAt   string `json:"generatedAt"`
	Snapshot      struct {
		ID           string `json:"id"`
		SnapshotTime string `json:"snapshotTime"`
		PeriodSeq    int    `json:"periodSeq"`
	} `json:"snapshot"`
	Summary struct {
		Headline string `json:"headline"`
	} `json:"summary"`
	Artifacts struct {
		Present []struct {
			Kind string `json:"kind"`
		} `json:"present"`
	} `json:"artifacts"`
}

// LoadAssessments scans *-assessment.json files under reportsDir.
func LoadAssessments(reportsDir string, reg exchangeregistry.Registry) ([]Assessment, error) {
	matches, err := filepath.Glob(filepath.Join(reportsDir, "*", "*-assessment.json"))
	if err != nil {
		return nil, err
	}
	var out []Assessment
	for _, path := range matches {
		a, err := parseAssessmentFile(path, reportsDir, reg)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		if a != nil {
			out = append(out, *a)
		}
	}
	sortAssessments(out)
	return out, nil
}

func parseAssessmentFile(path, reportsDir string, reg exchangeregistry.Registry) (*Assessment, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var doc assessmentJSON
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}
	exchange := strings.TrimSpace(doc.Exchange)
	snapshotID := strings.TrimSpace(doc.Snapshot.ID)
	if exchange == "" || snapshotID == "" {
		return nil, fmt.Errorf("missing exchange or snapshot id")
	}
	reportMD := filepath.Join(reportsDir, exchange, snapshotID+"-transparency-report.md")
	if _, err := os.Stat(reportMD); err != nil {
		return nil, nil
	}
	when, err := time.Parse(time.RFC3339, strings.TrimSpace(doc.Snapshot.SnapshotTime))
	if err != nil {
		return nil, fmt.Errorf("snapshotTime: %w", err)
	}
	genAt := time.Time{}
	if doc.GeneratedAt != "" {
		genAt, err = time.Parse(time.RFC3339, doc.GeneratedAt)
		if err != nil {
			return nil, fmt.Errorf("generatedAt: %w", err)
		}
	}
	kinds := make([]string, 0, len(doc.Artifacts.Present))
	for _, p := range doc.Artifacts.Present {
		if k := strings.TrimSpace(p.Kind); k != "" {
			kinds = append(kinds, k)
		}
	}
	stage := strings.TrimSpace(doc.PorStage)
	if stage == "" {
		stage = "Stage 0"
	}
	return &Assessment{
		Exchange:      exchange,
		ExchangeName:  exchangeDisplayName(exchange, reg),
		SnapshotID:    snapshotID,
		SnapshotTime:  when.UTC(),
		PeriodSeq:     doc.Snapshot.PeriodSeq,
		Stage:         stage,
		Gen:           strings.TrimSpace(doc.Gen),
		EvidenceLevel: strings.TrimSpace(doc.EvidenceLevel),
		Confidence:    strings.TrimSpace(doc.Confidence),
		EffectivePoR:  doc.EffectivePoR,
		GeneratedAt:   genAt.UTC(),
		Headline:      strings.TrimSpace(doc.Summary.Headline),
		ArtifactKinds: kinds,
		ReportRel:     "./" + exchange + "/" + snapshotID + "-transparency-report.md",
	}, nil
}

func sortAssessments(rows []Assessment) {
	sort.Slice(rows, func(i, j int) bool {
		if !rows[i].SnapshotTime.Equal(rows[j].SnapshotTime) {
			return rows[i].SnapshotTime.After(rows[j].SnapshotTime)
		}
		return rows[i].SnapshotID > rows[j].SnapshotID
	})
}

func latestByExchange(rows []Assessment, reg exchangeregistry.Registry) []Assessment {
	latest := map[string]Assessment{}
	for _, a := range rows {
		if cur, ok := latest[a.Exchange]; !ok || a.SnapshotTime.After(cur.SnapshotTime) {
			latest[a.Exchange] = a
		}
	}
	var order []string
	for _, ex := range reg.Exchanges {
		if _, ok := latest[ex.ID]; ok {
			order = append(order, ex.ID)
		}
	}
	for id := range latest {
		if !containsString(order, id) {
			order = append(order, id)
		}
	}
	out := make([]Assessment, 0, len(order))
	for _, id := range order {
		out = append(out, latest[id])
	}
	return out
}

func filterExchange(rows []Assessment, exchange string) []Assessment {
	var out []Assessment
	for _, a := range rows {
		if a.Exchange == exchange {
			out = append(out, a)
		}
	}
	sortAssessments(out)
	return out
}

func filterExchanges(rows []Assessment, skip ...string) []Assessment {
	skipSet := map[string]bool{}
	for _, s := range skip {
		skipSet[s] = true
	}
	var out []Assessment
	for _, a := range rows {
		if skipSet[a.Exchange] {
			continue
		}
		out = append(out, a)
	}
	sortAssessments(out)
	return out
}

func latestGeneratedAt(rows []Assessment) time.Time {
	var max time.Time
	for _, a := range rows {
		if a.GeneratedAt.After(max) {
			max = a.GeneratedAt
		}
	}
	return max
}

func assessmentsToReportRows(rows []Assessment) []ReportRow {
	out := make([]ReportRow, 0, len(rows))
	for _, a := range rows {
		out = append(out, ReportRow{
			ExchangeName: a.ExchangeName,
			SnapshotID:   a.SnapshotID,
			SnapshotTime: a.SnapshotTime,
			Stage:        a.Stage,
			ReportHref:   reportHTMLHref(a.ReportRel),
		})
	}
	return out
}

func reportHTMLHref(reportRel string) string {
	rel := strings.TrimPrefix(reportRel, "./")
	return filepath.ToSlash(filepath.Join("docs", "reports", strings.TrimSuffix(rel, ".md")+".html"))
}

func containsString(list []string, want string) bool {
	for _, s := range list {
		if s == want {
			return true
		}
	}
	return false
}

func hasArtifactKind(a Assessment, kind string) bool {
	for _, k := range a.ArtifactKinds {
		if k == kind {
			return true
		}
	}
	return false
}
