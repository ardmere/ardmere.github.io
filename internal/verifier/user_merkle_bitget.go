package verifier

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ardmere/ardmere/internal/exchanges/bitget/merkle"
)

// UserMerkleBitget validates Bitget's login-gated merkel_tree_bg.json.
// Bitget truncates SHA-256 to 64 bits, so PASS is weak inclusion evidence.
type UserMerkleBitget struct {
	ProofPath     string
	ProofSha256   string
	SummaryMerkle string
	AuditID       string
	SnapshotID    string
}

func (v UserMerkleBitget) Run() Verification {
	out := Verification{
		VerifierID: "user-merkle-proof",
		Version:    "bitget-1",
		SnapshotID: v.SnapshotID,
		VerifiedAt: time.Now().UTC(),
		Coverage:   0,
	}
	if v.ProofSha256 != "" {
		out.InputArtifacts = append(out.InputArtifacts, v.ProofSha256)
	}

	if v.ProofPath == "" {
		out.Verdict = VerdictUnverifiable
		out.Reason = "merkel_tree_bg.json not present; download from Bitget while logged in and import via -user-proof"
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
		out.Reason = "Bitget Merkle path validation failed"
		out.Findings = append(out.Findings, Finding{
			Subject: "merklePath",
			Field:   "validate",
			Status:  VerdictFail,
			Note:    "64-bit truncated SHA-256 path does not recompute to published root",
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
		Note:    "user inclusion proof matches Bitget's 64-bit truncated SHA-256 scheme; collision resistance is only ~32 bits",
	})
	out.Findings = append(out.Findings, Finding{
		Subject: "hashSecurity",
		Field:   "sha256Truncation",
		Claim:   "64-bit hash output",
		Actual:  "collision resistance about 32 bits",
		Status:  VerdictWarn,
		Note:    "Bitget truncates SHA-256 to 16 hex chars; Merkle PASS is weak evidence and should not be treated like full SHA-256",
	})

	if v.SummaryMerkle != "" && root != "" {
		st := VerdictPass
		note := "user proof root matches summary merkle root"
		if !bitgetRootEqual(v.SummaryMerkle, root) {
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
	out.Reason = "Bitget user Merkle path valid under weak 64-bit truncated-hash scheme; proves inclusion only for supplied account"
	return out
}

func bitgetRootEqual(summaryRoot, proofRoot string) bool {
	summaryRoot = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(summaryRoot)), "0x")
	proofRoot = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(proofRoot)), "0x")
	summaryRoot = strings.TrimLeft(summaryRoot, "0")
	if summaryRoot == "" {
		summaryRoot = "0"
	}
	return summaryRoot == proofRoot
}
