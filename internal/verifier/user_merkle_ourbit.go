package verifier

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ardmere/ardmere/internal/exchanges/ourbit/merkle"
)

// UserMerkleOurbit validates a login-gated Ourbit merkel_tree_ourbit.json Merkle path.
// Ourbit truncates SHA-256 to 32 hex characters (128 bits).
type UserMerkleOurbit struct {
	ProofPath     string
	ProofSha256   string
	SummaryMerkle string
	AuditID       string
	SnapshotID    string
}

func (v UserMerkleOurbit) Run() Verification {
	out := Verification{
		VerifierID: "user-merkle-proof",
		Version:    "ourbit-1",
		SnapshotID: v.SnapshotID,
		VerifiedAt: time.Now().UTC(),
		Coverage:   0,
	}
	if v.ProofSha256 != "" {
		out.InputArtifacts = append(out.InputArtifacts, v.ProofSha256)
	}

	if v.ProofPath == "" {
		out.Verdict = VerdictUnverifiable
		out.Reason = "merkel_tree_ourbit.json not present; download from Ourbit PoR page while logged in and import via -user-proof"
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
		out.Reason = "Ourbit Merkle path validation failed"
		out.Findings = append(out.Findings, Finding{
			Subject: "merklePath",
			Field:   "validate",
			Status:  VerdictFail,
			Note:    "128-bit truncated SHA-256 path does not recompute to published root",
		})
		return out
	}

	root, _ := merkle.RootHash(raw)
	proofID, _ := merkle.ProofID(raw)
	out.Coverage = 1.0
	out.Findings = append(out.Findings, Finding{
		Subject: "merklePath",
		Field:   "validate",
		Status:  VerdictPass,
		Note:    "user inclusion proof matches Ourbit's 128-bit truncated SHA-256 scheme; hash-only root check",
	})
	out.Findings = append(out.Findings, Finding{
		Subject: "hashSecurity",
		Field:   "sha256Truncation",
		Claim:   "128-bit hash output",
		Actual:  "collision resistance about 64 bits",
		Status:  VerdictWarn,
		Note:    "Ourbit truncates SHA-256 to 32 hex chars (same class as MEXC); weaker than full SHA-256",
	})

	if v.SummaryMerkle != "" && root != "" {
		st := VerdictPass
		note := "user proof root matches summary merkle root"
		if !ourbitRootEqual(v.SummaryMerkle, root) {
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
	} else if v.SummaryMerkle == "" {
		out.Findings = append(out.Findings, Finding{
			Subject: "merkleRoot",
			Field:   "summaryBind",
			Status:  VerdictWarn,
			Note:    "public stock_info has no period merkle root; inclusion verified against proof file root only",
		})
	}

	if v.AuditID != "" && proofID != "" {
		// AuditID is YYYYMMDD; proofId is PR-YYYY-MM-DD
		want := "PR-" + formatOurbitProofDate(v.AuditID)
		if want != proofID && v.AuditID != proofID {
			out.Findings = append(out.Findings, Finding{
				Subject: "proofId",
				Field:   "summaryBind",
				Claim:   v.AuditID,
				Actual:  proofID,
				Status:  VerdictWarn,
				Note:    "user proof proofId differs from summary snapshot id",
			})
		}
	}

	out.Verdict = VerdictPass
	out.Reason = "Ourbit user Merkle path valid under 128-bit truncated-hash scheme; proves inclusion only for supplied account"
	return out
}

func ourbitRootEqual(summaryRoot, proofRoot string) bool {
	summaryRoot = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(summaryRoot)), "0x")
	proofRoot = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(proofRoot)), "0x")
	summaryRoot = strings.TrimLeft(summaryRoot, "0")
	if summaryRoot == "" {
		summaryRoot = "0"
	}
	return summaryRoot == proofRoot
}

func formatOurbitProofDate(auditID string) string {
	// 20261101 → 2026-11-01
	if len(auditID) == 8 {
		return auditID[0:4] + "-" + auditID[4:6] + "-" + auditID[6:8]
	}
	return auditID
}
