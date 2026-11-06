// Package krakenapi fetches and parses Kraken Proof-of-Reserves summary artifacts.
// See docs/kraken-por-data-guide.md for data availability.
package krakenapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	userAgent = "ardmere/0.1 (+https://ardmere.org)"

	URLPoRPage          = "https://www.kraken.com/proof-of-reserves"
	URLLedgerLensAPI    = "https://widget.pro.ledgerlens.io/api/kraken-platform-spot-futures"
	URLLedgerLensWidget = "https://widget.pro.ledgerlens.io/kraken-platform-spot-futures"

	// Display floor used when Kraken CMS sets capPercentage=true (UI shows "105%+").
	cappedReserveFloor = 105.0
)

// Preferred coin order for display (matches Kraken PoR cards / LedgerLens liability order).
var preferredCoinOrder = []string{"BTC", "ETH", "SOL", "USDT", "USDC", "XRP", "ADA", "USDG"}

var (
	compactDateRe = regexp.MustCompile(`^\d{8}$`)
	initialPropsRE = regexp.MustCompile(`(?s)window\.__INITIAL_PROPS__=(\{.*?\})\s*;?\s*</script>`)
	auditExprRE    = regexp.MustCompile(`auditPeriod\s*==\s*"([^"]+)"`)
)

// SummaryBundle is the archived summarySnapshot for one Kraken PoR period.
type SummaryBundle struct {
	FetchedAt  time.Time       `json:"fetchedAt"`
	Source     string          `json:"source,omitempty"`
	SourceURLs []string        `json:"sourceUrls,omitempty"`
	ReportDate string          `json:"reportDate,omitempty"`
	SnapshotID string          `json:"snapshotId,omitempty"`
	LedgerLens json.RawMessage `json:"ledgerlens"`
	CMS        json.RawMessage `json:"cms"`
}

// ParsedSummary is normalized fields for por.Snapshot.
// Reserve rates are percentages (100 = 100%) for solvency-claim.
type ParsedSummary struct {
	AuditID          string
	ReportDate       string
	AuditTime        time.Time
	AuditTimeRaw     string
	MerkleRoot       string
	TotalReserveRate float64
	Accountant       string
	DatasetID        string
	AuditPeriodID    string
	CoinRows         []CoinRow
}

type CoinRow struct {
	Coin           string
	ReserveRate    float64 // effective rate used by solvency-claim (≥105 when capped)
	CapPercentage  bool
	DisplayRatio   string
	HideBalances   bool
}

type ledgerLensAPI struct {
	Slug      string              `json:"slug"`
	Datasets  []ledgerLensDataset `json:"datasets"`
	Organization struct {
		Name string `json:"name"`
	} `json:"organization"`
}

type ledgerLensDataset struct {
	DatasetID               string `json:"datasetId"`
	DatasetTimestamp        string `json:"datasetTimestamp"`
	DatasetDisplayTimestamp string `json:"datasetDisplayTimestamp"`
	HasReport               bool   `json:"hasReport"`
	MerkleTreeData          struct {
		CombinedRoot     string `json:"combinedRoot"`
		TotalLayers      int    `json:"totalLayers"`
		TotalLeaves      int64  `json:"totalLeaves"`
		HashingStrategy  string `json:"hashingStrategy"`
		DatasetFormat    string `json:"datasetFormat"`
		SortBalanced     bool   `json:"sortBalancedEntries"`
	} `json:"merkleTreeData"`
	Accountant struct {
		Name string `json:"name"`
	} `json:"accountant"`
	CoveredLiabilities []struct {
		Symbol string `json:"symbol"`
	} `json:"coveredLiabilities"`
	AssetTokens []struct {
		Symbol      string `json:"symbol"`
		DisplayText string `json:"displayText"`
	} `json:"assetTokens"`
}

type cmsCoinRow struct {
	Coin              string  `json:"coin"`
	ReserveRatioValue float64 `json:"reserveRatioValue"`
	CapPercentage     bool    `json:"capPercentage"`
	HideBalances      bool    `json:"hideBalances"`
}

type cmsBlock struct {
	AuditPeriodID    string       `json:"auditPeriodId"`
	AuditPeriodLabel string       `json:"auditPeriodLabel"`
	CoinRows         []cmsCoinRow `json:"coinRows"`
}

type ledgerLensBlock struct {
	Slug                    string          `json:"slug"`
	DatasetID               string          `json:"datasetId"`
	DatasetTimestamp        string          `json:"datasetTimestamp"`
	DatasetDisplayTimestamp string          `json:"datasetDisplayTimestamp"`
	Accountant              string          `json:"accountant"`
	MerkleTreeData          json.RawMessage `json:"merkleTreeData"`
	CoveredLiabilities      json.RawMessage `json:"coveredLiabilities,omitempty"`
	AssetTokens             json.RawMessage `json:"assetTokens,omitempty"`
	HasReport               bool            `json:"hasReport"`
}

// SummaryURL returns the primary machine-readable source used for artifact provenance.
func SummaryURL() string { return URLLedgerLensAPI }

// FetchPublicSummary downloads the latest quarterly Kraken PoR summary
// (LedgerLens Merkle root + CMS reserve ratios).
func FetchPublicSummary(ctx context.Context) ([]byte, ParsedSummary, error) {
	return FetchPublicSummaryByAuditID(ctx, "")
}

// FetchPublicSummaryByAuditID selects a period (e.g. "20260630", "2026-06-30", "30june").
// Empty auditID picks the newest LedgerLens dataset.
func FetchPublicSummaryByAuditID(ctx context.Context, auditID string) ([]byte, ParsedSummary, error) {
	llRaw, ll, err := fetchLedgerLens(ctx)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	if len(ll.Datasets) == 0 {
		return nil, ParsedSummary{}, fmt.Errorf("kraken ledgerlens: no datasets")
	}

	cmsRaw, periods, err := fetchCMSPeriods(ctx)
	if err != nil {
		return nil, ParsedSummary{}, err
	}

	ds, period, err := selectPeriod(ll.Datasets, periods, auditID)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	return buildBundle(ll, ds, period, llRaw, cmsRaw)
}

// ParseSummaryBytes parses an archived SummaryBundle.
func ParseSummaryBytes(raw []byte) (ParsedSummary, error) {
	var bundle SummaryBundle
	if err := json.Unmarshal(raw, &bundle); err != nil {
		return ParsedSummary{}, fmt.Errorf("decode kraken summary bundle: %w", err)
	}
	if len(bundle.LedgerLens) == 0 {
		return ParsedSummary{}, fmt.Errorf("kraken summary missing ledgerlens block")
	}
	var ll ledgerLensBlock
	if err := json.Unmarshal(bundle.LedgerLens, &ll); err != nil {
		return ParsedSummary{}, fmt.Errorf("decode ledgerlens block: %w", err)
	}
	var cms cmsBlock
	if len(bundle.CMS) > 0 {
		if err := json.Unmarshal(bundle.CMS, &cms); err != nil {
			return ParsedSummary{}, fmt.Errorf("decode cms block: %w", err)
		}
	}
	return parseBlocks(bundle.SnapshotID, bundle.ReportDate, ll, cms)
}

func fetchLedgerLens(ctx context.Context) ([]byte, ledgerLensAPI, error) {
	raw, err := getBytes(ctx, URLLedgerLensAPI, "application/json")
	if err != nil {
		return nil, ledgerLensAPI{}, err
	}
	var ll ledgerLensAPI
	if err := json.Unmarshal(raw, &ll); err != nil {
		return nil, ledgerLensAPI{}, fmt.Errorf("decode ledgerlens api: %w", err)
	}
	return raw, ll, nil
}

func fetchCMSPeriods(ctx context.Context) ([]byte, []cmsPeriod, error) {
	raw, err := getBytes(ctx, URLPoRPage, "text/html,application/xhtml+xml")
	if err != nil {
		return nil, nil, err
	}
	periods, err := parseCMSPeriods(raw)
	if err != nil {
		return nil, nil, err
	}
	return raw, periods, nil
}

type cmsPeriod struct {
	ID    string
	Label string
	Date  time.Time
	Rows  []cmsCoinRow
}

func parseCMSPeriods(html []byte) ([]cmsPeriod, error) {
	m := initialPropsRE.FindSubmatch(html)
	if m == nil {
		return nil, fmt.Errorf("kraken por page missing window.__INITIAL_PROPS__; capture JSON and use -summary-path")
	}
	var props map[string]json.RawMessage
	if err := json.Unmarshal(m[1], &props); err != nil {
		return nil, fmt.Errorf("decode __INITIAL_PROPS__: %w", err)
	}
	var data map[string]json.RawMessage
	if err := json.Unmarshal(props["data"], &data); err != nil {
		return nil, fmt.Errorf("decode cms data: %w", err)
	}
	var components []json.RawMessage
	if err := json.Unmarshal(data["field_page_components"], &components); err != nil {
		return nil, fmt.Errorf("decode field_page_components: %w", err)
	}
	if len(components) < 6 {
		return nil, fmt.Errorf("kraken cms: unexpected page component count %d", len(components))
	}

	options, err := extractAuditOptions(components[3])
	if err != nil {
		return nil, err
	}
	groups, err := extractPeriodGroups(components[5])
	if err != nil {
		return nil, err
	}

	out := make([]cmsPeriod, 0, len(options))
	for _, opt := range options {
		date, err := parseAuditLabel(opt.Label)
		if err != nil {
			return nil, fmt.Errorf("parse audit label %q: %w", opt.Label, err)
		}
		rows := groups[opt.ID]
		if len(rows) == 0 {
			return nil, fmt.Errorf("kraken cms: no coin cards for auditPeriod %q", opt.ID)
		}
		out = append(out, cmsPeriod{ID: opt.ID, Label: opt.Label, Date: date, Rows: rows})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("kraken cms: no audit periods")
	}
	return out, nil
}

type auditOption struct {
	ID    string
	Label string
}

func extractAuditOptions(comp json.RawMessage) ([]auditOption, error) {
	var node map[string]json.RawMessage
	if err := json.Unmarshal(comp, &node); err != nil {
		return nil, err
	}
	// components[3].field_children[0].field_children[1].field_options
	children0, err := firstChild(node, "field_children")
	if err != nil {
		return nil, err
	}
	var wrap map[string]json.RawMessage
	if err := json.Unmarshal(children0, &wrap); err != nil {
		return nil, err
	}
	dropdown, err := nthChild(wrap, "field_children", 1)
	if err != nil {
		return nil, err
	}
	var dd map[string]json.RawMessage
	if err := json.Unmarshal(dropdown, &dd); err != nil {
		return nil, err
	}
	var opts []struct {
		ID    string `json:"field_option_id"`
		Label string `json:"field_option_label"`
	}
	if err := json.Unmarshal(dd["field_options"], &opts); err != nil {
		return nil, fmt.Errorf("decode audit options: %w", err)
	}
	out := make([]auditOption, 0, len(opts))
	for _, o := range opts {
		if strings.TrimSpace(o.ID) == "" {
			continue
		}
		out = append(out, auditOption{ID: o.ID, Label: o.Label})
	}
	return out, nil
}

func extractPeriodGroups(comp json.RawMessage) (map[string][]cmsCoinRow, error) {
	var node map[string]json.RawMessage
	if err := json.Unmarshal(comp, &node); err != nil {
		return nil, err
	}
	// components[5].field_children[1].field_children[*] (period groups)
	wrapRaw, err := nthChild(node, "field_children", 1)
	if err != nil {
		return nil, err
	}
	var wrap map[string]json.RawMessage
	if err := json.Unmarshal(wrapRaw, &wrap); err != nil {
		return nil, err
	}
	var groups []json.RawMessage
	if err := json.Unmarshal(wrap["field_children"], &groups); err != nil {
		return nil, fmt.Errorf("decode period groups: %w", err)
	}
	out := map[string][]cmsCoinRow{}
	for _, g := range groups {
		var gm map[string]json.RawMessage
		if err := json.Unmarshal(g, &gm); err != nil {
			return nil, err
		}
		var expr string
		_ = json.Unmarshal(gm["field_expression"], &expr)
		m := auditExprRE.FindStringSubmatch(expr)
		if m == nil {
			continue
		}
		periodID := m[1]
		var cards []json.RawMessage
		if err := json.Unmarshal(gm["field_children"], &cards); err != nil {
			return nil, err
		}
		rows := make([]cmsCoinRow, 0, len(cards))
		for _, card := range cards {
			row, err := extractCoinCard(card)
			if err != nil {
				return nil, err
			}
			if row.Coin == "" {
				continue
			}
			rows = append(rows, row)
		}
		out[periodID] = rows
	}
	return out, nil
}

func extractCoinCard(card json.RawMessage) (cmsCoinRow, error) {
	var node map[string]json.RawMessage
	if err := json.Unmarshal(card, &node); err != nil {
		return cmsCoinRow{}, err
	}
	var values []json.RawMessage
	if err := json.Unmarshal(node["field_values"], &values); err != nil {
		return cmsCoinRow{}, err
	}
	props := map[string]any{}
	for _, v := range values {
		absorbDatasetRecord(v, props)
	}
	row := cmsCoinRow{}
	if s, ok := props["assetCode"].(string); ok {
		row.Coin = strings.ToUpper(strings.TrimSpace(s))
	}
	switch r := props["reserveRatioValue"].(type) {
	case float64:
		row.ReserveRatioValue = r
	case json.Number:
		f, _ := r.Float64()
		row.ReserveRatioValue = f
	case string:
		f, _ := strconv.ParseFloat(r, 64)
		row.ReserveRatioValue = f
	}
	if b, ok := props["capPercentage"].(bool); ok {
		row.CapPercentage = b
	}
	if b, ok := props["hideBalances"].(bool); ok {
		row.HideBalances = b
	}
	return row, nil
}

func absorbDatasetRecord(raw json.RawMessage, props map[string]any) {
	var rec map[string]json.RawMessage
	if err := json.Unmarshal(raw, &rec); err != nil {
		return
	}
	if children, ok := rec["children"]; ok {
		var kids []json.RawMessage
		if err := json.Unmarshal(children, &kids); err == nil {
			for _, kid := range kids {
				absorbDatasetRecord(kid, props)
			}
		}
	}
	var name string
	if err := json.Unmarshal(rec["datasetPropertyName"], &name); err != nil || name == "" {
		return
	}
	var val any
	if err := json.Unmarshal(rec["fieldValue"], &val); err != nil {
		return
	}
	props[name] = val
}

func firstChild(node map[string]json.RawMessage, key string) (json.RawMessage, error) {
	return nthChild(node, key, 0)
}

func nthChild(node map[string]json.RawMessage, key string, n int) (json.RawMessage, error) {
	var kids []json.RawMessage
	if err := json.Unmarshal(node[key], &kids); err != nil {
		return nil, fmt.Errorf("decode %s: %w", key, err)
	}
	if n < 0 || n >= len(kids) {
		return nil, fmt.Errorf("%s[%d] missing (len=%d)", key, n, len(kids))
	}
	return kids[n], nil
}

func selectPeriod(datasets []ledgerLensDataset, periods []cmsPeriod, auditID string) (ledgerLensDataset, cmsPeriod, error) {
	auditID = strings.TrimSpace(auditID)
	if auditID == "" {
		ds := datasets[0]
		date, err := datasetReportDate(ds)
		if err != nil {
			return ledgerLensDataset{}, cmsPeriod{}, err
		}
		period, ok := findPeriodByDate(periods, date)
		if !ok {
			return ledgerLensDataset{}, cmsPeriod{}, fmt.Errorf("kraken cms: no period matching latest ledgerlens date %s", date.Format("2006-01-02"))
		}
		return ds, period, nil
	}

	// Match by CMS period id (e.g. 30june).
	for _, p := range periods {
		if strings.EqualFold(p.ID, auditID) {
			ds, ok := findDatasetByDate(datasets, p.Date)
			if !ok {
				return ledgerLensDataset{}, cmsPeriod{}, fmt.Errorf("kraken ledgerlens: no dataset for cms period %s (%s)", p.ID, p.Date.Format("2006-01-02"))
			}
			return ds, p, nil
		}
	}

	date, err := normalizeAuditDate(auditID)
	if err != nil {
		return ledgerLensDataset{}, cmsPeriod{}, err
	}
	period, ok := findPeriodByDate(periods, date)
	if !ok {
		return ledgerLensDataset{}, cmsPeriod{}, fmt.Errorf("kraken cms: no period for %s", date.Format("2006-01-02"))
	}
	ds, ok := findDatasetByDate(datasets, date)
	if !ok {
		return ledgerLensDataset{}, cmsPeriod{}, fmt.Errorf("kraken ledgerlens: no dataset for %s", date.Format("2006-01-02"))
	}
	return ds, period, nil
}

func findPeriodByDate(periods []cmsPeriod, date time.Time) (cmsPeriod, bool) {
	want := date.UTC().Format("2006-01-02")
	for _, p := range periods {
		if p.Date.UTC().Format("2006-01-02") == want {
			return p, true
		}
	}
	return cmsPeriod{}, false
}

func findDatasetByDate(datasets []ledgerLensDataset, date time.Time) (ledgerLensDataset, bool) {
	want := date.UTC().Format("2006-01-02")
	for _, ds := range datasets {
		d, err := datasetReportDate(ds)
		if err != nil {
			continue
		}
		if d.UTC().Format("2006-01-02") == want {
			return ds, true
		}
	}
	return ledgerLensDataset{}, false
}

func datasetReportDate(ds ledgerLensDataset) (time.Time, error) {
	display := strings.TrimSpace(ds.DatasetDisplayTimestamp)
	if display != "" {
		// e.g. "06/30/2026   11:59 PM UTC"
		fields := strings.Fields(display)
		if len(fields) >= 1 {
			if t, err := time.Parse("01/02/2006", fields[0]); err == nil {
				return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC), nil
			}
		}
	}
	if ds.DatasetTimestamp != "" {
		t, err := time.Parse(time.RFC3339Nano, ds.DatasetTimestamp)
		if err != nil {
			t, err = time.Parse(time.RFC3339, ds.DatasetTimestamp)
		}
		if err != nil {
			return time.Time{}, fmt.Errorf("parse datasetTimestamp %q: %w", ds.DatasetTimestamp, err)
		}
		// LedgerLens sometimes stores end-of-day US offset as next UTC day; prefer calendar day of display when available.
		return time.Date(t.UTC().Year(), t.UTC().Month(), t.UTC().Day(), 0, 0, 0, 0, time.UTC), nil
	}
	return time.Time{}, fmt.Errorf("dataset %s missing timestamp", ds.DatasetID)
}

func buildBundle(ll ledgerLensAPI, ds ledgerLensDataset, period cmsPeriod, llRaw, cmsHTML []byte) ([]byte, ParsedSummary, error) {
	_ = llRaw
	_ = cmsHTML
	merkleRaw, err := json.Marshal(ds.MerkleTreeData)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	liabRaw, _ := json.Marshal(ds.CoveredLiabilities)
	tokRaw, _ := json.Marshal(minifyTokens(ds.AssetTokens))
	accountant := strings.TrimSpace(ds.Accountant.Name)
	if accountant == "" {
		accountant = strings.TrimSpace(ll.Organization.Name)
	}
	llBlock := ledgerLensBlock{
		Slug:                    ll.Slug,
		DatasetID:               ds.DatasetID,
		DatasetTimestamp:        ds.DatasetTimestamp,
		DatasetDisplayTimestamp: ds.DatasetDisplayTimestamp,
		Accountant:              accountant,
		MerkleTreeData:          merkleRaw,
		CoveredLiabilities:      liabRaw,
		AssetTokens:             tokRaw,
		HasReport:               ds.HasReport,
	}
	cms := cmsBlock{
		AuditPeriodID:    period.ID,
		AuditPeriodLabel: period.Label,
		CoinRows:         period.Rows,
	}
	reportDate := period.Date.UTC().Format("2006-01-02")
	snapshotID := period.Date.UTC().Format("20060102")
	parsed, err := parseBlocks(snapshotID, reportDate, llBlock, cms)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	llJSON, err := json.Marshal(llBlock)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	cmsJSON, err := json.Marshal(cms)
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	bundle := SummaryBundle{
		FetchedAt:  time.Now().UTC(),
		Source:     URLLedgerLensAPI,
		SourceURLs: []string{URLPoRPage, URLLedgerLensAPI, URLLedgerLensWidget},
		ReportDate: reportDate,
		SnapshotID: snapshotID,
		LedgerLens: llJSON,
		CMS:        cmsJSON,
	}
	out, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return nil, ParsedSummary{}, err
	}
	return out, parsed, nil
}

func minifyTokens(in []struct {
	Symbol      string `json:"symbol"`
	DisplayText string `json:"displayText"`
}) []map[string]string {
	out := make([]map[string]string, 0, len(in))
	for _, t := range in {
		out = append(out, map[string]string{
			"symbol":      t.Symbol,
			"displayText": t.DisplayText,
		})
	}
	return out
}

func parseBlocks(snapshotID, reportDate string, ll ledgerLensBlock, cms cmsBlock) (ParsedSummary, error) {
	var merkle struct {
		CombinedRoot string `json:"combinedRoot"`
	}
	if err := json.Unmarshal(ll.MerkleTreeData, &merkle); err != nil {
		return ParsedSummary{}, fmt.Errorf("decode merkleTreeData: %w", err)
	}
	root := strings.TrimSpace(strings.ToLower(merkle.CombinedRoot))
	if root == "" {
		return ParsedSummary{}, fmt.Errorf("kraken summary missing combinedRoot")
	}
	if reportDate == "" {
		if cms.AuditPeriodLabel != "" {
			if t, err := parseAuditLabel(cms.AuditPeriodLabel); err == nil {
				reportDate = t.Format("2006-01-02")
			}
		}
	}
	if reportDate == "" {
		return ParsedSummary{}, fmt.Errorf("kraken summary missing reportDate")
	}
	ts, err := time.Parse("2006-01-02", reportDate)
	if err != nil {
		return ParsedSummary{}, fmt.Errorf("parse reportDate %q: %w", reportDate, err)
	}
	ts = ts.UTC()
	if snapshotID == "" {
		snapshotID = ts.Format("20060102")
	}

	rows := make([]CoinRow, 0, len(cms.CoinRows))
	seen := map[string]CoinRow{}
	minRate := math.MaxFloat64
	for _, c := range cms.CoinRows {
		coin := strings.TrimSpace(strings.ToUpper(c.Coin))
		if coin == "" {
			continue
		}
		rate := c.ReserveRatioValue
		display := formatRatio(rate, c.CapPercentage)
		if c.CapPercentage {
			// UI publishes "105%+" when exact ratio is capped; treat as floor 105 for Stage 0 solvency-claim.
			rate = cappedReserveFloor
		}
		row := CoinRow{
			Coin:          coin,
			ReserveRate:   rate,
			CapPercentage: c.CapPercentage,
			DisplayRatio:  display,
			HideBalances:  c.HideBalances,
		}
		seen[coin] = row
		if rate < minRate {
			minRate = rate
		}
	}
	if len(seen) == 0 {
		return ParsedSummary{}, fmt.Errorf("kraken summary missing cms coin rows")
	}
	for _, coin := range preferredCoinOrder {
		if row, ok := seen[coin]; ok {
			rows = append(rows, row)
			delete(seen, coin)
		}
	}
	extras := make([]string, 0, len(seen))
	for coin := range seen {
		extras = append(extras, coin)
	}
	sort.Strings(extras)
	for _, coin := range extras {
		rows = append(rows, seen[coin])
	}
	if minRate == math.MaxFloat64 {
		minRate = 0
	}
	accountant := strings.TrimSpace(ll.Accountant)
	if accountant == "" {
		accountant = "unpublished"
	}

	return ParsedSummary{
		AuditID:          snapshotID,
		ReportDate:       reportDate,
		AuditTime:        ts,
		AuditTimeRaw:     ts.Format(time.RFC3339),
		MerkleRoot:       root,
		TotalReserveRate: minRate,
		Accountant:       accountant,
		DatasetID:        ll.DatasetID,
		AuditPeriodID:    cms.AuditPeriodID,
		CoinRows:         rows,
	}, nil
}

func formatRatio(rate float64, capped bool) string {
	if capped {
		return "105%+"
	}
	if rate == math.Trunc(rate) {
		return strconv.FormatFloat(rate, 'f', 0, 64) + "%"
	}
	return strconv.FormatFloat(rate, 'f', -1, 64) + "%"
}

func parseAuditLabel(label string) (time.Time, error) {
	label = strings.TrimSpace(label)
	if t, err := time.Parse("January 2, 2006", label); err == nil {
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC), nil
	}
	if t, err := time.Parse("January 02, 2006", label); err == nil {
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC), nil
	}
	return time.Time{}, fmt.Errorf("unsupported audit label")
}

func normalizeAuditDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "/", "-")
	if compactDateRe.MatchString(s) {
		t, err := time.Parse("20060102", s)
		if err != nil {
			return time.Time{}, err
		}
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC), nil
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC), nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return time.Date(t.UTC().Year(), t.UTC().Month(), t.UTC().Day(), 0, 0, 0, 0, time.UTC), nil
	}
	return time.Time{}, fmt.Errorf("unsupported audit id %q (want YYYY-MM-DD, YYYYMMDD, or cms period id)", s)
}

func getBytes(ctx context.Context, rawURL, accept string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", accept)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: HTTP %d: %s", rawURL, resp.StatusCode, truncate(string(body), 200))
	}
	return body, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

// BuildSummaryBundle is unused for Kraken live fetch but kept for import symmetry.
func BuildSummaryBundle(path string) ([]byte, error) {
	return os.ReadFile(path)
}
