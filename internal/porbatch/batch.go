package porbatch

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ardmere/ardmere/internal/exchangereg"
	"github.com/ardmere/ardmere/internal/exchanges/binance/bapi"
	"github.com/ardmere/ardmere/internal/exchanges/okx/okxapi"
	"github.com/ardmere/ardmere/internal/porrun"
	"github.com/ardmere/ardmere/internal/reportgen"
)

// Options controls a multi-snapshot batch run.
type Options struct {
	Exchanges   []string
	Count       int
	Artifacts   string
	ReportsDir  string
	SkipRPC     bool
	SkipZip     bool
	FullRPCOnly bool
}

// Result summarizes one snapshot run.
type Result struct {
	Exchange string
	Snapshot string
	AnchorOK bool
	ReportOK bool
	Error    string
}

// Run fetches, verifies, and writes reports for the last N snapshots per exchange.
func Run(ctx context.Context, opt Options) ([]Result, error) {
	if opt.Count <= 0 {
		opt.Count = 3
	}
	if opt.Artifacts == "" {
		opt.Artifacts = "./artifacts"
	}
	if opt.ReportsDir == "" {
		opt.ReportsDir = "./docs/reports"
	}
	if len(opt.Exchanges) == 0 {
		opt.Exchanges = []string{"okx", "binance", "gateio", "bybit", "bitget", "htx"}
	}

	var results []Result
	for _, ex := range opt.Exchanges {
		ids, err := listSnapshots(ctx, ex, opt.Count)
		if err != nil {
			results = append(results, Result{Exchange: ex, Error: err.Error()})
			continue
		}
		for i, id := range ids {
			skipRPC := opt.SkipRPC
			if opt.FullRPCOnly && i > 0 {
				skipRPC = true
			}
			r := runOne(ctx, ex, id, opt.Artifacts, opt.ReportsDir, skipRPC, opt.SkipZip)
			results = append(results, r)
		}
	}
	return results, nil
}

func listSnapshots(ctx context.Context, exchangeID string, count int) ([]string, error) {
	switch exchangeID {
	case "okx":
		list, err := okxapi.FetchAuditList(ctx)
		if err != nil {
			return nil, err
		}
		var ids []string
		for i := 0; i < len(list) && i < count; i++ {
			ids = append(ids, list[i].AuditID)
		}
		return ids, nil
	case "binance":
		entries, err := bapi.FetchRecentListEntries(ctx, count)
		if err != nil {
			return nil, err
		}
		var ids []string
		for _, e := range entries {
			ids = append(ids, e.AuditID)
		}
		return ids, nil
	case "gateio", "bybit", "bitget", "htx":
		if _, err := exchangereg.Get(exchangeID); err != nil {
			return nil, err
		}
		return []string{""}, nil
	default:
		return nil, fmt.Errorf("unsupported exchange %q", exchangeID)
	}
}

func runOne(ctx context.Context, exchangeID, auditID, artifactsBase, reportsDir string, skipRPC, skipZip bool) Result {
	res := Result{Exchange: exchangeID, Snapshot: auditID}
	log.Printf("batch: anchor %s audit=%s skip-rpc=%v", exchangeID, auditID, skipRPC)

	ar, err := porrun.RunAnchorPipeline(ctx, porrun.AnchorOpts{
		ExchangeID: exchangeID,
		AuditID:    auditID,
		OutDir:     artifactsBase,
		SkipRPC:    skipRPC,
		SkipZip:    skipZip,
	})
	if err != nil {
		res.Error = err.Error()
		return res
	}
	res.AnchorOK = true
	res.Snapshot = ar.SnapshotID

	if err := reportgen.Write(ctx, reportgen.Options{
		Exchange:   exchangeID,
		SnapshotID: ar.SnapshotID,
		Artifacts:  artifactsBase,
		ReportsDir: reportsDir,
	}); err != nil {
		res.Error = "report: " + err.Error()
		return res
	}
	res.ReportOK = true
	return res
}

// ParseExchangeList splits a comma-separated exchange list.
func ParseExchangeList(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// DefaultTimeout for batch runs.
const DefaultTimeout = 90 * time.Minute
