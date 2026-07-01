package verifier

import (
	"testing"

	"github.com/ardmere/ardmere/internal/verifier/onchainconfig"
)

func TestLedgerLiveSnapshotNoteSkipsWhenHistorical(t *testing.T) {
	spec := onchainconfig.LedgerSpec{LiveSnapshot: true, Kind: onchainconfig.LedgerSolanaSPL}
	if got := ledgerLiveSnapshotNote(spec, map[string]string{"historical_slot": "423478906"}); got != "" {
		t.Fatalf("expected empty when historical, got %q", got)
	}
	if got := ledgerLiveSnapshotNote(spec, nil); got == "" {
		t.Fatal("expected live snapshot note without historical")
	}
}
