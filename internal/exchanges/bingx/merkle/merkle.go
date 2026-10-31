// Package merkle validates BingX user Merkle inclusion proofs.
// Logic ported from github.com/BingX-API/BingX-proof-of-reserves por.py.
package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// Proof is a portable JSON shape for BingX leaf + path + published root.
// The official CLI takes CLI args; ardmere archives this object as userMerkleProof.
type Proof struct {
	MerkleLeaf string   `json:"merkleLeaf"`
	MerklePath string   `json:"merklePath"` // "l:hash,r:hash,..." or "hash:l,hash:r" (official README mixes forms)
	MerkleRoot string   `json:"merkleRoot"`
	AuditID    string   `json:"auditId,omitempty"`
	AssetLine  string   `json:"assetLine,omitempty"` // optional raw leaf preimage
	PathList   []string `json:"pathList,omitempty"`  // optional alternate path encoding
}

func ValidateProof(raw []byte) (bool, error) {
	var p Proof
	if err := json.Unmarshal(raw, &p); err != nil {
		return false, err
	}
	if p.MerkleLeaf == "" && p.AssetLine != "" {
		p.MerkleLeaf = HashString(p.AssetLine)
	}
	if p.MerkleLeaf == "" {
		return false, fmt.Errorf("proof missing merkleLeaf")
	}
	root, err := GenerateRoot(p.MerkleLeaf, pathString(p))
	if err != nil {
		return false, err
	}
	if p.MerkleRoot == "" {
		return false, fmt.Errorf("proof missing merkleRoot to compare")
	}
	return strings.EqualFold(root, strings.TrimPrefix(strings.TrimSpace(p.MerkleRoot), "0x")), nil
}

func RootHash(raw []byte) (string, error) {
	var p Proof
	if err := json.Unmarshal(raw, &p); err != nil {
		return "", err
	}
	if p.MerkleRoot != "" {
		return strings.ToLower(strings.TrimPrefix(strings.TrimSpace(p.MerkleRoot), "0x")), nil
	}
	leaf := p.MerkleLeaf
	if leaf == "" && p.AssetLine != "" {
		leaf = HashString(p.AssetLine)
	}
	return GenerateRoot(leaf, pathString(p))
}

func AuditID(raw []byte) (string, error) {
	var p Proof
	if err := json.Unmarshal(raw, &p); err != nil {
		return "", err
	}
	return p.AuditID, nil
}

func pathString(p Proof) string {
	if strings.TrimSpace(p.MerklePath) != "" {
		return p.MerklePath
	}
	return strings.Join(p.PathList, ",")
}

// HashString is SHA-256 hex of the UTF-8 input (official `por.py hash`).
func HashString(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// GenerateRoot recomputes the Merkle root from leaf + path (official `por.py verify`).
// Path entries accept either "l:hash" / "r:hash" (README) or "hash:l" / "hash:r" (script comment).
func GenerateRoot(merkleLeaf, path string) (string, error) {
	if merkleLeaf == "" {
		return "", fmt.Errorf("empty merkle leaf")
	}
	root := strings.ToLower(strings.TrimSpace(merkleLeaf))
	path = strings.TrimSpace(path)
	if path == "" {
		return root, nil
	}
	for _, part := range strings.Split(path, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		side, sibling, err := parsePathEntry(part)
		if err != nil {
			return "", err
		}
		sibling = strings.ToLower(sibling)
		switch side {
		case "l":
			root = HashString(sibling + root)
		case "r":
			root = HashString(root + sibling)
		default:
			return "", fmt.Errorf("invalid path side %q in %q", side, part)
		}
	}
	return root, nil
}

func parsePathEntry(part string) (side, sibling string, err error) {
	chunks := strings.Split(part, ":")
	if len(chunks) != 2 {
		return "", "", fmt.Errorf("path entry must be side:hash or hash:side, got %q", part)
	}
	a, b := strings.ToLower(strings.TrimSpace(chunks[0])), strings.ToLower(strings.TrimSpace(chunks[1]))
	switch {
	case a == "l" || a == "r":
		return a, b, nil
	case b == "l" || b == "r":
		return b, a, nil
	default:
		return "", "", fmt.Errorf("path entry missing side marker: %q", part)
	}
}
