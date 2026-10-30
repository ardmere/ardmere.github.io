package sitegen

import (
	"fmt"

	"github.com/ardmere/ardmere/internal/exchangeregistry"
)

// UpdateSite regenerates homepage, snapshot history, and exchange comparison pages.
func UpdateSite(opt SiteOptions) error {
	opt = mergeSiteOptions(opt)
	reg, err := exchangeregistry.Load(opt.RegistryPath)
	if err != nil {
		return err
	}
	rows, err := LoadAssessments(opt.ReportsDir, reg)
	if err != nil {
		return err
	}
	if err := updateHomepageFromAssessments(opt, reg, rows); err != nil {
		return fmt.Errorf("homepage: %w", err)
	}
	if err := WriteSnapshotHistory(opt.SnapshotHistoryPath, rows, reg); err != nil {
		return fmt.Errorf("snapshot history: %w", err)
	}
	if err := WriteExchangeComparison(opt.ExchangeComparisonPath, rows, reg); err != nil {
		return fmt.Errorf("exchange comparison: %w", err)
	}
	return nil
}

func mergeSiteOptions(opt SiteOptions) SiteOptions {
	def := defaultSiteOptions()
	if opt.IndexPath == "" {
		opt.IndexPath = def.IndexPath
	}
	if opt.ReportsDir == "" {
		opt.ReportsDir = def.ReportsDir
	}
	if opt.RegistryPath == "" {
		opt.RegistryPath = def.RegistryPath
	}
	if opt.SnapshotHistoryPath == "" {
		opt.SnapshotHistoryPath = def.SnapshotHistoryPath
	}
	if opt.ExchangeComparisonPath == "" {
		opt.ExchangeComparisonPath = def.ExchangeComparisonPath
	}
	if opt.HomepageLimit <= 0 {
		opt.HomepageLimit = def.HomepageLimit
	}
	return opt
}

func updateHomepageFromAssessments(opt SiteOptions, reg exchangeregistry.Registry, rows []Assessment) error {
	reportRows := assessmentsToReportRows(rows)
	if len(reportRows) > opt.HomepageLimit {
		reportRows = reportRows[:opt.HomepageLimit]
	}
	body, err := renderHomepageReportSet(reportRows)
	if err != nil {
		return err
	}
	return patchIndexReportSet(opt.IndexPath, body)
}
