// Package merkle validates Toobit Summation Merkle Tree (SMT) user proofs.
// Logic ported from github.com/toobit-docs/por-verifier (SHA256-SMT with
// SMT_LEAF / SMT_NODE domain separators).
package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"
)

const (
	leafPrefix = "SMT_LEAF"
	nodePrefix = "SMT_NODE"
)

// ProofFile is the offline proof JSON downloaded from Toobit (or demo testdata).
type ProofFile struct {
	Version    string         `json:"version"`
	AuditID    string         `json:"audit_id"`
	Algorithm  string         `json:"algorithm"`
	Precisions map[string]int `json:"precisions"`
	Self       SelfData       `json:"self"`
	Path       []PathNode     `json:"path"`
	Root       RootData       `json:"root"`
}

// SelfData contains the user's own leaf data.
type SelfData struct {
	Nonce    string            `json:"nonce"`
	Balances map[string]string `json:"balances"`
	LeafHash string            `json:"leaf_hash"`
}

// PathNode is a sibling on the Merkle path.
type PathNode struct {
	Position string            `json:"position"`
	Hash     string            `json:"hash"`
	Balances map[string]string `json:"balances"`
}

// RootData holds the published root.
type RootData struct {
	Hash          string            `json:"hash"`
	TotalBalances map[string]string `json:"total_balances"`
}

type balances map[string]*big.Int

// ValidateProof recomputes leaf + path per toobit-docs/por-verifier.
func ValidateProof(raw []byte) (bool, error) {
	p, err := parseProofFile(raw)
	if err != nil {
		return false, err
	}
	res, err := verify(p)
	if err != nil {
		return false, err
	}
	return res.valid, nil
}

// RootHash returns the published root hash from the proof file.
func RootHash(raw []byte) (string, error) {
	p, err := parseProofFile(raw)
	if err != nil {
		return "", err
	}
	root := strings.ToLower(strings.TrimSpace(p.Root.Hash))
	if root == "" {
		return "", fmt.Errorf("proof missing root.hash")
	}
	return root, nil
}

// AuditID returns the audit_id from the proof file.
func AuditID(raw []byte) (string, error) {
	p, err := parseProofFile(raw)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(p.AuditID), nil
}

func parseProofFile(raw []byte) (ProofFile, error) {
	var p ProofFile
	if err := json.Unmarshal(raw, &p); err != nil {
		return ProofFile{}, err
	}
	if p.Self.Nonce == "" || p.Self.LeafHash == "" || p.Root.Hash == "" {
		return ProofFile{}, fmt.Errorf("invalid toobit proof content")
	}
	return p, nil
}

type verifyResult struct {
	valid bool
}

func verify(p ProofFile) (verifyResult, error) {
	userBal, err := parseBalances(p.Self.Balances)
	if err != nil {
		return verifyResult{}, err
	}
	if !isValidSHA256Hex(p.Self.LeafHash) {
		return verifyResult{}, fmt.Errorf("invalid leaf hash format")
	}
	if !isValidSHA256Hex(p.Root.Hash) {
		return verifyResult{}, fmt.Errorf("invalid root hash format")
	}

	computedLeaf := computeLeafHash(p.Self.Nonce, userBal)
	leafOK := computedLeaf == strings.ToLower(strings.TrimSpace(p.Self.LeafHash))

	currentHash := strings.ToLower(strings.TrimSpace(p.Self.LeafHash))
	currentBal := cloneBalances(userBal)
	sumOK := true

	for i, sib := range p.Path {
		if !isValidSHA256Hex(sib.Hash) {
			return verifyResult{}, fmt.Errorf("invalid hash at path[%d]", i)
		}
		sibBal, err := parseBalances(sib.Balances)
		if err != nil {
			return verifyResult{}, fmt.Errorf("path[%d] balances: %w", i, err)
		}
		for asset, amount := range sibBal {
			if amount.Sign() < 0 {
				sumOK = false
				_ = asset
			}
		}
		sumBal := sumBalances(currentBal, sibBal)

		var leftHash, rightHash string
		switch strings.ToLower(strings.TrimSpace(sib.Position)) {
		case "left":
			leftHash = strings.ToLower(strings.TrimSpace(sib.Hash))
			rightHash = currentHash
		case "right":
			leftHash = currentHash
			rightHash = strings.ToLower(strings.TrimSpace(sib.Hash))
		default:
			return verifyResult{}, fmt.Errorf("path[%d]: invalid position %q", i, sib.Position)
		}
		currentHash = computeParentHash(leftHash, rightHash, sumBal)
		currentBal = sumBal
	}

	root := strings.ToLower(strings.TrimSpace(p.Root.Hash))
	pathOK := currentHash == root

	if len(p.Root.TotalBalances) > 0 {
		rootSums, err := parseBalances(p.Root.TotalBalances)
		if err != nil {
			return verifyResult{}, fmt.Errorf("root totals: %w", err)
		}
		for asset, expected := range rootSums {
			if got, ok := currentBal[asset]; !ok || got.Cmp(expected) != 0 {
				sumOK = false
			}
		}
	}

	return verifyResult{valid: leafOK && pathOK && sumOK}, nil
}

func computeLeafHash(nonce string, bal balances) string {
	keys := sortedAssets(bal)
	var b strings.Builder
	b.WriteString(leafPrefix)
	b.WriteString("|")
	b.WriteString(nonce)
	for _, asset := range keys {
		b.WriteString(fmt.Sprintf("|%s:%s", asset, bal[asset].String()))
	}
	sum := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(sum[:])
}

func computeParentHash(leftHash, rightHash string, sumBal balances) string {
	keys := sortedAssets(sumBal)
	var b strings.Builder
	b.WriteString(nodePrefix)
	b.WriteString("|")
	b.WriteString(leftHash)
	b.WriteString("|")
	b.WriteString(rightHash)
	for _, asset := range keys {
		b.WriteString(fmt.Sprintf("|%s:%s", asset, sumBal[asset].String()))
	}
	sum := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(sum[:])
}

func parseBalances(raw map[string]string) (balances, error) {
	out := make(balances, len(raw))
	for asset, amountStr := range raw {
		amount, ok := new(big.Int).SetString(strings.TrimSpace(amountStr), 10)
		if !ok {
			return nil, fmt.Errorf("invalid amount for %s: %q", asset, amountStr)
		}
		out[asset] = amount
	}
	return out, nil
}

func sortedAssets(bal balances) []string {
	keys := make([]string, 0, len(bal))
	for k := range bal {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func cloneBalances(bal balances) balances {
	out := make(balances, len(bal))
	for k, v := range bal {
		out[k] = new(big.Int).Set(v)
	}
	return out
}

func sumBalances(a, b balances) balances {
	out := cloneBalances(a)
	for asset, amount := range b {
		if existing, ok := out[asset]; ok {
			out[asset] = new(big.Int).Add(existing, amount)
		} else {
			out[asset] = new(big.Int).Set(amount)
		}
	}
	return out
}

func isValidSHA256Hex(s string) bool {
	s = strings.TrimSpace(s)
	if len(s) != 64 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}
