package sitegen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ardmere/ardmere/internal/exchangeregistry"
)

func TestListPublicReportsOrdersBySnapshotTime(t *testing.T) {
	dir := t.TempDir()
	writeAssessment(t, dir, "binance", "PR01JUN26", "2026-06-01T00:00:00Z", "Stage 0")
	writeAssessment(t, dir, "binance", "PR01JUL26", "2026-07-01T00:00:00Z", "Stage 0")
	writeAssessment(t, dir, "okx", "508399035", "2026-06-18T16:00:00Z", "Stage 1")

	reg := exchangeregistry.Registry{
		Exchanges: []exchangeregistry.Exchange{
			{ID: "binance", Name: "Binance"},
			{ID: "okx", Name: "OKX"},
		},
	}
	rows, err := ListPublicReports(dir, reg)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 3 {
		t.Fatalf("got %d rows, want 3", len(rows))
	}
	if rows[0].SnapshotID != "PR01JUL26" {
		t.Fatalf("first row=%s want PR01JUL26", rows[0].SnapshotID)
	}
	if rows[1].SnapshotID != "508399035" {
		t.Fatalf("second row=%s want 508399035", rows[1].SnapshotID)
	}
	if rows[0].ReportHref != "docs/reports/binance/PR01JUL26-transparency-report.html" {
		t.Fatalf("href=%q", rows[0].ReportHref)
	}
}

func TestUpdateHomepageReportSet(t *testing.T) {
	root := t.TempDir()
	reports := filepath.Join(root, "docs", "reports")
	regPath := filepath.Join(root, "config", "exchanges", "registry.yaml")
	indexPath := filepath.Join(root, "index.html")

	if err := os.MkdirAll(filepath.Join(root, "config", "exchanges"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(regPath, []byte(`version: 1
updated: "2026-07-01"
exchanges:
  - id: binance
    name: Binance
    porPage: https://example.com
`), 0o644); err != nil {
		t.Fatal(err)
	}
	writeAssessment(t, reports, "binance", "PR01JUL26", "2026-07-01T00:00:00Z", "Stage 0")

	index := `<div class="report-list">
              <!-- HOMEPAGE_REPORT_SET_START -->
              OLD
              <!-- HOMEPAGE_REPORT_SET_END -->
            </div>`
	if err := os.WriteFile(indexPath, []byte(index), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := UpdateHomepageReportSet(HomepageOptions{
		IndexPath:    indexPath,
		ReportsDir:   reports,
		RegistryPath: regPath,
		Limit:        6,
	}); err != nil {
		t.Fatal(err)
	}
	out, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if strings.Contains(s, "OLD") {
		t.Fatal("old content not replaced")
	}
	if !strings.Contains(s, "PR01JUL26") || !strings.Contains(s, "Binance") {
		t.Fatalf("missing generated row: %s", s)
	}
	if !strings.Contains(s, homepageReportSetStart) || !strings.Contains(s, homepageReportSetEnd) {
		t.Fatal("markers removed")
	}
}

func writeAssessment(t *testing.T, reportsDir, exchange, snapshotID, snapshotTime, stage string) {
	t.Helper()
	dir := filepath.Join(reportsDir, exchange)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	md := filepath.Join(dir, snapshotID+"-transparency-report.md")
	if err := os.WriteFile(md, []byte("# test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	jsonPath := filepath.Join(dir, snapshotID+"-assessment.json")
	payload := `{
  "exchange": "` + exchange + `",
  "porStage": "` + stage + `",
  "snapshot": {
    "id": "` + snapshotID + `",
    "snapshotTime": "` + snapshotTime + `"
  }
}`
	if err := os.WriteFile(jsonPath, []byte(payload), 0o644); err != nil {
		t.Fatal(err)
	}
	_ = time.Now()
}
