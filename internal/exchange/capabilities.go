package exchange

// Tier describes how much independent verification an exchange profile supports.
const (
	TierSummary = 1 // public summary only
	TierWallet  = 2 // wallet zip + address signatures
	TierOnchain = 3 // multi-chain onchain balance audit
)

// Capabilities documents what an exchange adapter can verify today.
type Capabilities struct {
	Tier             int
	WalletZip        bool
	AddressOwnership bool
	GlobalZK         bool
	OnchainHot       bool
	OnchainToken     bool
	OnchainLedger    bool
}
