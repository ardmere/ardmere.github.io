package verifier

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestBalanceWithinTolerance(t *testing.T) {
	claim := decimal.RequireFromString("23337.9217")
	actual := decimal.RequireFromString("23337.9216605")
	if !balanceWithinTolerance(actual, claim) {
		t.Fatalf("expected height-boundary ETH delta within tolerance")
	}
	large := decimal.RequireFromString("1000")
	small := decimal.RequireFromString("999.9")
	if balanceWithinTolerance(small, large) {
		t.Fatalf("expected large relative diff to fail")
	}
}

func TestSonicBalanceNote(t *testing.T) {
	claim := decimal.RequireFromString("82000048.8")
	actual := decimal.RequireFromString("48.8")
	note := sonicBalanceNote(claim, actual)
	if note == "" {
		t.Fatalf("expected diagnostic note for large Sonic gap")
	}
}
