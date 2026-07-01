package verifier

import (
	"fmt"
	"os"
	"time"

	"github.com/ardmere/ardmere/internal/exchanges/bybit/merkle"
)

// UserMerkleBybit validates a login-gated Bybit myProof.json Merkle path (v5 schema).
type UserMerkleBybit struct {
	ProofPath     string
	ProofSha256   string
	SummaryMerkle string
	AuditID       string
	SnapshotID    string
}

func (v UserMerkleBybit) Run() Verification {
	out := Verification{
		VerifierID: "user-merkle-proof",
		Version:    "bybit-1",
		SnapshotID: v.SnapshotID,
		VerifiedAt: time.Now().UTC(),
		Coverage: 0,
	}
	if v.ProofSha256 != "" {
		out.InputArtifacts = append(out.InputArtifacts, v.ProofSha256)
	}

	if v.ProofPath == "" {
		out.Verdict = VerdictUnverifiable
		out.Reason = "myProof.json not present; download from Bybit PoR page while logged in and import via -user-proof"
		return out
	}

	raw, err := os.ReadFile(v.ProofPath)
	if err != nil {
		out.Verdict = VerdictFail
		out.Reason = fmt.Sprintf("read user proof: %v", err)
		return out
	}

	ok, err := merkle.ValidateProofV5(raw)
	if err != nil {
		out.Verdict = VerdictFail
		out.Reason = fmt.Sprintf("parse user proof: %v", err)
		return out
	}
	if !ok {
		out.Verdict = VerdictFail
		out.Reason = "Merkle path validation failed (hash/balance/height chain)"
		out.Findings = append(out.Findings, Finding{
			Subject: "merklePath",
			Field:   "validate",
			Status:  VerdictFail,
			Note:    "SHA-256 binary tree path does not recompute to published root",
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
		Note:    "user inclusion proof cryptographically valid for v5 schema; does not prove exchange-wide solvency",
	})

	if v.SummaryMerkle != "" && root != "" {
		st := VerdictPass
		note := "user proof root matches summary merkle root"
		if !merkleHexEqual(v.SummaryMerkle, root) {
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
	out.Reason = "user Merkle path valid (Gen 1.5 SHA-256); proves inclusion only for supplied account"
	return out
}
