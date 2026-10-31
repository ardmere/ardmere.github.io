package verifier

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ardmere/ardmere/internal/exchanges/bingx/merkle"
)

// UserMerkleBingx validates a BingX user Merkle leaf+path against a published root.
type UserMerkleBingx struct {
	ProofPath     string
	ProofSha256   string
	SummaryMerkle string
	AuditID       string
	SnapshotID    string
}

func (v UserMerkleBingx) Run() Verification {
	out := Verification{
		VerifierID: "user-merkle-proof",
		Version:    "bingx-1",
		SnapshotID: v.SnapshotID,
		VerifiedAt: time.Now().UTC(),
		Coverage:   0,
	}
	if v.ProofSha256 != "" {
		out.InputArtifacts = append(out.InputArtifacts, v.ProofSha256)
	}
	if v.ProofPath == "" {
		out.Verdict = VerdictUnverifiable
		out.Reason = "BingX user Merkle proof not present; download leaf/path while logged in, wrap as JSON, import via -user-proof"
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
		out.Reason = "BingX Merkle path validation failed"
		out.Findings = append(out.Findings, Finding{
			Subject: "merklePath",
			Field:   "validate",
			Status:  VerdictFail,
			Note:    "SHA-256 path does not recompute to published root",
		})
		return out
	}
	root, _ := merkle.RootHash(raw)
	auditID, _ := merkle.AuditID(raw)
	out.Coverage = 1.0
	out.Findings = append(out.Findings, Finding{
		Subject: "merklePath",
		Field:   "validate",
		Status:  VerdictPass,
		Note:    "user inclusion proof matches BingX por.py SHA-256 scheme; proves inclusion only for supplied leaf",
	})
	if v.SummaryMerkle != "" && root != "" {
		st := VerdictPass
		note := "user proof root matches summary merkle root"
		if !bingxRootEqual(v.SummaryMerkle, root) {
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
	if v.AuditID != "" && auditID != "" && v.AuditID != auditID {
		out.Findings = append(out.Findings, Finding{
			Subject: "auditId",
			Field:   "summaryBind",
			Claim:   v.AuditID,
			Actual:  auditID,
			Status:  VerdictWarn,
			Note:    "user proof auditId differs from summary snapshot id",
		})
	}
	out.Verdict = VerdictPass
	out.Reason = "BingX user Merkle path valid; proves inclusion only for supplied account leaf"
	return out
}

func bingxRootEqual(summaryRoot, proofRoot string) bool {
	summaryRoot = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(summaryRoot)), "0x")
	proofRoot = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(proofRoot)), "0x")
	return summaryRoot == proofRoot
}
