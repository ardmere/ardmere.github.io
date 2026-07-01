package verifier

import (
	"testing"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/por"
)

func TestSolvencyClaimGatePass(t *testing.T) {
	v := SolvencyClaim{
		SummarySha256: "abc",
		Snapshot: por.Snapshot{
			Exchange: "gateio",
			ID:       "20231110",
			Extra: map[string]string{
				"totalReserveRate": "115.69",
			},
			CoinSummaries: []por.CoinSummary{{
				Coin:  "BTC",
				Extra: map[string]string{"reserveRate": "123.52"},
			}, {
				Coin:  "ETH",
				Extra: map[string]string{"reserveRate": "101.0"},
			}},
		},
	}.Run()

	if v.Verdict != VerdictPass {
		t.Fatalf("verdict=%s findings=%+v", v.Verdict, v.Findings)
	}
	if len(v.Findings) != 3 {
		t.Fatalf("findings=%+v", v.Findings)
	}
}

func TestSolvencyClaimGateFailTotal(t *testing.T) {
	v := SolvencyClaim{
		SummarySha256: "abc",
		Snapshot: por.Snapshot{
			Exchange: "gateio",
			ID:       "bad",
			Extra: map[string]string{
				"totalReserveRate": "99.5",
			},
		},
	}.Run()
	if v.Verdict != VerdictFail {
		t.Fatalf("verdict=%s", v.Verdict)
	}
}

func TestSolvencyClaimBinancePass(t *testing.T) {
	v := SolvencyClaim{
		SummarySha256: "abc",
		Snapshot: por.Snapshot{
			Exchange: "binance",
			ID:       "PR01JUN26",
			CoinSummaries: []por.CoinSummary{{
				Coin:              "BTC",
				CustomerLiability: decimal.NewFromInt(100),
				ExchangeLiability: decimal.NewFromInt(110),
			}},
		},
	}.Run()
	if v.Verdict != VerdictPass {
		t.Fatalf("verdict=%s findings=%+v", v.Verdict, v.Findings)
	}
}

func TestSolvencyClaimBinanceFail(t *testing.T) {
	v := SolvencyClaim{
		SummarySha256: "abc",
		Snapshot: por.Snapshot{
			Exchange: "binance",
			ID:       "PR01JUN26",
			CoinSummaries: []por.CoinSummary{{
				Coin:              "BTC",
				CustomerLiability: decimal.NewFromInt(110),
				ExchangeLiability: decimal.NewFromInt(100),
			}},
		},
	}.Run()
	if v.Verdict != VerdictFail {
		t.Fatalf("verdict=%s", v.Verdict)
	}
}
