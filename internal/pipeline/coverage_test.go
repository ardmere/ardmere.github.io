package pipeline

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/por"
	"github.com/ardmere/ardmere/internal/verifier"
)

type stubPrices map[string]float64

func (s stubPrices) PriceAt(_ context.Context, coinID string, at time.Time) (float64, time.Time, error) {
	p, ok := s[coinID]
	if !ok {
		return 0, time.Time{}, errPriceMissing(coinID)
	}
	return p, at, nil
}

type priceMissing string

func errPriceMissing(id string) error { return priceMissing(id) }
func (e priceMissing) Error() string  { return "missing " + string(e) }

func TestComputeCoverageRatios_OwnershipUnverifiableZeros(t *testing.T) {
	snap := por.Snapshot{
		Exchange:     "binance",
		ID:           "t1",
		SnapshotTime: time.Date(2026, 10, 31, 16, 0, 0, 0, time.UTC),
		CoinSummaries: []por.CoinSummary{{
			Coin:            "BTC",
			ExchangeReserve: decimal.NewFromInt(10),
		}},
	}
	vs := []verifier.Verification{{
		VerifierID: "address-ownership",
		Version:    "0",
		Verdict:    verifier.VerdictUnverifiable,
		Reason:     "No public download channel for wallet ownership signatures / proofs",
	}}
	got := ComputeCoverageRatios(context.Background(), CoverageInput{
		Snapshot:              snap,
		Verifications:         vs,
		PublishedAddressCount: 100,
		AddressBalances: map[string]map[string]decimal.Decimal{
			"0xabc": {"BTC": decimal.NewFromInt(1)},
		},
		Prices: stubPrices{"bitcoin": 100_000},
	})
	if got.WalletCoverageBps != 0 || got.ReserveCoverageBps != 0 {
		t.Fatalf("expected zeros, got wallet=%d reserve=%d", got.WalletCoverageBps, got.ReserveCoverageBps)
	}
	if got.Detail.ClaimedReserveUSD != "1000000.00" {
		t.Fatalf("claimed usd=%s", got.Detail.ClaimedReserveUSD)
	}
}

func TestComputeCoverageRatios_IntersectionAndPricing(t *testing.T) {
	snap := por.Snapshot{
		Exchange:     "okx",
		ID:           "t2",
		SnapshotTime: time.Date(2026, 10, 31, 16, 0, 0, 0, time.UTC),
		CoinSummaries: []por.CoinSummary{
			{Coin: "BTC", ExchangeReserve: decimal.NewFromInt(100)},
			{Coin: "WEIRD", ExchangeReserve: decimal.NewFromInt(50)},
		},
	}
	vs := []verifier.Verification{
		{
			VerifierID: "address-ownership",
			Version:    "okx-1.1",
			Verdict:    verifier.VerdictPass,
			Reason:     "verified 2/2",
		},
		{
			VerifierID: "onchain-balance-hot",
			Version:    "2.1",
			Verdict:    verifier.VerdictPartial,
			Findings: []verifier.Finding{
				{Subject: "0x1111111111111111111111111111111111111111", Field: "ETH@ETH", Status: verifier.VerdictPass, Claim: "1", Actual: "1"},
				{Subject: "summary", Field: "verified_pairs", Status: verifier.VerdictPass},
			},
		},
	}
	balances := map[string]map[string]decimal.Decimal{
		"0x1111111111111111111111111111111111111111": {"BTC": decimal.NewFromInt(4)},
		"0x2222222222222222222222222222222222222222": {"BTC": decimal.NewFromInt(6)},
	}
	got := ComputeCoverageRatios(context.Background(), CoverageInput{
		Snapshot:              snap,
		Verifications:         vs,
		PublishedAddressCount: 2,
		AddressBalances:       balances,
		Prices:                stubPrices{"bitcoin": 50_000},
	})
	// Only 0x1111 is in both ownership (all published) and balance-pass.
	if got.Detail.VerifiableAddressCount != 1 {
		t.Fatalf("verifiable=%d detail=%+v", got.Detail.VerifiableAddressCount, got.Detail)
	}
	if got.WalletCoverageBps != 5000 {
		t.Fatalf("walletCoverageBps=%d", got.WalletCoverageBps)
	}
	// verifiable USD = 4 BTC * 50k = 200k; claimed = 100 BTC * 50k = 5M; WEIRD excluded
	if got.Detail.VerifiableUSD != "200000.00" {
		t.Fatalf("verifiableUSD=%s", got.Detail.VerifiableUSD)
	}
	if got.Detail.ClaimedReserveUSD != "5000000.00" {
		t.Fatalf("claimed=%s excluded=%v", got.Detail.ClaimedReserveUSD, got.Detail.ExcludedAssets)
	}
	if got.ReserveCoverageBps != 400 { // 200k/5M = 4%
		t.Fatalf("reserveCoverageBps=%d", got.ReserveCoverageBps)
	}
	foundWeird := false
	for _, e := range got.Detail.ExcludedAssets {
		if e == "WEIRD" {
			foundWeird = true
		}
	}
	if !foundWeird {
		t.Fatalf("expected WEIRD excluded: %v", got.Detail.ExcludedAssets)
	}
}

func TestEncodePorStage(t *testing.T) {
	if EncodePorStage("Stage 1") != 2 || DecodePorStage(2) != "Stage 1" {
		t.Fatal("stage mapping")
	}
}
