package sitegen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// CustodyRiskIndex maps exchange/snapshot id to file-only custody risk score.
type CustodyRiskIndex struct {
	byKey map[string]custodyRiskRow
}

type custodyRiskScoresFile struct {
	Scores []custodyRiskRow `json:"scores"`
}

type custodyRiskRow struct {
	Exchange string `json:"exchange"`
	Snapshot string `json:"snapshot"`
	RiskScore *int  `json:"riskScore"`
	RiskLevel string `json:"riskLevel"`
}

func custodyRiskKey(exchange, snapshotID string) string {
	return exchange + "\x00" + snapshotID
}

// LoadCustodyRiskScores reads docs/reports/custody-risk-scores.json when present.
func LoadCustodyRiskScores(reportsDir string) (*CustodyRiskIndex, error) {
	path := "docs/reports/custody-risk-scores.json"
	if reportsDir != "" {
		path = filepath.Join(reportsDir, "custody-risk-scores.json")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &CustodyRiskIndex{byKey: map[string]custodyRiskRow{}}, nil
		}
		return nil, err
	}
	var doc custodyRiskScoresFile
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	idx := &CustodyRiskIndex{byKey: make(map[string]custodyRiskRow, len(doc.Scores))}
	for _, row := range doc.Scores {
		if row.Exchange == "" || row.Snapshot == "" {
			continue
		}
		idx.byKey[custodyRiskKey(row.Exchange, row.Snapshot)] = row
	}
	return idx, nil
}

func custodyRiskClass(level string) string {
	switch level {
	case "critical":
		return "custody-risk custody-risk-critical"
	case "high":
		return "custody-risk custody-risk-high"
	case "medium":
		return "custody-risk custody-risk-medium"
	case "low":
		return "custody-risk custody-risk-low"
	default:
		return "custody-risk"
	}
}

func applyCustodyRisk(rows []ReportRow, assessments []Assessment, custody *CustodyRiskIndex) {
	if custody == nil {
		return
	}
	for i := range rows {
		if i >= len(assessments) {
			break
		}
		a := assessments[i]
		if score, level, ok := custody.Lookup(a.Exchange, a.SnapshotID); ok {
			rows[i].CustodyRisk = score
			rows[i].CustodyLevel = level
		}
	}
}

func (c *CustodyRiskIndex) Lookup(exchange, snapshotID string) (score string, level string, ok bool) {
	if c == nil || c.byKey == nil {
		return "", "", false
	}
	row, ok := c.byKey[custodyRiskKey(exchange, snapshotID)]
	if !ok || row.RiskScore == nil {
		return "", "", false
	}
	return fmt.Sprintf("%d", *row.RiskScore), row.RiskLevel, true
}
