package sitegen

import (
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/ardmere/ardmere/internal/exchangeregistry"
)

const (
	homepageReportSetStart = "<!-- HOMEPAGE_REPORT_SET_START -->"
	homepageReportSetEnd   = "<!-- HOMEPAGE_REPORT_SET_END -->"
	defaultHomepageLimit   = 6
)

// HomepageOptions configures homepage report-list generation.
type HomepageOptions struct {
	IndexPath    string
	ReportsDir   string
	RegistryPath string
	Limit        int
}

// ReportRow is one public transparency report entry for the homepage list.
type ReportRow struct {
	ExchangeName string
	SnapshotID   string
	SnapshotTime time.Time
	Stage        string
	ReportHref   string
}

// ListPublicReports scans *-assessment.json files under reportsDir.
func ListPublicReports(reportsDir string, reg exchangeregistry.Registry) ([]ReportRow, error) {
	rows, err := LoadAssessments(reportsDir, reg)
	if err != nil {
		return nil, err
	}
	return assessmentsToReportRows(rows), nil
}

func exchangeDisplayName(id string, reg exchangeregistry.Registry) string {
	if ex := reg.Find(id); ex != nil && ex.Name != "" {
		return ex.Name
	}
	if id == "" {
		return ""
	}
	return strings.ToUpper(id[:1]) + id[1:]
}

// UpdateHomepageReportSet rewrites the marked section in index.html with the latest N reports.
func UpdateHomepageReportSet(opt HomepageOptions) error {
	siteOpt := mergeSiteOptions(SiteOptions{
		IndexPath:     opt.IndexPath,
		ReportsDir:    opt.ReportsDir,
		RegistryPath:  opt.RegistryPath,
		HomepageLimit: opt.Limit,
	})
	reg, err := exchangeregistry.Load(siteOpt.RegistryPath)
	if err != nil {
		return err
	}
	rows, err := LoadAssessments(siteOpt.ReportsDir, reg)
	if err != nil {
		return err
	}
	return updateHomepageFromAssessments(siteOpt, reg, rows)
}

func renderHomepageReportSet(rows []ReportRow) (string, error) {
	const tpl = `              <div class="report-list-head">
                <span>Exchange</span>
                <span>Snapshot</span>
                <span>Date</span>
                <span>Stage</span>
                <span class="report-list-action">Report</span>
              </div>
{{range .}}              <div class="report-row">
                <span class="exchange">{{.ExchangeName}}</span>
                <time class="report-snapshot" datetime="{{.SnapshotISO}}">{{.SnapshotID}}</time>
                <span class="report-date"><time datetime="{{.SnapshotISO}}">{{.SnapshotDate}}</time></span>
                <span class="{{.StageClass}}">{{.Stage}}</span>
                <a class="report-link" href="{{.ReportHref}}">View report</a>
              </div>
{{end}}`
	type view struct {
		ExchangeName string
		SnapshotID   string
		SnapshotISO  string
		SnapshotDate string
		Stage        string
		StageClass   string
		ReportHref   string
	}
	views := make([]view, 0, len(rows))
	for _, r := range rows {
		stageClass := "stage"
		if r.Stage == "Stage 1" {
			stageClass = "stage stage-1"
		}
		views = append(views, view{
			ExchangeName: r.ExchangeName,
			SnapshotID:   r.SnapshotID,
			SnapshotISO:  r.SnapshotTime.Format(time.RFC3339),
			SnapshotDate: r.SnapshotTime.Format("2006-01-02"),
			Stage:        r.Stage,
			StageClass:   stageClass,
			ReportHref:   r.ReportHref,
		})
	}
	t, err := template.New("homepage").Parse(tpl)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	if err := t.Execute(&b, views); err != nil {
		return "", err
	}
	return b.String(), nil
}

func patchIndexReportSet(indexPath, body string) error {
	raw, err := os.ReadFile(indexPath)
	if err != nil {
		return err
	}
	content := string(raw)
	start := strings.Index(content, homepageReportSetStart)
	end := strings.Index(content, homepageReportSetEnd)
	if start < 0 || end < 0 || end <= start {
		return fmt.Errorf("%s: missing homepage report set markers", indexPath)
	}
	start += len(homepageReportSetStart)
	var out strings.Builder
	out.WriteString(content[:start])
	out.WriteString("\n")
	out.WriteString(body)
	out.WriteString("              ")
	out.WriteString(homepageReportSetEnd)
	out.WriteString(content[end+len(homepageReportSetEnd):])
	return os.WriteFile(indexPath, []byte(out.String()), 0o644)
}
