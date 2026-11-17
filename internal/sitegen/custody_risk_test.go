package sitegen

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/ardmere/ardmere/internal/exchangeregistry"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		goMod := filepath.Join(dir, "go.mod")
		if raw, err := os.ReadFile(goMod); err == nil && strings.Contains(string(raw), "github.com/ardmere/ardmere") {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// TestCommittedSiteCustodyRiskScores asserts homepage and snapshot-history
// show custody scores whenever custody-risk-scores.json has a numeric score for
// that exchange/snapshot. Catches drift when scores are backfilled without
// running `go run ./cmd/por site`.
func TestCommittedSiteCustodyRiskScores(t *testing.T) {
	root := repoRoot(t)
	if root == "" {
		t.Skip("not inside ardmere repo")
	}
	opt := mergeSiteOptions(SiteOptions{
		IndexPath:           filepath.Join(root, "index.html"),
		ReportsDir:          filepath.Join(root, "docs/reports"),
		RegistryPath:        filepath.Join(root, "config/exchanges/registry.yaml"),
		SnapshotHistoryPath: filepath.Join(root, "docs/reports/snapshot-history.md"),
	})
	if _, err := os.Stat(filepath.Join(opt.ReportsDir, "custody-risk-scores.json")); err != nil {
		t.Skip("custody-risk-scores.json not present")
	}

	reg, err := exchangeregistry.Load(opt.RegistryPath)
	if err != nil {
		t.Fatal(err)
	}
	assessments, err := LoadAssessments(opt.ReportsDir, reg)
	if err != nil {
		t.Fatal(err)
	}
	custody, err := LoadCustodyRiskScores(opt.ReportsDir)
	if err != nil {
		t.Fatal(err)
	}

	indexHTML, err := os.ReadFile(opt.IndexPath)
	if err != nil {
		t.Fatal(err)
	}
	historyMD, err := os.ReadFile(opt.SnapshotHistoryPath)
	if err != nil {
		t.Fatal(err)
	}

	reportRows := assessmentsToReportRows(assessments)
	applyCustodyRisk(reportRows, assessments, custody)
	if len(reportRows) > opt.HomepageLimit {
		reportRows = reportRows[:opt.HomepageLimit]
	}

	rowSnippet := regexp.MustCompile(`(?s)<div class="report-row">.*?</div>`)

	for i, row := range reportRows {
		if i >= len(assessments) {
			break
		}
		a := assessments[i]
		score, _, ok := custody.Lookup(a.Exchange, a.SnapshotID)
		if !ok {
			continue
		}
		if row.CustodyRisk == "" {
			t.Fatalf("sitegen would render — for homepage %s/%s but custody index has %s", a.Exchange, a.SnapshotID, score)
		}
		if row.CustodyRisk != score {
			t.Fatalf("homepage row %s: want custody %s, sitegen has %s", a.SnapshotID, score, row.CustodyRisk)
		}
		var snippet string
		for _, block := range rowSnippet.FindAllString(string(indexHTML), -1) {
			if strings.Contains(block, row.SnapshotID) {
				snippet = block
				break
			}
		}
		if snippet == "" {
			t.Fatalf("index.html missing homepage row for %s", a.SnapshotID)
		}
		if !strings.Contains(snippet, score) {
			t.Fatalf("index.html stale for %s: expected transparency risk %q in row, got: %s", a.SnapshotID, score, snippet)
		}
	}

	for _, a := range assessments {
		if _, _, ok := custody.Lookup(a.Exchange, a.SnapshotID); !ok {
			continue
		}
		line := snapshotHistoryLine(string(historyMD), a.ExchangeName, a.SnapshotID)
		if line == "" {
			t.Fatalf("snapshot-history.md missing row for %s/%s", a.Exchange, a.SnapshotID)
		}
		if strings.Contains(line, "—") {
			t.Fatalf("snapshot-history.md shows — for %s/%s but custody index has a score: %s", a.Exchange, a.SnapshotID, line)
		}
	}
}

func snapshotHistoryLine(md, exchangeName, snapshotID string) string {
	for _, line := range strings.Split(md, "\n") {
		trim := strings.TrimSpace(line)
		if !strings.HasPrefix(trim, "|") || !strings.Contains(line, "`"+snapshotID+"`") {
			continue
		}
		fields := strings.Split(line, "|")
		if len(fields) < 3 {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(fields[1]), exchangeName) {
			return line
		}
	}
	return ""
}