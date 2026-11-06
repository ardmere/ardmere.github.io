package verifier

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ardmere/ardmere/internal/exchanges/bitunix/merkle"
)

// UserMerkleBitunix validates a login-gated Bitunix merkle_tree_bitunix.json path.
// Algorithm ported from github.com/BitunixOfficial/proof-of-reserves (full SHA-256).
type UserMerkleBitunix struct {
	ProofPath     string
	ProofSha256   string
	SummaryMerkle string
	AuditID       string
	SnapshotID    string
}

func (v UserMerkleBitunix) Run() Verification {
	out := Verification{
		VerifierID: "user-merkle-proof",
		Version:    "bitunix-1",
		SnapshotID: v.SnapshotID,
		VerifiedAt: time.Now().UTC(),
		Coverage:   0,
	}
	if v.ProofSha256 != "" {
		out.InputArtifacts = append(out.InputArtifacts, v.ProofSha256)
	}
	if v.ProofPath == "" {
		out.Verdict = VerdictUnverifiable
		out.Reason = "merkle_tree_bitunix.json not present; download from Bitunix PoR audit overview while logged in and import via -user-proof"
		return out
	}

	raw, err := os.ReadFile(v.ProofPath)
	if err != nil {
		out.Verdict = VerdictFail
		out.Reason = fmt.Sprintf("read user proof: %v", err)
		return out
	}
	ok, err := merkle.ValidateProof(raw)
	if err != nil {
		out.Verdict = VerdictFail
		out.Reason = fmt.Sprintf("parse user proof: %v", err)
		return out
	}
	if !ok {
		out.Verdict = VerdictFail
		out.Reason = "Bitunix Merkle path validation failed"
		out.Findings = append(out.Findings, Finding{
			Subject: "merklePath",
			Field:   "validate",
			Status:  VerdictFail,
			Note:    "full SHA-256 path does not recompute to published root",
		})
		return out
	}

	root, _ := merkle.RootHash(raw)
	out.Coverage = 1.0
	out.Findings = append(out.Findings, Finding{
		Subject: "merklePath",
		Field:   "validate",
		Status:  VerdictPass,
		Note:    "user inclusion proof matches BitunixOfficial/proof-of-reserves full SHA-256 scheme; proves inclusion only for supplied account",
	})

	if v.SummaryMerkle != "" && root != "" {
		st := VerdictPass
		note := "user proof root matches summary merkle root"
		if !bitunixRootEqual(v.SummaryMerkle, root) {
			st = VerdictFail
			note = "user proof root differs from summary merkle root"
		}
		out.Findings = append(out.Findings, Finding{
			Subject: "merkleRoot",
			Field:   "summaryBind",
			Claim:   normalizeMerkleHex(v.SummaryMerkle),
			Actual:  root,
			Status:  st,
			Note:    note,
		})
		if st == VerdictFail {
			out.Verdict = VerdictFail
			out.Reason = note
			return out
		}
	}

	out.Verdict = VerdictPass
	out.Reason = "Bitunix user Merkle path valid under full SHA-256 scheme; proves inclusion only for supplied account"
	return out
}

func bitunixRootEqual(summaryRoot, proofRoot string) bool {
	summaryRoot = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(summaryRoot)), "0x")
	proofRoot = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(proofRoot)), "0x")
	return summaryRoot == proofRoot
}
