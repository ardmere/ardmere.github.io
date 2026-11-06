// Package merkle validates Bitunix user Merkle proofs.
// Logic ported from github.com/BitunixOfficial/proof-of-reserves
// (full SHA-256; self + path TreeNode with height-descending aggregation).
package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

const (
	typeLeft  = 1
	typeRight = 2
	typeRoot  = 3
)

// Proof is the merkle_tree_bitunix.json shape (self + path).
type Proof struct {
	Self TreeNode   `json:"self"`
	Path []TreeNode `json:"path"`
}

// TreeNode mirrors com.bitunix.model.TreeNode.
type TreeNode struct {
	EncryptUID string            `json:"encryptUid"`
	MerkleLeaf string            `json:"merkleLeaf"`
	Nonce      string            `json:"nonce"`
	Balances   map[string]string `json:"balances"`
	Height     int               `json:"height"`
	Type       int               `json:"type"`
	Precisions map[string]int    `json:"precisions"`
}

type node struct {
	encryptUID string
	merkleLeaf string
	nonce      string
	balances   map[string]string
	height     int
	typ        int
	precisions map[string]int
}

// ValidateProof recomputes the leaf→root path per BitunixOfficial/proof-of-reserves.
func ValidateProof(raw []byte) (bool, error) {
	var p Proof
	if err := json.Unmarshal(raw, &p); err != nil {
		return false, err
	}
	return validate(p)
}

// RootHash returns the Merkle root hash from the last path node.
func RootHash(raw []byte) (string, error) {
	var p Proof
	if err := json.Unmarshal(raw, &p); err != nil {
		return "", err
	}
	if len(p.Path) == 0 {
		return "", fmt.Errorf("proof missing path")
	}
	root := strings.ToLower(strings.TrimSpace(p.Path[len(p.Path)-1].MerkleLeaf))
	if root == "" {
		return "", fmt.Errorf("proof missing root merkleLeaf")
	}
	return root, nil
}

func validate(p Proof) (bool, error) {
	if len(p.Path) < 2 {
		return false, fmt.Errorf("Must be at least two leafs to construct a Merkle tree")
	}
	if p.Self.MerkleLeaf == "" || p.Self.EncryptUID == "" || p.Self.Nonce == "" {
		return false, fmt.Errorf("invalid merkle tree content")
	}
	if !validateSelf(p.Self) {
		return false, nil
	}
	if !validatePathNode(p.Path[0]) {
		return false, nil
	}
	if p.Path[0].Type == p.Self.Type {
		return false, nil
	}

	precisions := p.Self.Precisions
	parent := createParent(toNode(p.Path[0]), toNode(p.Self), precisions)
	cur := toNode(p.Self)
	for i := 1; i < len(p.Path)-1; i++ {
		cur = parent
		parent = createParent(toNode(p.Path[i]), cur, precisions)
	}
	oldRoot := toNode(p.Path[len(p.Path)-1])
	return parent.merkleLeaf == oldRoot.merkleLeaf, nil
}

func validateSelf(self TreeNode) bool {
	if self.Balances == nil || self.Precisions == nil || self.Height == 0 {
		return false
	}
	if len(self.Balances) == 0 || len(self.Precisions) == 0 {
		return false
	}
	return createMerkelNodeLeaf(toNode(self)) == self.MerkleLeaf
}

func validatePathNode(n TreeNode) bool {
	if n.MerkleLeaf == "" || n.Height == 0 || n.Type < typeLeft || n.Type > typeRoot {
		return false
	}
	if n.Balances == nil || len(n.Balances) == 0 {
		return false
	}
	for _, v := range n.Balances {
		if v == "" {
			return false
		}
	}
	return true
}

func toNode(t TreeNode) node {
	return node{
		encryptUID: t.EncryptUID,
		merkleLeaf: t.MerkleLeaf,
		nonce:      t.Nonce,
		balances:   copyBalances(t.Balances),
		height:     t.Height,
		typ:        t.Type,
		precisions: copyPrecisions(t.Precisions),
	}
}

func createParent(sibling, self node, precisions map[string]int) node {
	if sibling.typ == typeLeft {
		return constructInternal(sibling, self, precisions)
	}
	return constructInternal(self, sibling, precisions)
}

func constructInternal(left, right node, precisions map[string]int) node {
	parent := node{
		balances: getTotalBalances(left.balances, right.balances, precisions),
		height:   left.height - 1,
	}
	parent.merkleLeaf = createMerkelParentLeaf(left, right, parent)
	return parent
}

func getTotalBalances(left, right map[string]string, precisions map[string]int) map[string]string {
	out := make(map[string]string, len(left))
	for k, v := range left {
		scale := precisions[k]
		sum := parseDec(v).Add(parseDec(right[k])).RoundDown(int32(scale))
		out[k] = sum.StringFixed(int32(scale))
	}
	return out
}

func createMerkelNodeLeaf(n node) string {
	// Fastjson Map serialization is compact with keys in a stable order;
	// encoding/json.Marshal sorts map keys alphabetically (matches sample).
	content := n.encryptUID + "," + n.nonce + "," + balancesJSON(n.balances)
	return sha256Hex(content)
}

func createMerkelParentLeaf(left, right, parent node) string {
	content := left.merkleLeaf + right.merkleLeaf + "," + balancesJSON(parent.balances) + "," + strconv.Itoa(parent.height)
	return sha256Hex(content)
}

func balancesJSON(balances map[string]string) string {
	raw, err := json.Marshal(balances)
	if err != nil {
		return "{}"
	}
	return string(raw)
}

func sha256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func parseDec(s string) decimal.Decimal {
	if strings.TrimSpace(s) == "" {
		return decimal.Zero
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero
	}
	return d
}

func copyBalances(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func copyPrecisions(in map[string]int) map[string]int {
	out := make(map[string]int, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
