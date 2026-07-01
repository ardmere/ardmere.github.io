package htxapi

import (
	"archive/zip"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	zkConfigName = "config.json"
	zkProofName  = "proof0.csv"
	zkVKName     = "zkpor500.vk.save"
)

// ZkConfig matches merkle_groth16 verifier config.json.
type ZkConfig struct {
	ProofTable    string `json:"ProofTable"`
	ZkKeyName     string `json:"ZkKeyName"`
	CexAssetsInfo []struct {
		TotalBalance int64  `json:"TotalBalance"`
		Symbol       string `json:"Symbol"`
		Index        int    `json:"Index"`
	} `json:"CexAssetsInfo"`
}

// ZkBundleMeta is metadata extracted from an HTX public-data zip.
type ZkBundleMeta struct {
	AuditID        string
	SnapshotTime   time.Time
	SnapshotTimeRaw string
	MerkleRoot     string
	BatchCount     int
	AssetCount     int
	Config         ZkConfig
}

// ParseZkBundleZip reads config.json and proof0.csv from a public-data.zip.
func ParseZkBundleZip(zipPath string) (ZkBundleMeta, error) {
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return ZkBundleMeta{}, err
	}
	defer zr.Close()

	cfgBytes, err := readZipBase(zr, zkConfigName)
	if err != nil {
		return ZkBundleMeta{}, err
	}
	var cfg ZkConfig
	if err := json.Unmarshal(cfgBytes, &cfg); err != nil {
		return ZkBundleMeta{}, fmt.Errorf("decode config.json: %w", err)
	}
	if cfg.ProofTable == "" {
		cfg.ProofTable = zkProofName
	}
	if cfg.ZkKeyName == "" {
		cfg.ZkKeyName = "zkpor500"
	}

	proofBytes, err := readZipBase(zr, zkProofName)
	if err != nil {
		return ZkBundleMeta{}, err
	}
	if _, err := readZipBase(zr, zkVKName); err != nil {
		return ZkBundleMeta{}, fmt.Errorf("missing %s in zk bundle", zkVKName)
	}

	root, batchCount, snapTime, err := merkleRootFromProofCSV(proofBytes)
	if err != nil {
		return ZkBundleMeta{}, err
	}

	auditID := snapTime.UTC().Format("20060102")
	if auditID == "00010101" {
		auditID = "zk-unknown"
	}

	return ZkBundleMeta{
		AuditID:         auditID,
		SnapshotTime:    snapTime,
		SnapshotTimeRaw: snapTime.UTC().Format(time.RFC3339),
		MerkleRoot:      root,
		BatchCount:      batchCount,
		AssetCount:      len(cfg.CexAssetsInfo),
		Config:          cfg,
	}, nil
}

// ParsedSummaryFromZk builds a minimal summary from zk bundle metadata.
func ParsedSummaryFromZk(meta ZkBundleMeta) ParsedSummary {
	out := ParsedSummary{
		AuditID:      meta.AuditID,
		AuditTime:    meta.SnapshotTime,
		AuditTimeRaw: meta.SnapshotTimeRaw,
		MerkleRoot:   meta.MerkleRoot,
	}
	for _, a := range meta.Config.CexAssetsInfo {
		if a.Symbol == "" {
			continue
		}
		out.CoinRows = append(out.CoinRows, CoinRow{
			Coin:          strings.ToUpper(a.Symbol),
			ReserveAmount: float64(a.TotalBalance),
			// Internal integer units from config; not comparable to on-chain without HTX scaling docs.
		})
	}
	sort.Slice(out.CoinRows, func(i, j int) bool { return out.CoinRows[i].Coin < out.CoinRows[j].Coin })
	return out
}

// SummaryBytesFromZk builds an archivable summary JSON from zk-derived fields.
func SummaryBytesFromZk(meta ZkBundleMeta) ([]byte, error) {
	info := map[string]any{
		"audit_id":         meta.AuditID,
		"audit_time":       meta.SnapshotTimeRaw,
		"merkle_root_hash": meta.MerkleRoot,
		"batch_count":      meta.BatchCount,
		"asset_count":      meta.AssetCount,
		"algorithm":        "groth16+zkpor500",
		"source":           "zk-bundle-derived",
	}
	infoRaw, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}
	return BuildSummaryBundle(infoRaw, nil, "zk-bundle", true)
}

func readZipBase(zr *zip.ReadCloser, base string) ([]byte, error) {
	for _, f := range zr.File {
		if filepath.Base(f.Name) == base {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}
	return nil, fmt.Errorf("%s not found in zip", base)
}

func merkleRootFromProofCSV(data []byte) (rootHex string, batchCount int, snapTime time.Time, err error) {
	r := csv.NewReader(strings.NewReader(string(data)))
	records, err := r.ReadAll()
	if err != nil {
		return "", 0, time.Time{}, fmt.Errorf("read proof csv: %w", err)
	}
	if len(records) < 2 {
		return "", 0, time.Time{}, fmt.Errorf("proof csv too short")
	}
	header := records[0]
	col := map[string]int{}
	for i, h := range header {
		col[h] = i
	}
	createdIdx, okCreated := col["CreatedAt"]
	rootsIdx, okRoots := col["AccountTreeRoots"]
	batchIdx, okBatch := col["BatchNumber"]
	if !okCreated || !okRoots {
		return "", 0, time.Time{}, fmt.Errorf("proof csv missing CreatedAt or AccountTreeRoots columns")
	}

	var lastRoots string
	maxBatch := -1
	for _, row := range records[1:] {
		if len(row) <= rootsIdx {
			continue
		}
		batchNum := -1
		if okBatch && len(row) > batchIdx {
			if n, e := parseInt(row[batchIdx]); e == nil {
				batchNum = n
			}
		}
		if batchNum >= maxBatch {
			maxBatch = batchNum
			lastRoots = row[rootsIdx]
			batchCount = batchNum + 1
		}
		if len(row) > createdIdx {
			if t, e := time.Parse(time.RFC3339Nano, row[createdIdx]); e == nil {
				snapTime = t
			}
		}
	}
	rootHex, err = finalAccountTreeRoot(lastRoots)
	if err != nil {
		return "", 0, time.Time{}, err
	}
	return rootHex, batchCount, snapTime, nil
}

func finalAccountTreeRoot(accountTreeRootsField string) (string, error) {
	// CSV column is JSON array of two base64 roots; final merkle root is element [1].
	var roots []string
	if err := json.Unmarshal([]byte(accountTreeRootsField), &roots); err != nil {
		return "", fmt.Errorf("decode AccountTreeRoots: %w", err)
	}
	if len(roots) < 2 {
		return "", fmt.Errorf("AccountTreeRoots need 2 elements, got %d", len(roots))
	}
	raw, err := decodeB64Flexible(roots[1])
	if err != nil {
		return "", err
	}
	if len(raw) != 32 {
		return "", fmt.Errorf("account tree root length %d", len(raw))
	}
	return fmt.Sprintf("%x", raw), nil
}

func decodeB64Flexible(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("empty base64")
	}
	if b, err := base64.StdEncoding.DecodeString(s); err == nil {
		return b, nil
	}
	if b, err := base64.RawStdEncoding.DecodeString(s); err == nil {
		return b, nil
	}
	return nil, fmt.Errorf("invalid base64 root")
}

func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(strings.TrimSpace(s), "%d", &n)
	return n, err
}

// ExtractZkBundleToDir unpacks zip into dir/public-data for HTX zkverifiermac.
func ExtractZkBundleToDir(zipPath, workDir string) (string, error) {
	outDir := filepath.Join(workDir, "public-data")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", err
	}
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer zr.Close()
	for _, name := range []string{zkConfigName, zkProofName, zkVKName} {
		data, err := readZipBase(zr, name)
		if err != nil {
			return "", err
		}
		if err := os.WriteFile(filepath.Join(outDir, name), data, 0o644); err != nil {
			return "", err
		}
	}
	return outDir, nil
}
