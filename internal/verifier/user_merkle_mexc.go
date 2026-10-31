package verifier

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ardmere/ardmere/internal/exchanges/mexc/merkle"
)

// UserMerkleMexc validates a login-gated MEXC myProof.json Merkle path.
// MEXC truncates SHA-256 to 32 hex characters (128 bits).
type UserMerkleMexc struct {
	ProofPath     string
	ProofSha256   string
	SummaryMerkle string
	AuditID       string
	SnapshotID    string
}

func (v UserMerkleMexc) Run() Verification {
	out := Verification{
		VerifierID: "user-merkle-proof",
		Version:    "mexc-1",
		SnapshotID: v.SnapshotID,
		VerifiedAt: time.Now().UTC(),
		Coverage:   0,
	}
	if v.ProofSha256 != "" {
		out.InputArtifacts = append(out.InputArtifacts, v.ProofSha256)
	}

	if v.ProofPath == "" {
		out.Verdict = VerdictUnverifiable
		out.Reason = "myProof.json not present; download from MEXC PoR page while logged in and import via -user-proof"
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
		out.Reason = "MEXC Merkle path validation failed"
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
		Note:    "user inclusion proof matches MEXC's 128-bit truncated SHA-256 scheme; hash-only root check (balances not re-bound at root)",
	})
	out.Findings = append(out.Findings, Finding{
		Subject: "hashSecurity",
		Field:   "sha256Truncation",
		Claim:   "128-bit hash output",
		Actual:  "collision resistance about 64 bits",
		Status:  VerdictWarn,
		Note:    "MEXC truncates SHA-256 to 32 hex chars; stronger than Bitget's 64-bit truncation but weaker than full SHA-256",
	})

	if v.SummaryMerkle != "" && root != "" {
		st := VerdictPass
		note := "user proof root matches summary merkle root"
		if !mexcRootEqual(v.SummaryMerkle, root) {
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

	if v.AuditID != "" && proofID != "" && v.AuditID != proofID {
		out.Findings = append(out.Findings, Finding{
			Subject: "proofId",
			Field:   "summaryBind",
			Claim:   v.AuditID,
			Actual:  proofID,
			Status:  VerdictWarn,
			Note:    "user proof proofId differs from summary snapshot id",
		})
	}

	out.Verdict = VerdictPass
	out.Reason = "MEXC user Merkle path valid under 128-bit truncated-hash scheme; proves inclusion only for supplied account"
	return out
}

func mexcRootEqual(summaryRoot, proofRoot string) bool {
	summaryRoot = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(summaryRoot)), "0x")
	proofRoot = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(proofRoot)), "0x")
	summaryRoot = strings.TrimLeft(summaryRoot, "0")
	if summaryRoot == "" {
		summaryRoot = "0"
	}
	return summaryRoot == proofRoot
}
