package verifier

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestAptosLedgerMismatch(t *testing.T) {
	claim := decimal.NewFromInt(1_000_000)
	cases := []struct {
		name   string
		actual decimal.Decimal
		comp   map[string]string
		want   Verdict
	}{
		{name: "pass", actual: claim, comp: nil, want: VerdictPass},
		{name: "surplus", actual: claim.Add(decimal.NewFromInt(100)), comp: nil, want: VerdictWarn},
		{name: "missing store", actual: decimal.Zero, comp: map[string]string{"coin_store_missing": "true"}, want: VerdictUnverifiable},
		{name: "partial", actual: decimal.NewFromInt(500_000), comp: nil, want: VerdictWarn},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := aptosLedgerMismatch(tc.actual, claim, tc.comp)
			if got != tc.want {
				t.Fatalf("got %s want %s", got, tc.want)
			}
		})
	}
}
