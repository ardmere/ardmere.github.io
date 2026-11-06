package verifier

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ardmere/ardmere/internal/exchanges/toobit/merkle"
)

// UserMerkleToobit validates a login-gated Toobit SMT proof JSON.
// Algorithm ported from github.com/toobit-docs/por-verifier (SHA256-SMT).
type UserMerkleToobit struct {
	ProofPath     string
	ProofSha256   string
	SummaryMerkle string
	AuditID       string
	SnapshotID    string
}

func (v UserMerkleToobit) Run() Verification {
	out := Verification{
		VerifierID: "user-merkle-proof",
		Version:    "toobit-1",
		SnapshotID: v.SnapshotID,
		VerifiedAt: time.Now().UTC(),
		Coverage:   0,
	}
	if v.ProofSha256 != "" {
		out.InputArtifacts = append(out.InputArtifacts, v.ProofSha256)
	}
	if v.ProofPath == "" {
		out.Verdict = VerdictUnverifiable
		out.Reason = "Toobit proof JSON not present; download proof_of_reserve_{auditId}.json while logged in and import via -user-proof"
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
		out.Reason = "Toobit SMT Merkle path validation failed"
		out.Findings = append(out.Findings, Finding{
			Subject: "merklePath",
			Field:   "validate",
			Status:  VerdictFail,
			Note:    "SHA256-SMT path does not recompute to published root",
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
		Note:    "user inclusion proof matches toobit-docs/por-verifier SHA256-SMT scheme; proves inclusion only for supplied account",
	})

	if v.SummaryMerkle != "" && root != "" {
		st := VerdictPass
		note := "user proof root matches summary merkle root"
		if !toobitRootEqual(v.SummaryMerkle, root) {
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
	if v.AuditID != "" && auditID != "" && v.AuditID != auditID && !strings.EqualFold(v.AuditID, auditID) {
		// Official audit ids look like Toobit_POR_202610; calendar snapshot ids are YYYYMMDD.
		out.Findings = append(out.Findings, Finding{
			Subject: "auditId",
			Field:   "summaryBind",
			Claim:   v.AuditID,
			Actual:  auditID,
			Status:  VerdictWarn,
			Note:    "user proof audit_id differs from summary snapshot id (expected for calendar vs official id)",
		})
	}

	out.Verdict = VerdictPass
	out.Reason = "Toobit user SMT Merkle path valid; proves inclusion only for supplied account"
	return out
}

func toobitRootEqual(summaryRoot, proofRoot string) bool {
	summaryRoot = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(summaryRoot)), "0x")
	proofRoot = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(proofRoot)), "0x")
	return summaryRoot == proofRoot
}
