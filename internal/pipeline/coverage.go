package pipeline

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/coingecko"
	"github.com/ardmere/ardmere/internal/por"
	"github.com/ardmere/ardmere/internal/verifier"
	"github.com/ardmere/ardmere/internal/walletzip"
)

// CoverageDetail is the assessment-facing breakdown of the two on-chain ratios.
type CoverageDetail struct {
	PublishedAddressCount   int      `json:"publishedAddressCount"`
	VerifiableAddressCount  int      `json:"verifiableAddressCount"`
	OwnershipPassCount      int      `json:"ownershipPassCount"`
	BalancePassCount        int      `json:"balancePassCount"`
	WalletCoverageBps       uint16   `json:"walletCoverageBps"`
	ReserveCoverageBps      uint16   `json:"reserveCoverageBps"`
	VerifiableUSD           string   `json:"verifiableUSD"`
	ClaimedReserveUSD       string   `json:"claimedReserveUSD"`
	CoinGeckoAsOf           string   `json:"coinGeckoAsOf,omitempty"`
	PriceSource             string   `json:"priceSource"`
	ExcludedAssets          []string `json:"excludedAssets,omitempty"`
	Notes                   []string `json:"notes,omitempty"`
	DenominatorUnverifiable bool     `json:"denominatorUnverifiable,omitempty"`
}

// CoverageRatios is what the anchor calldata and assessment carry for schema v4.
type CoverageRatios struct {
	WalletCoverageBps  uint16
	ReserveCoverageBps uint16
	Detail             CoverageDetail
}

// PriceOracle is satisfied by *coingecko.Client.
type PriceOracle interface {
	PriceAt(ctx context.Context, coinID string, at time.Time) (price float64, asOf time.Time, err error)
}

// CoverageInput gathers verifier outputs + published wallet data for ratio math.
type CoverageInput struct {
	Snapshot      por.Snapshot
	Verifications []verifier.Verification
	// PublishedAddressCount is the exchange-published wallet_address_list size.
	// When 0 and PublishedUnknown is false, the zip was empty / missing.
	PublishedAddressCount int
	PublishedUnknown      bool // true when no published address set exists (UNVERIFIABLE denominator)
	// AddressBalances maps lower(address) → coin → amount for reserve numerator.
	AddressBalances map[string]map[string]decimal.Decimal
	Prices          PriceOracle
}

// ComputeCoverageRatios derives walletCoverageBps and reserveCoverageBps.
//
// Definitions (locked):
//   - verifiable address = ownership PASS AND on-chain balance PASS (any of
//     hot/token/ledger/deposit covering that address)
//   - ownership stub/UNVERIFIABLE → both ratios 0 (no balance-only fallback)
//   - published address denominator missing → walletCoverageBps 0 + note
//   - reserve USD uses CoinGecko at snapshotTime on claimed asset amounts;
//     exchange self-reported USD totals are never reused
func ComputeCoverageRatios(ctx context.Context, in CoverageInput) CoverageRatios {
	detail := CoverageDetail{
		PublishedAddressCount: in.PublishedAddressCount,
		PriceSource:           "coingecko_asset_amounts_at_snapshotTime",
	}
	out := CoverageRatios{Detail: detail}

	own := findVerification(in.Verifications, "address-ownership")
	if own == nil || own.Verdict == verifier.VerdictUnverifiable || isOwnershipStub(own) {
		note := "address-ownership is stub/UNVERIFIABLE; walletCoverageBps=0 and reserveCoverageBps=0 (no balance-only fallback)"
		if own != nil && own.Reason != "" {
			note = note + "; " + own.Reason
		}
		out.Detail.Notes = append(out.Detail.Notes, note)
		if in.PublishedUnknown || in.PublishedAddressCount == 0 {
			out.Detail.DenominatorUnverifiable = true
			out.Detail.Notes = append(out.Detail.Notes, "published wallet_address_list denominator is UNVERIFIABLE / empty")
		}
		// Still price claimed reserves so the dry-run / assessment can show the denominator.
		claimed, excluded, asOf, notes := claimedReserveUSD(ctx, in)
		out.Detail.ClaimedReserveUSD = formatUSD(claimed)
		out.Detail.VerifiableUSD = "0"
		out.Detail.ExcludedAssets = excluded
		out.Detail.CoinGeckoAsOf = asOf
		out.Detail.Notes = append(out.Detail.Notes, notes...)
		out.Detail.WalletCoverageBps = 0
		out.Detail.ReserveCoverageBps = 0
		return out
	}

	if in.PublishedUnknown {
		out.Detail.DenominatorUnverifiable = true
		out.Detail.Notes = append(out.Detail.Notes, "exchange published no wallet_address_list; walletCoverageBps=0 (UNVERIFIABLE denominator)")
		out.WalletCoverageBps = 0
	}

	ownershipPass := ownershipPassSet(own, in)
	out.Detail.OwnershipPassCount = len(ownershipPass)

	balancePass := balancePassSet(in.Verifications)
	out.Detail.BalancePassCount = len(balancePass)

	verifiable := map[string]struct{}{}
	for addr := range ownershipPass {
		if _, ok := balancePass[addr]; ok {
			verifiable[addr] = struct{}{}
		}
	}
	out.Detail.VerifiableAddressCount = len(verifiable)

	if !in.PublishedUnknown && in.PublishedAddressCount > 0 {
		out.WalletCoverageBps = ratioBps(float64(len(verifiable)), float64(in.PublishedAddressCount))
	} else if in.PublishedAddressCount == 0 {
		out.Detail.DenominatorUnverifiable = true
		out.Detail.Notes = append(out.Detail.Notes, "published address count is 0; walletCoverageBps=0")
		out.WalletCoverageBps = 0
	}

	if len(balancePass) == 0 {
		out.Detail.Notes = append(out.Detail.Notes,
			"no on-chain balance PASS findings (hot/token/ledger/deposit); verifiable address set is empty")
	}

	claimed, excludedClaimed, asOf1, notes1 := claimedReserveUSD(ctx, in)
	verUSD, excludedVer, asOf2, notes2 := verifiableReserveUSD(ctx, in, verifiable)
	excluded := uniqueStrings(append(excludedClaimed, excludedVer...))
	asOf := asOf1
	if asOf2 != "" {
		asOf = asOf2
	}
	out.Detail.ClaimedReserveUSD = formatUSD(claimed)
	out.Detail.VerifiableUSD = formatUSD(verUSD)
	out.Detail.ExcludedAssets = excluded
	out.Detail.CoinGeckoAsOf = asOf
	out.Detail.Notes = append(out.Detail.Notes, notes1...)
	out.Detail.Notes = append(out.Detail.Notes, notes2...)
	out.Detail.Notes = append(out.Detail.Notes,
		"reserve USD uses CoinGecko prices × exchange claimed asset amounts; exchange self-reported USD totals are not reused")

	if claimed > 0 {
		out.ReserveCoverageBps = ratioBps(verUSD, claimed)
	} else {
		out.ReserveCoverageBps = 0
		out.Detail.Notes = append(out.Detail.Notes, "claimed reserve USD denominator is 0 after exclusions")
	}
	out.Detail.WalletCoverageBps = out.WalletCoverageBps
	out.Detail.ReserveCoverageBps = out.ReserveCoverageBps
	return out
}

func findVerification(vs []verifier.Verification, id string) *verifier.Verification {
	for i := range vs {
		if vs[i].VerifierID == id {
			return &vs[i]
		}
	}
	return nil
}

func isOwnershipStub(v *verifier.Verification) bool {
	if v == nil {
		return true
	}
	if strings.HasPrefix(v.Version, "0") && v.Version != "0.0-ok" {
		// version "0" stubs
		if v.Version == "0" || v.Version == "0.0" {
			return true
		}
	}
	r := strings.ToLower(v.Reason)
	return strings.Contains(r, "not implemented") ||
		strings.Contains(r, "stub") ||
		strings.Contains(r, "no public") && strings.Contains(r, "ownership")
}

func ownershipPassSet(own *verifier.Verification, in CoverageInput) map[string]struct{} {
	out := map[string]struct{}{}
	failed := map[string]struct{}{}
	for _, f := range own.Findings {
		if f.Status == verifier.VerdictFail && looksLikeAddress(f.Subject) {
			failed[normalizeAddr(f.Subject)] = struct{}{}
		}
		if f.Status == verifier.VerdictPass && looksLikeAddress(f.Subject) {
			out[normalizeAddr(f.Subject)] = struct{}{}
		}
	}
	// OKX (and similar) emit only FAIL findings when overall PASS — treat all
	// published addresses as ownership-pass except explicit FAIL subjects.
	if own.Verdict == verifier.VerdictPass && len(out) == 0 {
		if in.AddressBalances != nil {
			for addr := range in.AddressBalances {
				if _, bad := failed[addr]; !bad {
					out[addr] = struct{}{}
				}
			}
		} else if in.PublishedAddressCount > 0 {
			// No zip loaded (e.g. assessment-only dry-run): synthesize count-only
			// markers so OwnershipPassCount reflects the published set. These
			// synthetic keys cannot intersect balance-pass addresses.
			for i := 0; i < in.PublishedAddressCount-len(failed); i++ {
				out[fmt.Sprintf("__published_own_pass_%d", i)] = struct{}{}
			}
		}
	}
	return out
}

func balancePassSet(vs []verifier.Verification) map[string]struct{} {
	out := map[string]struct{}{}
	for _, v := range vs {
		switch v.VerifierID {
		case "onchain-balance-hot", "onchain-balance-token", "onchain-balance-ledger", "onchain-balance-deposit":
		default:
			continue
		}
		for _, f := range v.Findings {
			if f.Status == verifier.VerdictPass && looksLikeAddress(f.Subject) {
				out[normalizeAddr(f.Subject)] = struct{}{}
			}
		}
	}
	return out
}

func claimedReserveUSD(ctx context.Context, in CoverageInput) (usd float64, excluded []string, asOf string, notes []string) {
	if in.Prices == nil {
		notes = append(notes, "no CoinGecko price oracle configured; claimed reserve USD left at 0")
		return 0, nil, "", notes
	}
	at := in.Snapshot.SnapshotTime.UTC()
	var asOfTimes []time.Time
	for _, c := range in.Snapshot.CoinSummaries {
		amt := c.ExchangeReserve.Add(c.ThirdPartyReserve)
		if amt.IsZero() {
			continue
		}
		id := coingecko.ResolveID(c.Coin)
		if id == "" {
			excluded = append(excluded, c.Coin)
			continue
		}
		price, t, err := in.Prices.PriceAt(ctx, id, at)
		if err != nil {
			excluded = append(excluded, c.Coin+" (price error: "+err.Error()+")")
			continue
		}
		f, _ := amt.Float64()
		usd += f * price
		asOfTimes = append(asOfTimes, t)
	}
	if len(excluded) > 0 {
		notes = append(notes, "excluded unmapped/unpriced claimed assets from reserve denominator: "+strings.Join(excluded, ", "))
	}
	if len(asOfTimes) > 0 {
		asOf = asOfTimes[0].UTC().Format(time.RFC3339)
	}
	return usd, uniqueStrings(excluded), asOf, notes
}

func verifiableReserveUSD(ctx context.Context, in CoverageInput, verifiable map[string]struct{}) (usd float64, excluded []string, asOf string, notes []string) {
	if in.Prices == nil || len(verifiable) == 0 || in.AddressBalances == nil {
		return 0, nil, "", notes
	}
	at := in.Snapshot.SnapshotTime.UTC()
	totals := map[string]decimal.Decimal{}
	for addr := range verifiable {
		for coin, amt := range in.AddressBalances[addr] {
			totals[coin] = totals[coin].Add(amt)
		}
	}
	var asOfTimes []time.Time
	for coin, amt := range totals {
		if amt.IsZero() {
			continue
		}
		id := coingecko.ResolveID(coin)
		if id == "" {
			excluded = append(excluded, coin)
			continue
		}
		price, t, err := in.Prices.PriceAt(ctx, id, at)
		if err != nil {
			excluded = append(excluded, coin+" (price error: "+err.Error()+")")
			continue
		}
		f, _ := amt.Float64()
		usd += f * price
		asOfTimes = append(asOfTimes, t)
	}
	if len(excluded) > 0 {
		notes = append(notes, "excluded unmapped/unpriced verifiable assets from reserve numerator: "+strings.Join(excluded, ", "))
	}
	if len(asOfTimes) > 0 {
		asOf = asOfTimes[0].UTC().Format(time.RFC3339)
	}
	return usd, uniqueStrings(excluded), asOf, notes
}

func ratioBps(num, den float64) uint16 {
	if den <= 0 || num <= 0 {
		return 0
	}
	bps := (num / den) * 10000
	if bps > 10000 {
		bps = 10000
	}
	return uint16(math.Round(bps))
}

func formatUSD(v float64) string {
	return strconvFormatFloat(v)
}

func strconvFormatFloat(v float64) string {
	return fmt.Sprintf("%.2f", v)
}

func normalizeAddr(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func looksLikeAddress(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" || s == "summary" || s == "onchain-sample" {
		return false
	}
	if strings.Contains(s, "|") && !strings.HasPrefix(s, "0x") {
		// unsupported pair keys like "BTC|BTC"
		return false
	}
	if strings.HasPrefix(s, "0x") && len(s) >= 42 {
		return true
	}
	// Non-EVM addresses (BTC, TRX, SOL, …) — require some length and no spaces.
	if len(s) >= 20 && !strings.Contains(s, " ") && !strings.Contains(s, "@") {
		return true
	}
	return false
}

func uniqueStrings(in []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, s := range in {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

// CollectPublishedAddresses builds address→coin→amount from a wallet aggregate path helper.
func CollectPublishedAddresses(rows []walletzip.Row) (count int, balances map[string]map[string]decimal.Decimal) {
	balances = map[string]map[string]decimal.Decimal{}
	seen := map[string]struct{}{}
	for _, r := range rows {
		addr := normalizeAddr(r.Address)
		if addr == "" {
			continue
		}
		if _, ok := seen[addr]; !ok {
			seen[addr] = struct{}{}
		}
		if balances[addr] == nil {
			balances[addr] = map[string]decimal.Decimal{}
		}
		coin := strings.ToUpper(r.Coin)
		balances[addr][coin] = balances[addr][coin].Add(r.Balance)
	}
	return len(seen), balances
}

// EncodePorStage maps assessment stage strings to on-chain uint8.
func EncodePorStage(stage string) uint8 {
	switch strings.TrimSpace(stage) {
	case "Pre-Stage":
		return 0
	case "Stage 0":
		return 1
	case "Stage 1":
		return 2
	case "Stage 2":
		return 3
	default:
		return 0
	}
}

// DecodePorStage is the inverse of EncodePorStage for skills/docs.
func DecodePorStage(v uint8) string {
	switch v {
	case 0:
		return "Pre-Stage"
	case 1:
		return "Stage 0"
	case 2:
		return "Stage 1"
	case 3:
		return "Stage 2"
	default:
		return fmt.Sprintf("unknown(%d)", v)
	}
}
