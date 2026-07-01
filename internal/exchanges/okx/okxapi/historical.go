package okxapi

import (
	"archive/zip"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// FetchSummaryByAuditID returns summary JSON for a specific OKX audit id.
// Uses the detail page when it matches; otherwise builds from audit list + liability zip merkle root.
func FetchSummaryByAuditID(ctx context.Context, auditID string) (SummaryBundle, []byte, error) {
	auditID = strings.TrimSpace(auditID)
	if auditID == "" {
		return SummaryBundle{}, nil, fmt.Errorf("audit id required")
	}

	current, raw, err := FetchCurrentSummary(ctx)
	if err == nil && current.Audit.AuditID == auditID {
		return current, raw, nil
	}

	rec, err := findAuditRecord(ctx, auditID)
	if err != nil {
		return SummaryBundle{}, nil, err
	}

	tmp, err := os.MkdirTemp("", "okx-liab-*")
	if err != nil {
		return SummaryBundle{}, nil, err
	}
	defer os.RemoveAll(tmp)

	zipPath, _, _, err := DownloadFile(ctx, rec.MerkleTreeDownload, tmp, ".zip")
	if err != nil {
		return SummaryBundle{}, nil, fmt.Errorf("download liability zip: %w", err)
	}
	merkle, err := MerkleHashFromLiabilityZip(zipPath)
	if err != nil {
		return SummaryBundle{}, nil, err
	}

	bundle := SummaryBundle{
		FetchedAt: time.Now().UTC().Format(time.RFC3339),
		Source:    "audit-list+liability-zip",
		DetailURL: URLDownload,
		Audit: AuditRootInfo{
			AuditID:    auditID,
			MerkleHash: merkle,
			CreateTime: rec.CreateTime,
			TypeNum:    rec.TypeNum,
		},
	}
	out, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return SummaryBundle{}, nil, err
	}
	return bundle, out, nil
}

func findAuditRecord(ctx context.Context, auditID string) (AuditRecord, error) {
	list, err := FetchAuditList(ctx)
	if err != nil {
		return AuditRecord{}, err
	}
	for _, rec := range list {
		if rec.AuditID == auditID {
			return rec, nil
		}
	}
	return AuditRecord{}, fmt.Errorf("audit %s not found in download list", auditID)
}

// MerkleHashFromLiabilityZip reads sum_proof_data.json inside the liability zip.
func MerkleHashFromLiabilityZip(zipPath string) (string, error) {
	raw, err := readSumProofData(zipPath)
	if err != nil {
		return "", err
	}
	var doc struct {
		Proof struct {
			PublicInputs []json.Number `json:"public_inputs"`
		} `json:"proof"`
		General map[string]any `json:"general"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return "", fmt.Errorf("decode sum_proof_data.json: %w", err)
	}
	if len(doc.Proof.PublicInputs) < 6 {
		return "", fmt.Errorf("sum_proof_data.json: expected >=6 public_inputs, got %d", len(doc.Proof.PublicInputs))
	}
	merkle, err := merkleRootHexFromLimbs(doc.Proof.PublicInputs[2:6])
	if err != nil {
		return "", err
	}
	if round := proofRoundString(doc.General); round != "" {
		_ = round
	}
	return merkle, nil
}

func merkleRootHexFromLimbs(limbs []json.Number) (string, error) {
	if len(limbs) != 4 {
		return "", fmt.Errorf("merkle root: need 4 limbs, got %d", len(limbs))
	}
	var buf [32]byte
	for i, limb := range limbs {
		u, err := strconv.ParseUint(limb.String(), 10, 64)
		if err != nil {
			return "", fmt.Errorf("merkle limb %d: %w", i, err)
		}
		binary.LittleEndian.PutUint64(buf[i*8:], u)
	}
	return fmt.Sprintf("%x", buf[:]), nil
}

func proofRoundString(general map[string]any) string {
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
