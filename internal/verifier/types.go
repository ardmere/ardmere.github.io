// Package verifier defines the plugin contract and core types for ardmere's
// PoR verifications. See docs/verifier-architecture.md §4.
package verifier

import "time"

type Verdict string

const (
	VerdictPass         Verdict = "PASS"
	VerdictFail         Verdict = "FAIL"
	VerdictWarn         Verdict = "WARN"
	VerdictPartial      Verdict = "PARTIAL"
	VerdictUnverifiable Verdict = "UNVERIFIABLE"
)

// Finding is one fact-level outcome inside a Verification.
type Finding struct {
	Subject    string            `json:"subject"`          // e.g. "BTC", "0x57a0…", "merkleRoot"
	Field      string            `json:"field"`            // e.g. "exchangeBalance"
	Claim      string            `json:"claim,omitempty"`  // value as claimed by Binance
	Actual     string            `json:"actual,omitempty"` // value as we measured/recomputed
	Status     Verdict           `json:"status"`
	Note       string            `json:"note,omitempty"`
	Components map[string]string `json:"components,omitempty"` // e.g. liquid/staked/unbonding split
}

// Verification is a single execution of one verifier over one snapshot.
// It is what gets serialized into the verification bundle and Merkle-rooted
// for on-chain anchoring (see §6 of the architecture doc).
type Verification struct {
	VerifierID     string    `json:"verifierId"` // "internal-consistency@1"
	Version        string    `json:"version"`    // "1.0"
	SnapshotID     string    `json:"snapshotId"` // "PR01APR26"
	VerifiedAt     time.Time `json:"verifiedAt"`
	InputArtifacts []string  `json:"inputArtifacts"` // sha256 hex
	Verdict        Verdict   `json:"verdict"`
	Coverage       float64   `json:"coverage"` // 0..1, fraction of subjects checked
	Findings       []Finding `json:"findings,omitempty"`
	Reason         string    `json:"reason,omitempty"` // for UNVERIFIABLE / PARTIAL stubs
}
