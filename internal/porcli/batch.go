package porcli

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ardmere/ardmere/internal/porbatch"
	"github.com/ardmere/ardmere/internal/reportgen"
)

func runBatch(args []string) int {
	fs := flag.NewFlagSet("batch", flag.ExitOnError)
	exchanges := fs.String("exchanges", "okx,binance,gateio,bybit,bitget,htx", "comma-separated exchange ids")
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
	return 0
}
