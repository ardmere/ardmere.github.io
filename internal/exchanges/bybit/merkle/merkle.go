// Package merkle validates Bybit user Merkle proofs (Gen 1.5 SHA-256 binary tree).
// Logic ported from github.com/bybit-exchange/merkle-proof GenericMerkleTree + Balance40V5.
package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/shopspring/decimal"
)

// Proof is the myProof.json shape from Bybit (self + path siblings).
type Proof struct {
	Self SelfNode   `json:"self"`
	Path []PathNode `json:"path"`
}

type SelfNode struct {
	AuditID  string            `json:"auditId"`
	Type     int               `json:"type"`
	Hash     string            `json:"hash"`
	Height   int               `json:"height"`
	Nonce    string            `json:"nonce"`
	UserHash string            `json:"userHash"`
	Balance  map[string]string `json:"balance"`
}

type PathNode struct {
	Type    int               `json:"type"`
	Hash    string            `json:"hash"`
	Height  int               `json:"height"`
	Nonce   string            `json:"nonce"`
	Balance map[string]string `json:"balance"`
}

type balanceV5 map[string]string

func (b balanceV5) validate() bool {
	for _, asset := range Asset40V5 {
		v, ok := b[asset]
		if !ok {
			return false
		}
		if _, err := decimal.NewFromString(v); err != nil {
			return false
		}
	}
	return true
}

func (b balanceV5) concat() string {
	out := ""
	for _, asset := range Asset40V5 {
		out += b[asset]
	}
	return out
}

func (b balanceV5) add(other balanceV5) balanceV5 {
	res := make(balanceV5, len(Asset40V5))
	for _, asset := range Asset40V5 {
		a, _ := decimal.NewFromString(b[asset])
		c, _ := decimal.NewFromString(other[asset])
		res[asset] = a.Add(c).RoundDown(balanceScale).StringFixed(balanceScale)
	}
	return res
}

func (b balanceV5) equal(other balanceV5) bool {
	for _, asset := range Asset40V5 {
		if b[asset] != other[asset] {
			return false
		}
	}
	return true
}

func toBalanceV5(m map[string]string) balanceV5 {
	b := make(balanceV5, len(Asset40V5))
	for _, asset := range Asset40V5 {
		b[asset] = m[asset]
	}
	return b
}

func sha256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func (s SelfNode) validateSelf() bool {
	if s.UserHash == "" {
		return false
	}
	bal := toBalanceV5(s.Balance)
	if !bal.validate() {
		return false
	}
	return sha256Hex(s.UserHash+bal.concat()) == s.Hash
}

func (p PathNode) validateNode() bool {
	bal := toBalanceV5(p.Balance)
	if !bal.validate() || p.Hash == "" || p.Type < 1 || p.Type > 4 || p.Height < 1 {
		return false
	}
	return true
}

func combinePath(leftHash, rightHash string, b1, b2 balanceV5, height, nodeType int) (hash string, bal balanceV5, h int, t int) {
	if !b1.validate() || !b2.validate() {
		return "", nil, 0, 0
	}
	bal = b1.add(b2)
	data := leftHash + rightHash + bal.concat() + fmt.Sprintf("%d", height)
	return sha256Hex(data), bal, height, nodeType
}

func siblingType(sib PathNode) int {
	if sib.Type == 3 {
		return 3
	}
	if sib.Type == 1 {
		return 2
	}
	return 1
}

// ValidateProofV5 checks a v5-schema user proof (40 online assets).
func ValidateProofV5(raw []byte) (bool, error) {
	var p Proof
	if err := json.Unmarshal(raw, &p); err != nil {
		return false, err
	}
	return validateTree(p), nil
}

func validateTree(p Proof) bool {
	if p.Self.UserHash == "" || len(p.Path) == 0 {
		return false
	}
	if !p.Self.validateSelf() || !p.Path[0].validateNode() {
		return false
	}
	if p.Path[0].Type == p.Self.Type {
		return false
	}

	selfBal := toBalanceV5(p.Self.Balance)
	path0Bal := toBalanceV5(p.Path[0].Balance)

	var left, right string
	if p.Self.Type == 1 {
		left, right = p.Self.Hash, p.Path[0].Hash
	} else {
		left, right = p.Path[0].Hash, p.Self.Hash
	}

	height := 1
	var nodeType int
	if len(p.Path) > 1 {
		nodeType = siblingType(p.Path[1])
	}
	nodeHash, nodeBal, nodeH, nodeT := combinePath(left, right, selfBal, path0Bal, height+1, nodeType)
	if nodeHash == "" {
		return false
	}
	node := struct {
		hash string
		bal  balanceV5
		h    int
		t    int
	}{nodeHash, nodeBal, nodeH, nodeT}

	for i := 1; i < len(p.Path)-1; i++ {
		height++
		if !p.Path[i].validateNode() {
			return false
		}
		sibBal := toBalanceV5(p.Path[i].Balance)
		var l, r string
		if p.Path[i].Type == 1 {
			l, r = p.Path[i].Hash, node.hash
		} else {
			l, r = node.hash, p.Path[i].Hash
		}
		nt := siblingType(p.Path[i+1])
		nh, nb, nhH, nhT := combinePath(l, r, node.bal, sibBal, height+1, nt)
		if nh == "" {
			return false
		}
		node = struct {
			hash string
			bal  balanceV5
			h    int
			t    int
		}{nh, nb, nhH, nhT}
	}

	root := p.Path[len(p.Path)-1]
	rootBal := toBalanceV5(root.Balance)
	return node.hash == root.Hash && node.bal.equal(rootBal) && node.h == root.Height
}

// RootHash returns the Merkle root hash from a valid proof path (last path node).
func RootHash(raw []byte) (string, error) {
	var p Proof
	if err := json.Unmarshal(raw, &p); err != nil {
		return "", err
	}
	if len(p.Path) == 0 {
		return "", fmt.Errorf("proof missing path")
	}
	return p.Path[len(p.Path)-1].Hash, nil
}

// AuditID returns the audit id from the self node when present.
func AuditID(raw []byte) (string, error) {
	var p Proof
	if err := json.Unmarshal(raw, &p); err != nil {
		return "", err
	}
	return p.Self.AuditID, nil
}
