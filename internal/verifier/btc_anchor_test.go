package verifier

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ardmere/ardmere/internal/por"
	"github.com/ardmere/ardmere/internal/rpc"
)

func TestBTCAnchorNoTimeAnchor(t *testing.T) {
	v := BTCAnchor{
		SummarySha256: "abc",
		Snapshot: por.Snapshot{
			Exchange: "gateio",
			ID:       "20260316",
		},
	}.Run(context.Background())

	if v.Verdict != VerdictUnverifiable {
		t.Fatalf("verdict=%s reason=%q", v.Verdict, v.Reason)
	}
	if v.Reason != StubReason("gateio", "btc-anchor") {
		t.Fatalf("reason=%q", v.Reason)
	}
}

func TestBTCAnchorPassWithinTolerance(t *testing.T) {
	t.Setenv("RPC_CACHE_DIR", t.TempDir())
	claim := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	blockTS := claim.Add(-22 * time.Minute).Unix()

	srv := esploraTestServer(t, blockTS)
	defer srv.Close()

	old := rpc.EsploraBases
	rpc.EsploraBases = []string{srv.URL}
	defer func() { rpc.EsploraBases = old }()

	v := BTCAnchor{
		SummarySha256: "abc",
		Version:       "1",
		Snapshot: por.Snapshot{
			Exchange:     "binance",
			ID:           "PR01JUN26",
			SnapshotTime: claim,
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

func TestBTCAnchorFailOutsideTolerance(t *testing.T) {
	t.Setenv("RPC_CACHE_DIR", t.TempDir())
	claim := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	blockTS := claim.Add(-45 * time.Minute).Unix()

	srv := esploraTestServer(t, blockTS)
	defer srv.Close()

	old := rpc.EsploraBases
	rpc.EsploraBases = []string{srv.URL}
	defer func() { rpc.EsploraBases = old }()

	v := BTCAnchor{
		SummarySha256: "abc",
		Version:       "1",
		Snapshot: por.Snapshot{
			Exchange:     "binance",
			ID:           "PR01JUN26",
			SnapshotTime: claim,
			TimeAnchor: &por.TimeAnchor{
				Kind:   "btc_block",
				Height: 951913,
			},
		},
		Ledger: rpc.NewLedger(),
	}.Run(context.Background())

	if v.Verdict != VerdictFail {
		t.Fatalf("verdict=%s findings=%+v", v.Verdict, v.Findings)
	}
}

func esploraTestServer(t *testing.T, blockTS int64) *httptest.Server {
	t.Helper()
	const hash = "00000000000000000000d40f10282220ab5904b995daf315010f1c47f8eb32a2"
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/block-height/951913":
			_, _ = w.Write([]byte(hash))
		case "/block/" + hash:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":        hash,
				"height":    951913,
				"timestamp": blockTS,
			})
		default:
			http.NotFound(w, r)
		}
	}))
}
