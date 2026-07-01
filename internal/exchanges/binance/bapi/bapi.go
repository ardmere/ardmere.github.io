// Package bapi is a thin client over Binance's public Proof-of-Reserves
// BAPI endpoints. See docs/binance-por-data-guide.md.
package bapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	urlSnapshot     = "https://www.binance.com/bapi/apex/v1/public/apex/market/query/userReserveAuditProofSnapshot"
	urlSnapshotList = "https://www.binance.com/bapi/apex/v1/public/apex/market/query/auditProofSnapshotCondition"
	urlDownload     = "https://www.binance.com/bapi/apex/v1/public/apex/market/por/getDownloadUrl?auditId="
	userAgent       = "ardmere/0.1 (+https://ardmere.org)"
)

// CoinRow mirrors one entry of snapshotDataList.
type CoinRow struct {
	Coin                  string   `json:"coin"`
	Ratio                 string   `json:"ratio"`
	CustomerLiability     float64  `json:"customerLiability"`
	BinanceLiability      float64  `json:"binanceLiability"`
	ExchangeBalance       float64  `json:"exchangeBalance"`
	ThirdPartyCustody     float64  `json:"thirdPartyCustody"`
	MarginInsurance       *float64 `json:"marginInsurance"`
	FutureInsurance       *float64 `json:"futureInsurance"`
	CustomerLiabilityUsdt float64  `json:"customerLiabilityUsdt"`
}

// Snapshot is the parsed payload of the userReserveAuditProofSnapshot endpoint.
type Snapshot struct {
	SnapshotTime     string    `json:"snapshotTime"`
	MerkleRootHash   string    `json:"merkleRootHash"`
	Auditor          string    `json:"auditor"`
	AuditorLink      string    `json:"auditorLink"`
	AuditID          string    `json:"auditId"`
	AuditDate        string    `json:"auditDate"`
	SnapshotDataList []CoinRow `json:"snapshotDataList"`
}

type snapResp struct {
	Code string   `json:"code"`
	Data Snapshot `json:"data"`
}

// SnapshotMeta holds derived metadata for on-chain anchoring.
type SnapshotMeta struct {
	PeriodSeq      uint32 // 1 = oldest snapshot in Binance history
	BTCBlockHeight uint32
	SnapshotTime   time.Time
}

var snapshotListEntry = regexp.MustCompile(`^(\d{2}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} UTC) \| BTC Block Height (\d+)$`)

// FetchSnapshotMeta looks up periodSeq and BTC block height for auditID
// using the public snapshot history list (newest-first).
func FetchSnapshotMeta(ctx context.Context, auditID string, snap Snapshot) (SnapshotMeta, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlSnapshotList, nil)
	if err != nil {
		return SnapshotMeta{}, err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := httpClient().Do(req)
	if err != nil {
		return SnapshotMeta{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return SnapshotMeta{}, fmt.Errorf("snapshot list HTTP %d", resp.StatusCode)
	}
	var lr struct {
		Code string   `json:"code"`
		Data []string `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&lr); err != nil {
		return SnapshotMeta{}, err
	}
	if lr.Code != "000000" || len(lr.Data) == 0 {
		return SnapshotMeta{}, fmt.Errorf("snapshot list empty")
	}

	meta := SnapshotMeta{}
	if t, err := ParseSnapshotTime(snap.SnapshotTime); err == nil {
		meta.SnapshotTime = t
	}

	for i, entry := range lr.Data {
		m := snapshotListEntry.FindStringSubmatch(entry)
		if m == nil {
			continue
		}
		if m[1] != snap.SnapshotTime {
			continue
		}
		h, _ := strconv.ParseUint(m[2], 10, 32)
		meta.BTCBlockHeight = uint32(h)
		// list is newest-first; oldest snapshot = period 1
		meta.PeriodSeq = uint32(len(lr.Data) - i)
		return meta, nil
	}
	return meta, fmt.Errorf("auditId %s snapshotTime %q not found in history list", auditID, snap.SnapshotTime)
}

// ParseSnapshotTime parses Binance snapshot time, e.g. "01/04/26 00:00:00 UTC".
func ParseSnapshotTime(s string) (time.Time, error) {
	return time.Parse("01/02/06 15:04:05 MST", strings.TrimSpace(s))
}

type dlResp struct {
	Code string `json:"code"`
	Data string `json:"data"`
}

// LoadSnapshot reads a saved BAPI snapshot JSON from disk.
func LoadSnapshot(path string) (Snapshot, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, err
	}
	var sr snapResp
	if err := json.Unmarshal(raw, &sr); err != nil {
		var direct Snapshot
		if err2 := json.Unmarshal(raw, &direct); err2 != nil {
			return Snapshot{}, fmt.Errorf("decode BAPI snapshot %s: %w", path, err)
		}
		return direct, nil
	}
	if sr.Data.AuditID == "" {
		return Snapshot{}, fmt.Errorf("BAPI snapshot %s has no auditId", path)
	}
	return sr.Data, nil
}

// FetchSnapshot returns the current snapshot plus the raw response bytes
// (so callers can content-address the artifact byte-for-byte).
func FetchSnapshot(ctx context.Context) ([]byte, Snapshot, error) {
	body := []byte(`{"time":"","pageIndex":0,"pageSize":0}`)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlSnapshot, bytes.NewReader(body))
	if err != nil {
		return nil, Snapshot{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := httpClient().Do(req)
	if err != nil {
		return nil, Snapshot{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, Snapshot{}, fmt.Errorf("BAPI snapshot HTTP %d", resp.StatusCode)
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, Snapshot{}, err
	}
	var sr snapResp
	if err := json.Unmarshal(raw, &sr); err != nil {
		return nil, Snapshot{}, fmt.Errorf("decode BAPI snapshot: %w", err)
	}
	if sr.Code != "000000" || sr.Data.AuditID == "" {
		return nil, Snapshot{}, fmt.Errorf("BAPI snapshot returned no data: %s", string(raw))
	}
	return raw, sr.Data, nil
}

// FetchWalletZipURL resolves the wallet-address ZIP URL for a given auditId.
func FetchWalletZipURL(ctx context.Context, auditID string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlDownload+auditID, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := httpClient().Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("getDownloadUrl HTTP %d", resp.StatusCode)
	}
	var dr dlResp
	if err := json.NewDecoder(resp.Body).Decode(&dr); err != nil {
		return "", err
	}
	if dr.Code != "000000" || dr.Data == "" {
		return "", fmt.Errorf("getDownloadUrl returned no data for %s", auditID)
	}
	return dr.Data, nil
}

func httpClient() *http.Client { return &http.Client{Timeout: 5 * time.Minute} }

// SnapshotURL returns the public BAPI snapshot URL (constant).
func SnapshotURL() string { return urlSnapshot }
