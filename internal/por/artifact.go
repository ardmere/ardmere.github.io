package por

// Artifact kind constants. Legacy Binance-specific kinds remain accepted when parsing bundles.
const (
	KindSummarySnapshot     = "summarySnapshot"
	KindBapiSnapshot        = "bapiSnapshot" // legacy Binance summary
	KindWalletZip           = "walletZip"
	KindWalletAddressBundle = "walletAddressBundle"
	KindGlobalProofBundle   = "globalProofBundle" // Gate zkmerkle_cex tar.gz contents
	KindUserMerkleProof     = "userMerkleProof"   // Bybit myProof.json (login-gated)
)

// IsSummaryKind reports whether kind is a normalized or legacy summary artifact.
func IsSummaryKind(kind string) bool {
	return kind == KindSummarySnapshot || kind == KindBapiSnapshot
}

// IsWalletKind reports whether kind is a normalized or legacy wallet bundle artifact.
func IsWalletKind(kind string) bool {
	return kind == KindWalletZip || kind == KindWalletAddressBundle
}
