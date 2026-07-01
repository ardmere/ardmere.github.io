package verifier

import (
	"github.com/shopspring/decimal"
)

// aptosLedgerMismatch decides verdict for APT|APT rows after CoinStore+FA query.
func aptosLedgerMismatch(actual, claim decimal.Decimal, components map[string]string) (Verdict, string) {
	if balanceWithinTolerance(actual, claim) {
		return VerdictPass, ""
	}
	_, surplus := classifyBalanceMismatch(actual, claim)
	if surplus {
		return VerdictWarn, "chain observed > CSV row allocation; Aptos CoinStore+FA APT vs CSV row split"
	}
	if components["coin_store_missing"] == "true" && claim.GreaterThan(decimal.NewFromFloat(1)) {
		return VerdictUnverifiable,
			"no AptosCoin/FungibleStore APT at ledger_version; account empty or CSV row is internal allocation"
	}
	if actual.IsPositive() {
		return VerdictWarn, "partial CoinStore+FA APT vs CSV row allocation; other custody forms not queried"
	}
	if claim.GreaterThan(decimal.NewFromFloat(1)) {
		return VerdictUnverifiable, "zero Aptos CoinStore+FA APT at ledger_version vs non-trivial CSV claim"
	}
	return VerdictWarn, "Aptos CoinStore+FA APT below CSV row allocation"
}
