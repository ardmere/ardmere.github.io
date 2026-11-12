package pipeline

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/por"
)

func TestComputeReserveFigures_amountsAndThirdParty(t *testing.T) {
	snap := por.Snapshot{
		SnapshotTime: time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC),
		CoinSummaries: []por.CoinSummary{
			{
				Coin:              "BTC",
				CustomerLiability: decimal.NewFromInt(2),
				ExchangeReserve:   decimal.NewFromInt(1),
				ThirdPartyReserve: decimal.NewFromInt(1),
				Extra:             map[string]string{"reserveRate": "150"},
			},
			{
				Coin:              "ETH",
				CustomerLiability: decimal.NewFromInt(10),
				ExchangeReserve:   decimal.NewFromInt(11),
			},
			{
				Coin:              "NOTACOIN",
				CustomerLiability: decimal.NewFromInt(1000),
				ExchangeReserve:   decimal.NewFromInt(1000),
			},
			{Coin: "SOL", Extra: map[string]string{"reserveRate": "101"}},
		},
	}
	prices := stubPrices{"bitcoin": 100, "ethereum": 10} // DOGE unmapped
	got := ComputeReserveFigures(context.Background(), snap, prices)
	if got.TotalLiabilityUsd == nil || got.TotalReserveUsd == nil {
		t.Fatalf("totals nil: %+v", got)
	}
	// BTC 2*100 + ETH 10*10 = 300 liability; reserve (1+1)*100 + 11*10 = 310
	if *got.TotalLiabilityUsd != 300 || *got.TotalReserveUsd != 310 {
		t.Fatalf("usd liab=%v res=%v", *got.TotalLiabilityUsd, *got.TotalReserveUsd)
	}
	if len(got.UnpricedCoins) != 1 || got.UnpricedCoins[0] != "NOTACOIN" {
		t.Fatalf("unpriced: %#v", got.UnpricedCoins)
	}
	if len(got.MajorCoinReserveRates) != 2 {
		t.Fatalf("rates: %#v", got.MajorCoinReserveRates)
	}
	btc := got.MajorCoinReserveRates[0]
	if btc.Coin != "BTC" || btc.Basis != "amounts" || !btc.SelfReported || btc.ReserveRate.String() != "1" {
		t.Fatalf("btc rate: %+v", btc)
	}
	eth := got.MajorCoinReserveRates[1]
	if eth.Basis != "amounts" || eth.ReserveRate.String() != "1.1" {
		t.Fatalf("eth rate: %+v", eth)
	}
	btcBps, ethBps, usdtBps, usdcBps := got.RatioSlots()
	if btcBps != 10000 || ethBps != 11000 || usdtBps != ReserveRatioUnavailable || usdcBps != ReserveRatioUnavailable {
		t.Fatalf("bps %d %d %d %d", btcBps, ethBps, usdtBps, usdcBps)
	}
	if OnChainUSD(got.TotalLiabilityUsd).Cmp(big.NewInt(300)) != 0 {
		t.Fatalf("onchain usd %s", OnChainUSD(got.TotalLiabilityUsd))
	}
}

func TestComputeReserveFigures_reportedRateOnly(t *testing.T) {
	snap := por.Snapshot{
		Extra: map[string]string{"amountsPublished": "false", "totalReserveRate": "118.49"},
		CoinSummaries: []por.CoinSummary{
			{Coin: "BTC", Extra: map[string]string{"reserveRate": "102.5", "reserveAmount": "unpublished", "liabilityAmount": "unpublished"}},
			{Coin: "USDT", Extra: map[string]string{"capitalRatio": "106%"}},
			{Coin: "NOTACOIN", Extra: map[string]string{"reserveRate": "99"}},
		},
	}
	got := ComputeReserveFigures(context.Background(), snap, stubPrices{"bitcoin": 1})
	if got.TotalLiabilityUsd != nil || got.TotalReserveUsd != nil || got.AmountsPublished {
		t.Fatalf("expected null totals, got %+v", got)
	}
	if len(got.UnpricedCoins) != 0 {
		t.Fatalf("unpriced: %#v", got.UnpricedCoins)
	}
	if len(got.MajorCoinReserveRates) != 2 {
		t.Fatalf("rates: %#v", got.MajorCoinReserveRates)
	}
	if got.MajorCoinReserveRates[0].Basis != "reported" || got.MajorCoinReserveRates[0].ReserveRate.String() != "1.025" {
		t.Fatalf("btc: %+v", got.MajorCoinReserveRates[0])
	}
	if got.MajorCoinReserveRates[1].Coin != "USDT" || got.MajorCoinReserveRates[1].ReserveRate.String() != "1.06" {
		t.Fatalf("usdt: %+v", got.MajorCoinReserveRates[1])
	}
	if OnChainUSD(nil).Sign() != 0 {
		t.Fatal("nil usd must encode 0")
	}
	if OnChainRatioBps(got.MajorCoinReserveRates, "USDC") != ReserveRatioUnavailable {
		t.Fatal("missing coin must be 65535")
	}
}

func TestComputeReserveFigures_noPriceOracleLeavesTotalsNull(t *testing.T) {
	snap := por.Snapshot{
		CoinSummaries: []por.CoinSummary{{
			Coin:              "BTC",
			CustomerLiability: decimal.NewFromInt(1),
			ExchangeReserve:   decimal.NewFromInt(2),
		}},
	}
	got := ComputeReserveFigures(context.Background(), snap, nil)
	if got.TotalLiabilityUsd != nil || !got.AmountsPublished {
		t.Fatalf("%+v", got)
	}
	if len(got.UnpricedCoins) != 1 || got.MajorCoinReserveRates[0].Basis != "amounts" {
		t.Fatalf("%+v", got)
	}
}

func TestComputeReserveFigures_missingReservesUnavailable(t *testing.T) {
	snap := por.Snapshot{
		SnapshotTime: time.Date(2026, 10, 27, 6, 0, 0, 0, time.UTC),
		CoinSummaries: []por.CoinSummary{
			{
				Coin:              "BTC",
				CustomerLiability: decimal.NewFromInt(2),
				Extra:             map[string]string{"reserveAmount": "unpublished", "reserveRate": "unpublished"},
			},
			{
				Coin:              "ETH",
				CustomerLiability: decimal.NewFromInt(10),
				Extra:             map[string]string{"reserveAmount": "unpublished", "reserveRate": "unpublished"},
			},
			{
				Coin:              "USDT",
				CustomerLiability: decimal.NewFromInt(100),
				Extra:             map[string]string{"reserveAmount": "unpublished", "reserveRate": "unpublished"},
			},
			{
				Coin:              "USDC",
				CustomerLiability: decimal.NewFromInt(50),
				Extra:             map[string]string{"reserveAmount": "unpublished", "reserveRate": "unpublished"},
			},
		},
	}
	prices := stubPrices{"bitcoin": 100, "ethereum": 10, "tether": 1, "usd-coin": 1}
	got := ComputeReserveFigures(context.Background(), snap, prices)
	if !got.AmountsPublished {
		t.Fatalf("liabilities were published: %+v", got)
	}
	if got.TotalLiabilityUsd == nil || *got.TotalLiabilityUsd != 450 {
		// 2*100 + 10*10 + 100*1 + 50*1 = 450
		t.Fatalf("liab=%v want 450", got.TotalLiabilityUsd)
	}
	if got.TotalReserveUsd != nil {
		t.Fatalf("missing reserves must stay nil, got %v", *got.TotalReserveUsd)
	}
	if len(got.MajorCoinReserveRates) != 0 {
		t.Fatalf("major rates must be omitted when unpublished, got %#v", got.MajorCoinReserveRates)
	}
	btc, eth, usdt, usdc := got.RatioSlots()
	if btc != ReserveRatioUnavailable || eth != ReserveRatioUnavailable || usdt != ReserveRatioUnavailable || usdc != ReserveRatioUnavailable {
		t.Fatalf("on-chain slots want 65535, got %d %d %d %d", btc, eth, usdt, usdc)
	}
	if OnChainUSD(got.TotalReserveUsd).Sign() != 0 {
		t.Fatal("unavailable reserve must encode on-chain as 0 USD")
	}
}

func TestComputeReserveFigures_publishedZeroReserveDistinctFromMissing(t *testing.T) {
	snap := por.Snapshot{
		SnapshotTime: time.Date(2026, 10, 27, 6, 0, 0, 0, time.UTC),
		CoinSummaries: []por.CoinSummary{{
			Coin:              "BTC",
			CustomerLiability: decimal.NewFromInt(1),
			ExchangeReserve:   decimal.Zero, // published zero
		}},
	}
	got := ComputeReserveFigures(context.Background(), snap, stubPrices{"bitcoin": 50})
	if got.TotalLiabilityUsd == nil || *got.TotalLiabilityUsd != 50 {
		t.Fatalf("liab=%v", got.TotalLiabilityUsd)
	}
	if got.TotalReserveUsd == nil || *got.TotalReserveUsd != 0 {
		t.Fatalf("published zero reserve must be 0, not nil: %+v", got.TotalReserveUsd)
	}
	if len(got.MajorCoinReserveRates) != 1 || got.MajorCoinReserveRates[0].ReserveRate.String() != "0" {
		t.Fatalf("published zero rate: %#v", got.MajorCoinReserveRates)
	}
}

func TestOnChainRatioBps_saturatesBelowSentinel(t *testing.T) {
	rates := []MajorCoinReserveRate{{Coin: "BTC", ReserveRate: "10", Basis: "amounts", SelfReported: true}}
	if OnChainRatioBps(rates, "BTC") != ReserveRatioBpsMax {
		t.Fatalf("got %d", OnChainRatioBps(rates, "BTC"))
	}
}
