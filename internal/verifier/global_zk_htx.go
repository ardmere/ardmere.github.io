package verifier

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ardmere/ardmere/internal/exchanges/htx/htxapi"
)

const htxZKBoundaryNote = "HTX Groth16 proof verifies batch SMT insertion and platform cumDebt≤cumEquity only; " +
	"does not prove per-user solvency, asset sum semantics, on-chain reserves, or summary↔config binding without independent fetch"

// GlobalZKProofHTX validates HTX Groth16 liability bundle (public-data.zip).
type GlobalZKProofHTX struct {
	SnapshotID      string
	SummarySha256   string
	ProofBundlePath string
	ProofBundleSha  string
	MerkleRoot      string
}

func (v GlobalZKProofHTX) Run() Verification {
	out := Verification{
		VerifierID:     "global-zk-proof",
		Version:        "htx-1",
		SnapshotID:     v.SnapshotID,
		VerifiedAt:     time.Now().UTC(),
		InputArtifacts: artifactInputs(v.SummarySha256, v.ProofBundleSha),
	}

	if v.ProofBundlePath == "" {
		out.Verdict = VerdictUnverifiable
		out.Reason = "HTX public-data zk bundle not present; download from PoR Reports or GitHub release public-data.zip"
		return out
	}

	meta, err := htxapi.ParseZkBundleZip(v.ProofBundlePath)
	if err != nil {
		out.Verdict = VerdictFail
		out.Reason = fmt.Sprintf("parse zk bundle: %v", err)
		return out
	}

	out.Findings = append(out.Findings,
		Finding{
			Subject: "public-data.zip",
			Field:   "structure",
			Status:  VerdictPass,
			Note:    fmt.Sprintf("config.json + proof0.csv + zkpor500.vk.save; %d batches, %d assets", meta.BatchCount, meta.AssetCount),
		},
		Finding{
			Subject: "zk-liability",
			Field:   "proof_boundary",
			Status:  VerdictPass,
			Note:    htxZKBoundaryNote,
		},
	)

	if v.MerkleRoot != "" {
		if merkleHexEqual(v.MerkleRoot, meta.MerkleRoot) {
			out.Findings = append(out.Findings, Finding{
				Subject: "merkleRoot",
				Field:   "summary_vs_zk_csv",
				Claim:   normalizeMerkleHex(v.MerkleRoot),
				Actual:  meta.MerkleRoot,
				Status:  VerdictPass,
				Note:    "summary merkle root matches final AccountTreeRoots in proof0.csv",
			})
		} else {
			out.Findings = append(out.Findings, Finding{
				Subject: "merkleRoot",
				Field:   "summary_vs_zk_csv",
				Claim:   normalizeMerkleHex(v.MerkleRoot),
				Actual:  meta.MerkleRoot,
				Status:  VerdictFail,
			})
			out.Verdict = VerdictFail
			out.Reason = "summary merkle root does not match zk proof csv final account tree root"
			return out
		}
	}

	if v.SnapshotID != "" && meta.AuditID != "" && v.SnapshotID != meta.AuditID {
		out.Findings = append(out.Findings, Finding{
			Subject: "auditId",
			Field:   "snapshot_id",
			Claim:   v.SnapshotID,
			Actual:  meta.AuditID,
			Status:  VerdictFail,
		})
		out.Verdict = VerdictFail
		out.Reason = fmt.Sprintf("snapshot id %s != zk-derived audit id %s", v.SnapshotID, meta.AuditID)
		return out
	}

	bin := os.Getenv("HTX_ZK_VERIFIER")
	if bin == "" {
		out.Verdict = VerdictPartial
		out.Coverage = 0.7
		out.Reason = "structure and merkle bind OK; set HTX_ZK_VERIFIER to zkverifiermac path for Groth16 cryptographic verification"
		return out
	}

	workDir, err := os.MkdirTemp("", "htx-zk-*")
	if err != nil {
		out.Verdict = VerdictFail
		out.Reason = fmt.Sprintf("temp dir: %v", err)
		return out
	}
	defer os.RemoveAll(workDir)

	if _, err := htxapi.ExtractZkBundleToDir(v.ProofBundlePath, workDir); err != nil {
		out.Verdict = VerdictFail
		out.Reason = fmt.Sprintf("extract zk bundle: %v", err)
		return out
	}

	cmd := exec.Command(bin)
	cmd.Dir = workDir
	outBytes, err := cmd.CombinedOutput()
	text := string(outBytes)
	if err != nil {
		out.Verdict = VerdictFail
		out.Reason = fmt.Sprintf("HTX zk verifier: %v (%s)", err, trimOutput(outBytes))
		return out
	}
	if strings.Contains(text, "All proofs verify passed!!!") {
		out.Verdict = VerdictPass
		out.Coverage = 1.0
		out.Reason = "HTX zkverifier succeeded; see proof boundary finding for cryptographic limits"
		if root := parseHTXVerifierRoot(text); root != "" && v.MerkleRoot != "" && !merkleHexEqual(v.MerkleRoot, root) {
			out.Findings = append(out.Findings, Finding{
				Subject: "merkleRoot",
				Field:   "verifier_vs_summary",
				Claim:   normalizeMerkleHex(v.MerkleRoot),
				Actual:  root,
				Status:  VerdictFail,
			})
			out.Verdict = VerdictFail
			out.Reason = "zkverifier merkle root differs from summary"
		}
		return out
	}

	out.Verdict = VerdictPartial
	out.Coverage = 0.85
	out.Reason = trimOutput(outBytes)
	return out
}

func parseHTXVerifierRoot(text string) string {
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "account merkle tree root is ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "account merkle tree root is "))
		}
	}
	return ""
}
