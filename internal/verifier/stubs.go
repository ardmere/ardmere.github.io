package verifier

import (
	"strings"
	"time"
)

// StubReason returns the default UNVERIFIABLE reason for a verifier id.
func StubReason(exchange, verifierID string) string {
	switch verifierID {
	case "internal-consistency":
		if exchange == "gateio" || exchange == "htx" || exchange == "bybit" || exchange == "bitget" {
			return "No public wallet address CSV; cannot reconcile summary vs address aggregates"
		}
		return "Wallet address bundle not available"
	case "btc-anchor":
		if exchange == "gateio" || exchange == "htx" || exchange == "bybit" || exchange == "bitget" {
			return "Summary does not declare a BTC block height time anchor"
		}
		return "BTC block time anchor verifier not implemented"
	case "onchain-balance-hot":
		if exchange == "gateio" || exchange == "htx" || exchange == "bybit" || exchange == "bitget" {
			return "No public HotCold wallet address list"
		}
		return "Wallet address bundle not available for on-chain audit"
	case "onchain-balance-token":
		if exchange == "gateio" || exchange == "htx" || exchange == "bybit" || exchange == "bitget" {
			return "No public wallet address list for ERC20/BEP20 audit"
		}
		return "Wallet address bundle not available for token audit"
	case "address-ownership":
		if exchange == "okx" {
			return "Wallet address bundle not available for signature verification"
		}
		return "No public download channel for wallet ownership signatures / proofs"
	case "global-zk-proof":
		if exchange == "okx" {
			return "Liability zk-STARK zip not present; download from OKX CDN or set OKX_ZK_STARK_VALIDATOR"
		}
		if exchange == "gateio" {
			return "Global zkmerkle_cex tar.gz requires login; not available from public API"
		}
		if exchange == "htx" {
			return "HTX public-data zk bundle not present; download from PoR Reports or GitHub public-data.zip"
		}
		if exchange == "bybit" {
			return "Bybit uses pure Merkle Tree (no ZK); user myProof.json proves inclusion only"
		}
		if exchange == "bitget" {
			return "Bitget uses pure Merkle Tree with 64-bit truncated SHA-256 (no ZK)"
		}
		return "Global proof.csv / verifying key not publicly distributed"
	case "third-party-attestation":
		if exchange == "bybit" {
			return "Hacken monthly PoR reports are public but not machine-readable in this pipeline"
		}
		if exchange == "bitget" {
			return "No clearly identified independent third-party PoR attestation in machine-readable form"
		}
		return "No public third-party attestation report available"
	case "cross-chain-wrapped":
		return "No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode)"
	default:
		return "Verifier not available for this exchange snapshot"
	}
}

// Stub produces an UNVERIFIABLE verification for a profile-listed stub verifier.
func Stub(verifierRef, exchange, snapshotID, summarySha, reason string) Verification {
	id, version := ParseVerifierRef(verifierRef)
	if reason == "" {
		reason = StubReason(exchange, id)
	}
	inputs := []string{}
	if summarySha != "" {
		inputs = append(inputs, summarySha)
	}
	return Verification{
		VerifierID:     id,
		Version:        version,
		SnapshotID:     snapshotID,
		VerifiedAt:     time.Now().UTC(),
		InputArtifacts: inputs,
		Verdict:        VerdictUnverifiable,
		Reason:         reason,
	}
}

// ParseVerifierRef splits "artifact-integrity@1" into id and version.
func ParseVerifierRef(ref string) (id, version string) {
	id = ref
	version = "0.0"
	if i := strings.LastIndex(ref, "@"); i >= 0 {
		id = ref[:i]
		version = ref[i+1:]
	}
	return id, version
}
