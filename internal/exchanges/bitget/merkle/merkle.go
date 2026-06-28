// Package merkle validates Bitget user Merkle proofs.
//
// Bitget's public verifier truncates SHA-256 to the first 16 hex characters
// (64 bits). ardmere reproduces that behavior for compatibility, but verifier
// findings must treat a PASS as weak evidence only.
package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/shopspring/decimal"
)

const HashHexLen = 16

type Proof struct {
	AuditID    string            `json:"auditId"`
	EncryptUID string            `json:"encryptUid"`
	Nonce      string            `json:"nonce"`
	Hash       string            `json:"hash"`
	Balances   map[string]string `json:"balances"`
	MerkelPath []PathNode        `json:"merkelPath"`
	MerklePath []PathNode        `json:"merklePath"`
	Root       RootNode          `json:"root"`
}

type PathNode struct {
	Hash     string            `json:"hash"`
	Balances map[string]string `json:"balances"`
	Level    int               `json:"level"`
	Position string            `json:"position"`
}

type RootNode struct {
	Hash     string            `json:"hash"`
	Balances map[string]string `json:"balances"`
	Level    int               `json:"level"`
}

func ValidateProof(raw []byte) (bool, error) {
	var p Proof
	if err := json.Unmarshal(raw, &p); err != nil {
		return false, err
	}
	return validate(p), nil
}

func RootHash(raw []byte) (string, error) {
	var p Proof
	if err := json.Unmarshal(raw, &p); err != nil {
		return "", err
	}
	if p.Root.Hash == "" {
		return "", fmt.Errorf("proof missing root.hash")
	}
	return p.Root.Hash, nil
}

func AuditID(raw []byte) (string, error) {
	var p Proof
	if err := json.Unmarshal(raw, &p); err != nil {
		return "", err
	}
	return p.AuditID, nil
}

func validate(p Proof) bool {
	if p.EncryptUID == "" || p.Nonce == "" || len(p.Balances) == 0 {
		return false
	}
	path := p.MerkelPath
	if len(path) == 0 {
		path = p.MerklePath
	}
	if len(path) == 0 || p.Root.Hash == "" {
		return false
	}

	currentBalances, ok := normalizeBalances(p.Balances)
	if !ok {
		return false
	}
	currentHash := p.Hash
	if currentHash == "" {
		currentHash = leafHash(p.EncryptUID, p.Nonce, currentBalances)
	} else if currentHash != leafHash(p.EncryptUID, p.Nonce, currentBalances) {
		return false
	}

	for _, sibling := range path {
		sibBalances, ok := normalizeBalances(sibling.Balances)
		if !ok || sibling.Hash == "" || sibling.Level <= 0 {
			return false
		}
		combined := addBalances(currentBalances, sibBalances)
		pos := strings.ToLower(strings.TrimSpace(sibling.Position))
		switch pos {
		case "left", "l":
			currentHash = internalHash(sibling.Hash, currentHash, combined, sibling.Level)
		default:
			currentHash = internalHash(currentHash, sibling.Hash, combined, sibling.Level)
		}
		currentBalances = combined
	}

	rootBalances, ok := normalizeBalances(p.Root.Balances)
	if !ok {
		return false
	}
	return currentHash == p.Root.Hash && balancesEqual(currentBalances, rootBalances)
}

func leafHash(encryptUID, nonce string, balances map[string]string) string {
	return truncatedSHA256(encryptUID + "," + nonce + "," + jsonBalances(balances))
}

func internalHash(leftHash, rightHash string, balances map[string]string, level int) string {
	return truncatedSHA256(leftHash + rightHash + "," + jsonBalances(balances) + "," + fmt.Sprintf("%d", level))
}

func truncatedSHA256(input string) string {
	sum := sha256.Sum256([]byte(input))
	return hex.EncodeToString(sum[:])[:HashHexLen]
}

func normalizeBalances(in map[string]string) (map[string]string, bool) {
	out := make(map[string]string, len(in))
	for k, v := range in {
		d, err := decimal.NewFromString(v)
		if err != nil {
			return nil, false
		}
		out[k] = d.String()
	}
	return out, true
}

func addBalances(a, b map[string]string) map[string]string {
	out := make(map[string]string, len(a)+len(b))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		left, _ := decimal.NewFromString(out[k])
		right, _ := decimal.NewFromString(v)
		out[k] = left.Add(right).String()
	}
	return out
}

func balancesEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, av := range a {
		if b[k] != av {
			return false
		}
	}
	return true
}

func jsonBalances(balances map[string]string) string {
	keys := make([]string, 0, len(balances))
	for k := range balances {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%q:%q", k, balances[k]))
	}
	return "{" + strings.Join(parts, ",") + "}"
}
