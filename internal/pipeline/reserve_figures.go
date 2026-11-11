package pipeline

import (
	"context"
	"encoding/json"
	"math"
	"math/big"
	"sort"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/coingecko"
	"github.com/ardmere/ardmere/internal/por"
)

const (
	// ReserveRatioUnavailable is the on-chain uint16 sentinel (schema v7).
	// 10000 bps = 100%. 65535 = this coin has no row, or no computable ratio.
	ReserveRatioUnavailable uint16 = 65535
	// ReserveRatioBpsMax is the largest ratio that still fits under the sentinel.
	ReserveRatioBpsMax uint16 = 65534
)

// MajorCoins are the reserve-ratio slots anchored on-chain, in calldata order.
var MajorCoins = []string{"BTC", "ETH", "USDT", "USDC"}

// MajorCoinReserveRate is one self-reported reserve ratio for a major coin.
// ReserveRate is reserve/liability where 1 = 100% (exchange percentages ÷ 100).
type MajorCoinReserveRate struct {
	Coin         string      `json:"coin"`
	ReserveRate  json.Number `json:"reserveRate"`
	Basis        string      `json:"basis"` // "amounts" or "reported"
	SelfReported bool        `json:"selfReported"`
}

// ReserveFigures is the assessment block for liabilities, reserves, and major-coin ratios.
// USD totals are nil when the snapshot did not publish amounts, or when every
// amount-bearing coin was unpriced. They are self-reported exchange amounts
// priced with CoinGecko at snapshotTime — not a solvency check.
type ReserveFigures struct {
	TotalLiabilityUsd     *float64               `json:"totalLiabilityUsd"`
	TotalReserveUsd       *float64               `json:"totalReserveUsd"`
	UnpricedCoins         []string               `json:"unpricedCoins"`
	MajorCoinReserveRates []MajorCoinReserveRate `json:"majorCoinReserveRates"`
	// SourceMissing is true when the summary artifact could not be read.
	// Totals stay nil; the report must not claim the exchange published no amounts.
	SourceMissing bool `json:"-"`
	// AmountsPublished is true when at least one coin row has a customer
	// liability or an exchange reserve (self-custody + third-party).
	AmountsPublished bool   `json:"-"`
	PriceSource      string `json:"-"`
}

// ReserveFiguresSourceMissing is the assessment shape when no summary was loaded.
func ReserveFiguresSourceMissing() ReserveFigures {
	return ReserveFigures{
		UnpricedCoins:         []string{},
		MajorCoinReserveRates: []MajorCoinReserveRate{},
		SourceMissing:         true,
		PriceSource:           "coingecko_asset_amounts_at_snapshotTime",
	}
}

// ComputeReserveFigures prices self-reported per-coin amounts at snapshotTime.
// Major-coin ratios use exchange reserve / customer liability when both amounts
// exist; otherwise the exchange-reported rate already parsed into CoinSummary.Extra.
func ComputeReserveFigures(ctx context.Context, snap por.Snapshot, prices PriceOracle) ReserveFigures {
	out := ReserveFigures{
		UnpricedCoins:         []string{},
		MajorCoinReserveRates: majorRates(snap),
		PriceSource:           "coingecko_asset_amounts_at_snapshotTime",
	}
	if snapshotHidesAmounts(snap) {
		return out
	}

	type row struct {
		coin      string
		liability decimal.Decimal
		reserve   decimal.Decimal
	}
	var rows []row
	for _, c := range snap.CoinSummaries {
		if coinHidesAmounts(c) {
			continue
		}
		liab := c.CustomerLiability
		res := c.ExchangeReserve.Add(c.ThirdPartyReserve)
		if liab.IsZero() && res.IsZero() {
			continue
		}
		rows = append(rows, row{coin: c.Coin, liability: liab, reserve: res})
	}
	if len(rows) == 0 {
		return out
	}
	out.AmountsPublished = true
	if prices == nil {
		for _, r := range rows {
			out.UnpricedCoins = append(out.UnpricedCoins, r.coin)
		}
		out.UnpricedCoins = uniqueSorted(out.UnpricedCoins)
		return out
	}

	at := snap.SnapshotTime.UTC()
	var liabUSD, resUSD float64
	priced := 0
	for _, r := range rows {
		id := coingecko.ResolveID(r.coin)
		if id == "" {
			out.UnpricedCoins = append(out.UnpricedCoins, r.coin)
			continue
		}
		price, _, err := prices.PriceAt(ctx, id, at)
		if err != nil || price <= 0 {
			note := r.coin
			if err != nil {
				note = r.coin + " (price error: " + err.Error() + ")"
			}
			out.UnpricedCoins = append(out.UnpricedCoins, note)
			continue
		}
		lf, _ := r.liability.Float64()
		rf, _ := r.reserve.Float64()
		liabUSD += lf * price
		resUSD += rf * price
		priced++
	}
	out.UnpricedCoins = uniqueSorted(out.UnpricedCoins)
	if priced == 0 {
		return out
	}
	out.TotalLiabilityUsd = &liabUSD
	out.TotalReserveUsd = &resUSD
	return out
}

func majorRates(snap por.Snapshot) []MajorCoinReserveRate {
	byCoin := map[string]por.CoinSummary{}
	for _, c := range snap.CoinSummaries {
		key := strings.ToUpper(strings.TrimSpace(c.Coin))
		if _, ok := byCoin[key]; ok {
			continue
		}
		byCoin[key] = c
	}
	var out []MajorCoinReserveRate
	for _, coin := range MajorCoins {
		c, ok := byCoin[coin]
		if !ok {
			continue
		}
		rate, basis, ok := coinReserveRatio(c)
		if !ok {
			continue
		}
		out = append(out, MajorCoinReserveRate{
			Coin:         coin,
			ReserveRate:  formatRate(rate),
			Basis:        basis,
			SelfReported: true,
		})
	}
	if out == nil {
		out = []MajorCoinReserveRate{}
	}
	return out
}

func coinReserveRatio(c por.CoinSummary) (rate float64, basis string, ok bool) {
	if !coinHidesAmounts(c) {
		liab := c.CustomerLiability
		res := c.ExchangeReserve.Add(c.ThirdPartyReserve)
		if liab.IsPositive() && !res.IsNegative() {
			f, _ := res.Div(liab).Float64()
			return f, "amounts", true
		}
	}
	if reported, ok := reportedRatio(c); ok {
		return reported, "reported", true
	}
	return 0, "", false
}

func snapshotHidesAmounts(snap por.Snapshot) bool {
	return snap.Extra != nil && snap.Extra["amountsPublished"] == "false"
}

func coinHidesAmounts(c por.CoinSummary) bool {
	if c.Extra == nil {
		return false
	}
	if c.Extra["reserveAmount"] == "unpublished" || c.Extra["liabilityAmount"] == "unpublished" {
		return true
	}
	return false
}

// reportedRatio reads CoinSummary.Extra reserveRate or capitalRatio.
// Those parsers store a percentage (100 = 100%), sometimes with a trailing %.
func reportedRatio(c por.CoinSummary) (float64, bool) {
	if c.Extra == nil {
		return 0, false
	}
	for _, key := range []string{"reserveRate", "capitalRatio"} {
		raw := strings.TrimSpace(c.Extra[key])
		if raw == "" || raw == "unpublished" {
			continue
		}
		raw = strings.TrimSuffix(raw, "%")
		raw = strings.TrimSuffix(raw, "+")
		raw = strings.TrimSpace(raw)
		d, err := decimal.NewFromString(raw)
		if err != nil || d.IsNegative() {
			continue
		}
		f, _ := d.Div(decimal.NewFromInt(100)).Float64()
		return f, true
	}
	return 0, false
}

// OnChainUSD encodes a USD total as uint128 whole dollars.
// nil (unavailable) and a non-positive sum encode as 0.
// A positive sum that rounds to 0 encodes as 1 so it is not the unavailable sentinel.
func OnChainUSD(v *float64) *big.Int {
	if v == nil || math.IsNaN(*v) || math.IsInf(*v, 0) || *v <= 0 {
		return big.NewInt(0)
	}
	rounded := math.Round(*v)
	if rounded < 1 {
		return big.NewInt(1)
	}
	if rounded > math.MaxInt64 {
		return big.NewInt(math.MaxInt64)
	}
	return big.NewInt(int64(rounded))
}

// OnChainRatioBps returns the uint16 slot for one major coin.
// Missing coin or missing ratio → 65535. 10000 = 100%. Values above 655.34% saturate at 65534.
func OnChainRatioBps(rates []MajorCoinReserveRate, coin string) uint16 {
	for _, r := range rates {
		if !strings.EqualFold(r.Coin, coin) {
			continue
		}
		f, err := r.ReserveRate.Float64()
		if err != nil || math.IsNaN(f) || math.IsInf(f, 0) || f < 0 {
			return ReserveRatioUnavailable
		}
		bps := math.Round(f * 10000)
		if bps > float64(ReserveRatioBpsMax) {
			return ReserveRatioBpsMax
		}
		return uint16(bps)
	}
	return ReserveRatioUnavailable
}

// AssessmentMap is the four JSON fields, with null totals when unavailable.
func (f ReserveFigures) AssessmentMap() map[string]any {
	var liab, res any
	if f.TotalLiabilityUsd != nil {
		liab = json.Number(strconv.FormatFloat(round2(*f.TotalLiabilityUsd), 'f', 2, 64))
	}
	if f.TotalReserveUsd != nil {
		res = json.Number(strconv.FormatFloat(round2(*f.TotalReserveUsd), 'f', 2, 64))
	}
	rates := f.MajorCoinReserveRates
	if rates == nil {
		rates = []MajorCoinReserveRate{}
	}
	unpriced := f.UnpricedCoins
	if unpriced == nil {
		unpriced = []string{}
	}
	return map[string]any{
		"totalLiabilityUsd":     liab,
		"totalReserveUsd":       res,
		"unpricedCoins":         unpriced,
		"majorCoinReserveRates": rates,
	}
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func formatRate(v float64) json.Number {
	s := strconv.FormatFloat(v, 'f', 8, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "" || s == "-" {
		s = "0"
	}
	return json.Number(s)
}

func uniqueSorted(in []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, s := range in {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

// RatioSlots is BTC, ETH, USDT, USDC bps in that order.
func (f ReserveFigures) RatioSlots() (btc, eth, usdt, usdc uint16) {
	return OnChainRatioBps(f.MajorCoinReserveRates, "BTC"),
		OnChainRatioBps(f.MajorCoinReserveRates, "ETH"),
		OnChainRatioBps(f.MajorCoinReserveRates, "USDT"),
		OnChainRatioBps(f.MajorCoinReserveRates, "USDC")
}
