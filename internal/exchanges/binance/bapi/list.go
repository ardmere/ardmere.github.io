package bapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

var monthTokens = []string{"JAN", "FEB", "MAR", "APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC"}

// ListEntry is one row from auditProofSnapshotCondition.
type ListEntry struct {
	SnapshotTime     string
	SnapshotTimeUTC  time.Time
	BTCBlockHeight   uint32
	AuditID          string
}

// FetchRecentListEntries returns the newest N snapshot list rows with derived audit ids.
func FetchRecentListEntries(ctx context.Context, limit int) ([]ListEntry, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlSnapshotList, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("snapshot list HTTP %d", resp.StatusCode)
	}
	var lr struct {
		Code string   `json:"code"`
		Data []string `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&lr); err != nil {
		return nil, err
	}
	if lr.Code != "000000" || len(lr.Data) == 0 {
		return nil, fmt.Errorf("snapshot list empty")
	}
	if limit <= 0 || limit > len(lr.Data) {
		limit = len(lr.Data)
	}
	out := make([]ListEntry, 0, limit)
	for i, entry := range lr.Data {
		if i >= limit {
			break
		}
		m := snapshotListEntry.FindStringSubmatch(entry)
		if m == nil {
			continue
		}
		t, err := ParseSnapshotTime(m[1])
		if err != nil {
			continue
		}
		h, _ := strconv.ParseUint(m[2], 10, 32)
		out = append(out, ListEntry{
			SnapshotTime:    m[1],
			SnapshotTimeUTC: t,
			BTCBlockHeight:  uint32(h),
			AuditID:         AuditIDFromSnapshotTime(m[1]),
		})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no parseable snapshot list entries")
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].SnapshotTimeUTC.After(out[j].SnapshotTimeUTC)
	})
	return out, nil
}

// AuditIDFromSnapshotTime maps Binance list time to audit id (DD/MM/YY in list strings).
// Example: "01/04/26 00:00:00 UTC" -> PR01APR26
func AuditIDFromSnapshotTime(snapshotTime string) string {
	parts := strings.Fields(snapshotTime)
	if len(parts) == 0 {
		return ""
	}
	dmy := strings.Split(parts[0], "/")
	if len(dmy) != 3 {
		return ""
	}
	dd, mm, yy := dmy[0], dmy[1], dmy[2]
	m, _ := strconv.Atoi(mm)
	if m < 1 || m > 12 {
		return ""
	}
	return fmt.Sprintf("PR%s%s%s", dd, monthTokens[m-1], yy)
}
