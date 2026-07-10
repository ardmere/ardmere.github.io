package verifier

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/ardmere/ardmere/internal/por"
	"github.com/ardmere/ardmere/internal/rpc"
)

func TestBTCAnchorLivePR01JUN26(t *testing.T) {
	if os.Getenv("BTC_ANCHOR_LIVE") == "" {
		t.Skip("set BTC_ANCHOR_LIVE=1 to run live BTC anchor probe")
	}

	v := BTCAnchor{
		SummarySha256: "live",
		Version:       "1",
		Snapshot: por.Snapshot{
			Exchange:     "binance",
			ID:           "PR01JUN26",
			SnapshotTime: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
			TimeAnchor: &por.TimeAnchor{
				Kind:   "btc_block",
				Height: 951913,
			},
		},
		Ledger: rpc.NewLedger(),
	}.Run(context.Background())

	if v.Verdict != VerdictPass {
		t.Fatalf("verdict=%s findings=%+v reason=%q", v.Verdict, v.Findings, v.Reason)
	}
}
