package verifier

// onchainPairCoverage reports which on-chain verifier configs include a (coin|network) pair.
func onchainPairCoverage(exchange, coin, network string) (native, token, ledger bool) {
	key := coin + "|" + network
	_, native = loadNativeSupported(exchange)[key]
	_, token = loadTokenSupportedFor(exchange)[key]
	_, ledger = loadLedgerSupportedFor(exchange)[key]
	return native, token, ledger
}

// shouldReportUnsupportedHot is true when no verifier config covers this pair.
func shouldReportUnsupportedHot(exchange, coin, network string) bool {
	native, token, ledger := onchainPairCoverage(exchange, coin, network)
	return !native && !token && !ledger
}

// shouldReportUnsupportedToken is true when token verifier does not cover the pair and
// no other on-chain verifier owns it.
func shouldReportUnsupportedToken(exchange, coin, network string) bool {
	key := coin + "|" + network
	if _, ok := loadTokenSupportedFor(exchange)[key]; ok {
		return false
	}
	if !isEVMTokenCandidate(coin) {
		return false
	}
	native, _, ledger := onchainPairCoverage(exchange, coin, network)
	return !native && !ledger
}

// shouldReportUnsupportedLedger is true when ledger verifier does not cover the pair and
// no other on-chain verifier owns it.
func shouldReportUnsupportedLedger(exchange, coin, network string) bool {
	key := coin + "|" + network
	if _, ok := loadLedgerSupportedFor(exchange)[key]; ok {
		return false
	}
	if !isLedgerCandidate(coin) {
		return false
	}
	native, token, _ := onchainPairCoverage(exchange, coin, network)
	return !native && !token
}
