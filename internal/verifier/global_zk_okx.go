package verifier

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// GlobalZKProofOKX validates OKX liability zk-STARK bundle (sum_proof_data.json).
type GlobalZKProofOKX struct {
	SnapshotID      string
	SummarySha256   string
	ProofBundlePath string
	ProofBundleSha  string
	MerkleRoot      string
}

type okxSumProof struct {
	General      map[string]any `json:"general"`
	RootVDDigest map[string]any `json:"root_vd_digest"`
	CircuitsInfo any            `json:"circuits_info"`
	Proof        okxProofEnvelope `json:"proof"`
}

func (v GlobalZKProofOKX) Run() Verification {
	out := Verification{
		VerifierID:     "global-zk-proof",
		Version:        "okx-1",
		SnapshotID:     v.SnapshotID,
		VerifiedAt:     time.Now().UTC(),
		InputArtifacts: artifactInputs(v.SummarySha256, v.ProofBundleSha),
	}

	if v.ProofBundlePath == "" {
		out.Verdict = VerdictUnverifiable
		out.Reason = "liability zk-STARK zip not present in artifact bundle"
		return out
	}

	proofJSON, err := readSumProofData(v.ProofBundlePath)
	if err != nil {
		out.Verdict = VerdictFail
		out.Reason = fmt.Sprintf("read sum_proof_data.json: %v", err)
		return out
	}
	var parsed okxSumProof
	if err := json.Unmarshal(proofJSON, &parsed); err != nil {
		out.Verdict = VerdictFail
		out.Reason = fmt.Sprintf("decode sum_proof_data.json: %v", err)
		return out
	}
	if len(parsed.Proof.PublicInputs) == 0 || parsed.Proof.Proof == nil {
		out.Verdict = VerdictFail
		out.Reason = "sum_proof_data.json missing proof or public_inputs"
		return out
	}
	if parsed.RootVDDigest == nil {
		out.Verdict = VerdictFail
		out.Reason = "sum_proof_data.json missing root_vd_digest"
		return out
	}

	out.Findings = append(out.Findings, Finding{
		Subject: "sum_proof_data.json",
		Field:   "structure",
		Status:  VerdictPass,
		Note:    "required zk-STARK fields present",
	})

	agg, err := okxZKAggregatesFromProof(parsed.Proof)
	if err != nil {
		out.Verdict = VerdictFail
		out.Reason = fmt.Sprintf("parse zk public inputs: %v", err)
		return out
	}
	out.Findings = append(out.Findings,
		Finding{
			Subject: "zk-liability",
			Field:   "aggregate_equity",
			Actual:  fmt.Sprintf("%d", agg.Equity),
			Status:  VerdictPass,
			Note:    "all-token aggregate from zk public inputs (internal units)",
		},
		Finding{
			Subject: "zk-liability",
			Field:   "aggregate_debt",
			Actual:  fmt.Sprintf("%d", agg.Debt),
			Status:  VerdictPass,
		},
		Finding{
			Subject: "zk-liability",
			Field:   "aggregate_liability",
			Actual:  fmt.Sprintf("%d", agg.Liability),
			Status:  VerdictPass,
			Note:    "equity - debt; not directly comparable to per-coin summary liabilityBalances",
		},
	)

	if v.MerkleRoot != "" {
		if merkleHexEqual(v.MerkleRoot, agg.MerkleHex) {
			out.Findings = append(out.Findings, Finding{
				Subject: "merkleRoot",
				Field:   "summary_vs_zk_proof",
				Claim:   normalizeMerkleHex(v.MerkleRoot),
				Actual:  agg.MerkleHex,
				Status:  VerdictPass,
				Note:    "summary merkleHash matches zk proof public root hash",
			})
		} else {
			out.Findings = append(out.Findings, Finding{
				Subject: "merkleRoot",
				Field:   "summary_vs_zk_proof",
				Claim:   normalizeMerkleHex(v.MerkleRoot),
				Actual:  agg.MerkleHex,
				Status:  VerdictFail,
			})
			out.Verdict = VerdictFail
			out.Reason = "summary merkleHash does not match zk proof root"
			return out
		}
	}

	if roundStr := okxProofRoundString(parsed.General); roundStr != "" && v.SnapshotID != "" {
		if roundStr == v.SnapshotID {
			out.Findings = append(out.Findings, Finding{
				Subject: "round_num",
				Field:   "snapshot_id",
				Claim:   v.SnapshotID,
				Actual:  roundStr,
				Status:  VerdictPass,
			})
		} else {
			out.Findings = append(out.Findings, Finding{
				Subject: "round_num",
				Field:   "snapshot_id",
				Claim:   v.SnapshotID,
				Actual:  roundStr,
				Status:  VerdictFail,
			})
			out.Verdict = VerdictFail
			out.Reason = fmt.Sprintf("proof round_num %s != snapshot %s", roundStr, v.SnapshotID)
			return out
		}
	}

	bin := os.Getenv("OKX_ZK_STARK_VALIDATOR")
	if bin == "" {
		out.Verdict = VerdictPartial
		out.Coverage = 0.75
		out.Reason = "structure and summary binding OK; set OKX_ZK_STARK_VALIDATOR for cryptographic verification"
		return out
	}

	dir, err := extractProofDir(v.ProofBundlePath)
	if err != nil {
		out.Verdict = VerdictFail
		out.Reason = fmt.Sprintf("extract proof zip: %v", err)
		return out
	}
	defer os.RemoveAll(dir)

	proofPath := filepath.Join(dir, "sum_proof_data.json")
	if _, err := os.Stat(proofPath); err != nil {
		out.Verdict = VerdictFail
		out.Reason = fmt.Sprintf("sum_proof_data.json missing after extract: %v", err)
		return out
	}

	cmd := exec.Command(bin, "verify-global", "--proof-path", proofPath)
	cmd.Stdin = strings.NewReader("\n")
	outBytes, err := cmd.CombinedOutput()
	if err != nil {
		out.Verdict = VerdictFail
		out.Reason = fmt.Sprintf("zkSTARKValidator: %v (%s)", err, trimOutput(outBytes))
		return out
	}
	text := string(outBytes)
	if strings.Contains(text, "successfully verify the global proof") || strings.Contains(text, "Execution result: Ok(())") {
		out.Verdict = VerdictPass
		out.Coverage = 1.0
		out.Reason = "zkSTARKValidator verify-global succeeded; summary merkle root bound to zk proof"
		return out
	}
	out.Verdict = VerdictPartial
	out.Coverage = 0.85
	out.Reason = trimOutput(outBytes)
	return out
}

func okxProofRoundString(general map[string]any) string {
	v, ok := general["round_num"]
	if !ok || v == nil {
		return ""
	}
	switch n := v.(type) {
	case float64:
		return strconv.FormatInt(int64(n), 10)
	case json.Number:
		return n.String()
	case string:
		return strings.TrimSpace(n)
	case int:
		return strconv.Itoa(n)
	case int64:
		return strconv.FormatInt(n, 10)
	default:
		return fmt.Sprint(v)
	}
}

func readSumProofData(zipPath string) ([]byte, error) {
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	for _, f := range zr.File {
		if strings.EqualFold(filepath.Base(f.Name), "sum_proof_data.json") {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}
	return nil, fmt.Errorf("sum_proof_data.json not found in %s", zipPath)
}

func extractProofDir(zipPath string) (string, error) {
	dir, err := os.MkdirTemp("", "okx-zk-*")
	if err != nil {
		return "", err
	}
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		os.RemoveAll(dir)
		return "", err
	}
	defer zr.Close()
	for _, f := range zr.File {
		dst := filepath.Join(dir, filepath.Base(f.Name))
		rc, err := f.Open()
		if err != nil {
			os.RemoveAll(dir)
			return "", err
		}
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			os.RemoveAll(dir)
			return "", err
		}
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			os.RemoveAll(dir)
			return "", err
		}
	}
	return dir, nil
}

func artifactInputs(parts ...string) []string {
	var out []string
	for _, p := range parts {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func trimOutput(b []byte) string {
	s := strings.TrimSpace(string(b))
	if len(s) > 400 {
		return s[:400] + "..."
	}
	return s
}
