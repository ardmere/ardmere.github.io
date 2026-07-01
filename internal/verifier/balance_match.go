package verifier

import "github.com/shopspring/decimal"

var (
	balanceAbsTol = decimal.NewFromFloat(1e-4)
	balanceRelTol = decimal.NewFromFloat(2e-7) // CSV display rounding; e.g. ETH@ARBITRUM 2122.637163
)

func balanceWithinTolerance(actual, claim decimal.Decimal) bool {
	diff := actual.Sub(claim).Abs()
	if diff.LessThanOrEqual(balanceAbsTol) {
		return true
	}
	maxAbs := decimal.Max(actual.Abs(), claim.Abs())
	if maxAbs.IsZero() {
		return true
	}
	return diff.Div(maxAbs).LessThanOrEqual(balanceRelTol)
}
