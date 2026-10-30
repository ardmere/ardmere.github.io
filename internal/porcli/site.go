package porcli

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ardmere/ardmere/internal/sitegen"
)

func runSite(args []string) int {
	if len(args) == 0 {
		return runSiteAll(nil)
	}
	switch args[0] {
	case "homepage":
		return runSiteHomepage(args[1:])
	case "all":
		return runSiteAll(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown site command %q (try: por site, por site all, por site homepage)\n", args[0])
		return 2
	}
}

func runSiteAll(args []string) int {
	fs := flag.NewFlagSet("site", flag.ExitOnError)
	index := fs.String("index", "index.html", "homepage HTML path")
	reports := fs.String("reports", "./docs/reports", "reports directory")
	registry := fs.String("registry", "config/exchanges/registry.yaml", "exchange registry path")
	snapshotHistory := fs.String("snapshot-history", "docs/reports/snapshot-history.md", "snapshot history markdown path")
	exchangeComparison := fs.String("exchange-comparison", "docs/reports/exchange-comparison.md", "exchange comparison markdown path")
	limit := fs.Int("limit", 6, "maximum homepage report rows")
	_ = fs.Parse(args)

	if err := sitegen.UpdateSite(sitegen.SiteOptions{
		IndexPath:              *index,
		ReportsDir:             *reports,
		RegistryPath:           *registry,
		SnapshotHistoryPath:    *snapshotHistory,
		ExchangeComparisonPath: *exchangeComparison,
		HomepageLimit:          *limit,
	}); err != nil {
		log.Fatalf("site: %v", err)
	}
	fmt.Printf("updated site pages (%s, %s, %s, homepage limit=%d)\n",
		*index, *snapshotHistory, *exchangeComparison, *limit)
	return 0
}

func runSiteHomepage(args []string) int {
	fs := flag.NewFlagSet("site homepage", flag.ExitOnError)
	index := fs.String("index", "index.html", "homepage HTML path")
	reports := fs.String("reports", "./docs/reports", "reports directory")
	registry := fs.String("registry", "config/exchanges/registry.yaml", "exchange registry path")
	limit := fs.Int("limit", 6, "maximum report rows")
	_ = fs.Parse(args)

	if err := sitegen.UpdateHomepageReportSet(sitegen.HomepageOptions{
		IndexPath:    *index,
		ReportsDir:   *reports,
		RegistryPath: *registry,
		Limit:        *limit,
	}); err != nil {
		log.Fatalf("site homepage: %v", err)
	}
	fmt.Printf("updated homepage report set (%s, limit=%d)\n", *index, *limit)
	return 0
}

func refreshSitePages() {
	if err := sitegen.UpdateSite(sitegen.SiteOptions{}); err != nil {
		log.Printf("site pages: %v", err)
	}
}
