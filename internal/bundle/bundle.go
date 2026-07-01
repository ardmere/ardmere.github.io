// Package bundle handles the artifact-bundle and verification-bundle
// manifest types + their canonical Merkle roots used for on-chain anchoring.
package bundle

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"

	"github.com/ardmere/ardmere/internal/verifier"
)

type Artifact struct {
	Kind      string `json:"kind"`
	SHA256    string `json:"sha256"`
	URL       string `json:"url"`
	FetchedAt string `json:"fetchedAt"`
	SizeBytes int64  `json:"sizeBytes"`
	LocalPath string `json:"localPath,omitempty"`
}

// ArtifactBundle is what we Merkle-root for the "artifact-bundle" anchor.
type ArtifactBundle struct {
	Exchange     string     `json:"exchange"`     // e.g. "binance"
	SnapshotID   string     `json:"snapshotId"`
	PeriodSeq    uint32     `json:"periodSeq"`
	SnapshotTime string     `json:"snapshotTime"` // RFC3339
	BTCBlockHeight uint32   `json:"btcBlockHeight"`
	ExchangeMerkleRoot string `json:"exchangeMerkleRoot,omitempty"`
	BuiltAt      string     `json:"builtAt"`
	Kind         string     `json:"kind"` // "artifact-bundle"
	Artifacts    []Artifact `json:"artifacts"`
	MerkleRoot   string     `json:"merkleRoot"` // 0x-prefixed hex
}

// VerificationBundle is what we Merkle-root for the "verification-bundle" anchor.
type VerificationBundle struct {
	Exchange           string                  `json:"exchange"`
	SnapshotID         string                  `json:"snapshotId"`
	PeriodSeq          uint32                  `json:"periodSeq"`
	SnapshotTime       string                  `json:"snapshotTime"`
	BTCBlockHeight     uint32                  `json:"btcBlockHeight"`
	ExchangeMerkleRoot string                  `json:"exchangeMerkleRoot,omitempty"`
	VerdictSummary     uint8                   `json:"verdictSummary"`
	CoverageBps        uint16                  `json:"coverageBps"`
	BuiltAt            string                  `json:"builtAt"`
	Kind               string                  `json:"kind"` // "verification-bundle"
	Verifications      []verifier.Verification `json:"verifications"`
	MerkleRoot         string                  `json:"merkleRoot"`
}

// ArtifactRoot computes a deterministic sha256 Merkle root over a set of
// artifacts. Leaves are: sha256("<sha256>|<url>|<fetchedAt>"), sorted by
// the artifact sha256 ascending. Last odd leaf is duplicated.
func ArtifactRoot(arts []Artifact) [32]byte {
	if len(arts) == 0 {
		return [32]byte{}
	}
	sorted := make([]Artifact, len(arts))
	copy(sorted, arts)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].SHA256 < sorted[j].SHA256 })

	leaves := make([][32]byte, len(sorted))
	for i, a := range sorted {
		s := a.SHA256 + "|" + a.URL + "|" + a.FetchedAt
		leaves[i] = sha256.Sum256([]byte(s))
	}
	return reduceLevels(leaves)
}

// VerificationRoot computes a Merkle root over a set of verifications.
// Leaf = sha256(canonical_json(verification)). Sorted by VerifierID ascending.
func VerificationRoot(vs []verifier.Verification) [32]byte {
	if len(vs) == 0 {
		return [32]byte{}
	}
	sorted := make([]verifier.Verification, len(vs))
	copy(sorted, vs)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].VerifierID < sorted[j].VerifierID })

	leaves := make([][32]byte, len(sorted))
	for i, v := range sorted {
		buf, _ := json.Marshal(v) // Go's json.Marshal is deterministic for fixed structs
		leaves[i] = sha256.Sum256(buf)
	}
	return reduceLevels(leaves)
}

func reduceLevels(level [][32]byte) [32]byte {
	for len(level) > 1 {
		if len(level)%2 == 1 {
			level = append(level, level[len(level)-1])
		}
		next := make([][32]byte, 0, len(level)/2)
		for i := 0; i < len(level); i += 2 {
			pair := append([]byte{}, level[i][:]...)
			pair = append(pair, level[i+1][:]...)
			next = append(next, sha256.Sum256(pair))
		}
		level = next
	}
	return level[0]
}

// HexRoot returns the 0x-prefixed hex representation of a 32-byte root.
func HexRoot(r [32]byte) string { return "0x" + hex.EncodeToString(r[:]) }
