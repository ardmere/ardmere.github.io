package reportgen

// proofSetupClass describes whether an exchange's global ZK stack needs a
// trusted setup ceremony (SNARK-class) or uses transparent setup (STARK-class).
type proofSetupClass string

const (
	proofSetupNone        proofSetupClass = "none"        // no global ZK in public set (Gen 1 Merkle)
	proofSetupTransparent proofSetupClass = "transparent" // zk-STARK / transparent-setup systems
	proofSetupTrusted     proofSetupClass = "trusted"     // zk-SNARK / ceremony-required systems
)

func exchangeProofSetup(exchangeID string) proofSetupClass {
	switch exchangeID {
	case "okx":
		return proofSetupTransparent
	case "binance", "gateio":
		return proofSetupTrusted
	case "htx":
		return proofSetupTrusted
	default:
		return proofSetupNone
	}
}
