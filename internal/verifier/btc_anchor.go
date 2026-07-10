package verifier

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/ardmere/ardmere/internal/por"
	"github.com/ardmere/ardmere/internal/rpc"
)

const btcAnchorMaxDrift = 30 * time.Minute

// BTCAnchor checks that a snapshot-declared BTC block height timestamp is
// consistent with the exchange's declared snapshot time (within 30 minutes).
type BTCAnchor struct {
	SummarySha256 string
	Snapshot      por.Snapshot
	Version       string
	Ledger        *rpc.LedgerClient
}

func (v BTCAnchor) Run(ctx context.Context) Verification {
	ver := v.Version
	if ver == "" {
		ver = "1"
	}
	out := Verification{
		VerifierID:     "btc-anchor",
		Version:        ver,
		SnapshotID:     v.Snapshot.ID,
		VerifiedAt:     time.Now().UTC(),
		InputArtifacts: artifactInputs(v.SummarySha256),
		Coverage:       1.0,
	}

	anchor := v.Snapshot.TimeAnchor
	if anchor == nil || anchor.Height == 0 {
		out.Verdict = VerdictUnverifiable
		out.Coverage = 0
		out.Reason = StubReason(v.Snapshot.Exchange, "btc-anchor")
		return out
	}
	if anchor.Kind != "" && anchor.Kind != "btc_block" {
		out.Verdict = VerdictUnverifiable
		out.Coverage = 0
		out.Reason = fmt.Sprintf("unsupported time anchor kind %q", anchor.Kind)
		return out
	}

	ledger := v.Ledger
	if ledger == nil {
		ledger = rpc.NewLedger()
	}

	height := int64(anchor.Height)
	blockTime, provider, err := ledger.BTCBlockTimeAtHeight(ctx, height)
	if err != nil {
		out.Verdict = VerdictFail
		out.Reason = err.Error()
		out.Findings = append(out.Findings, Finding{
			Subject: fmt.Sprintf("btc:%d", height),
			Field:   "blockTimestamp",
			Claim:   v.Snapshot.SnapshotTime.UTC().Format(time.RFC3339),
			Status:  VerdictFail,
			Note:    fmt.Sprintf("chain query failed via %s: %v", provider, err),
		})
		return out
	}

	claimTime := v.Snapshot.SnapshotTime.UTC()
	drift := blockTime.Sub(claimTime)
	driftAbs := time.Duration(math.Abs(float64(drift)))

	st := VerdictPass
	note := fmt.Sprintf("BTC block %d timestamp via %s; drift %s (max %s)", height, provider, driftAbs, btcAnchorMaxDrift)
	if driftAbs > btcAnchorMaxDrift {
		st = VerdictFail
	}

	out.Findings = append(out.Findings, Finding{
		Subject: fmt.Sprintf("btc:%d", height),
		Field:   "snapshotTime",
		Claim:   claimTime.Format(time.RFC3339),
		Actual:  blockTime.Format(time.RFC3339),
		Status:  st,
		Note:    note,
	})
	out.Verdict = st
	return out
}
