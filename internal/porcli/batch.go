package porcli

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ardmere/ardmere/internal/porbatch"
	"github.com/ardmere/ardmere/internal/reportgen"
)

func runBatch(args []string) int {
	fs := flag.NewFlagSet("batch", flag.ExitOnError)
	exchanges := fs.String("exchanges", "okx,binance,gateio,bybit,bitget,htx,mexc,kucoin,bingx,poloniex,ourbit,coinex,bitunix,woox,toobit,phemex,btcc,kraken", "comma-separated exchange ids")
	count := fs.Int("count", 3, "snapshots per exchange (where supported)")
	artifacts := fs.String("artifacts", "./artifacts", "artifacts root")
	reports := fs.String("reports", "./docs/reports", "reports output directory")
	skipRPC := fs.Bool("skip-rpc", false, "skip all on-chain verifiers")
	fullRPCOnly := fs.Bool("full-rpc-latest-only", true, "run full on-chain checks only on newest snapshot per exchange")
	skipZip := fs.Bool("skip-zip", false, "skip wallet zip download")
	_ = fs.Parse(args)

	ctx, cancel := context.WithTimeout(context.Background(), porbatch.DefaultTimeout)
	defer cancel()

	results, err := porbatch.Run(ctx, porbatch.Options{
		Exchanges:   porbatch.ParseExchangeList(*exchanges),
		Count:       *count,
		Artifacts:   *artifacts,
		ReportsDir:  *reports,
		SkipRPC:     *skipRPC,
		SkipZip:     *skipZip,
		FullRPCOnly: *fullRPCOnly,
	})
	if err != nil {
		log.Fatalf("batch: %v", err)
	}

	var failed int
	for _, r := range results {
		status := "OK"
		if r.Error != "" {
			status = "FAIL: " + r.Error
			failed++
		}
		fmt.Printf("%-8s %-16s %s\n", r.Exchange, r.Snapshot, status)
	}
	if failed > 0 {
		return 1
	}
	refreshSitePages()
	return 0
}

func runReport(args []string) int {
	fs := flag.NewFlagSet("report", flag.ExitOnError)
	exchangeID := fs.String("exchange", "", "exchange id (required)")
	snapshotID := fs.String("snapshot", "", "snapshot id (required)")
	artifacts := fs.String("artifacts", "./artifacts", "artifacts root")
	reports := fs.String("reports", "./docs/reports", "reports output directory")
	_ = fs.Parse(args)

	if strings.TrimSpace(*exchangeID) == "" || strings.TrimSpace(*snapshotID) == "" {
		fmt.Fprintln(os.Stderr, "report requires -exchange and -snapshot")
		return 2
	}
	if err := reportgen.Write(context.Background(), reportgen.Options{
		Exchange:   *exchangeID,
		SnapshotID: *snapshotID,
		Artifacts:  *artifacts,
		ReportsDir: *reports,
	}); err != nil {
		log.Fatalf("report: %v", err)
	}
	fmt.Printf("wrote reports for %s %s\n", *exchangeID, *snapshotID)
	refreshSitePages()
	return 0
}

func runPatchReserveFigures(args []string) int {
	fs := flag.NewFlagSet("patch-reserve-figures", flag.ExitOnError)
	reports := fs.String("reports", "./docs/reports", "reports directory")
	fixtures := fs.String("fixtures", "./fixtures", "fixture summaries")
	artifacts := fs.String("artifacts", "./artifacts", "cached artifact bundles")
	network := fs.Bool("network", false, "use the current Binance or OKX summary when the audit id matches")
	only := fs.String("only", "", "comma-separated exchange/snapshotID filters")
	_ = fs.Parse(args)

	opt := reportgen.PatchOptions{
		ReportsDir:   *reports,
		FixturesDir:  *fixtures,
		ArtifactsDir: *artifacts,
		AllowNetwork: *network,
	}
	if strings.TrimSpace(*only) != "" {
		opt.Only = map[string]struct{}{}
		for _, item := range strings.Split(*only, ",") {
			item = strings.TrimSpace(item)
			if item != "" {
				opt.Only[item] = struct{}{}
			}
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Hour)
	defer cancel()
	results, err := reportgen.PatchReserveFigures(ctx, opt)
	if err != nil {
		log.Fatalf("patch-reserve-figures: %v", err)
	}
	var priced int
	for _, r := range results {
		liab, res := "null", "null"
		if r.TotalLiabilityUsd != nil {
			liab = fmt.Sprintf("%.2f", *r.TotalLiabilityUsd)
			priced++
		}
		if r.TotalReserveUsd != nil {
			res = fmt.Sprintf("%.2f", *r.TotalReserveUsd)
		}
		fmt.Printf("%s %s source=%s amounts=%v liab=%s reserve=%s major=%d unpriced=%d\n",
			r.Exchange, r.SnapshotID, r.Source, r.AmountsPublished, liab, res, r.MajorCoins, r.Unpriced)
	}
	fmt.Printf("patched %d assessment(s), %d with non-null totals\n", len(results), priced)
	return 0
}

func runRefreshReserveDisplay(args []string) int {
	fs := flag.NewFlagSet("refresh-reserve-display", flag.ExitOnError)
	reports := fs.String("reports", "./docs/reports", "reports output directory")
	_ = fs.Parse(args)

	results, err := reportgen.RefreshReserveDisplay(*reports)
	if err != nil {
		log.Fatalf("refresh-reserve-display: %v", err)
	}
	var n int
	for _, r := range results {
		if r.Changed {
			n++
			fmt.Printf("updated %s %s\n", r.Exchange, r.SnapshotID)
		}
	}
	fmt.Printf("refresh-reserve-display: %d changed of %d\n", n, len(results))
	return 0
}

func runRefreshFrequency(args []string) int {
	fs := flag.NewFlagSet("refresh-frequency", flag.ExitOnError)
	reports := fs.String("reports", "./docs/reports", "reports output directory")
	_ = fs.Parse(args)

	results, err := reportgen.RefreshFrequency(*reports)
	if err != nil {
		log.Fatalf("refresh-frequency: %v", err)
	}

	var changed int
	for _, r := range results {
		if r.Changed {
			changed++
			fmt.Printf("updated %s %s\n", r.Exchange, r.SnapshotID)
		}
	}
	fmt.Printf("checked %d report(s), updated %d\n", len(results), changed)
	return 0
}
