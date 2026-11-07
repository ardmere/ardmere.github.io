// Package merkle validates Ourbit user Merkle proofs.
// Logic ported from github.com/OurbitOfficial/merkle_validate
// (SHA-256 truncated to 32 hex chars / 128 bits).
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

const (
	HashHexLen = 32

	roleEmpty = 0
	roleLeft  = 1
	roleRight = 2
	roleRoot  = 3
)

// Proof is the merkel_tree_ourbit.json shape from Ourbit (self + path).
type Proof struct {
	Self          TreeNode   `json:"self"`
	Path          []TreeNode `json:"path"`
	AllCurrencies []string   `json:"allCurrencies"`
}

// TreeNode mirrors com.ourbit.model.TreeNode.
type TreeNode struct {
	SnapshotDate string                     `json:"snapshotDate"`
	Balances     map[string]decimal.Decimal `json:"balances"`
	Hash         string                     `json:"hash"`
	Level        int                        `json:"level"`
	Role         int                        `json:"role"`
	UID          string                     `json:"uid"`

	leftChildHashId  string
	rightChildHashId string
}

func (n *TreeNode) UnmarshalJSON(raw []byte) error {
	var aux struct {
		SnapshotDate string          `json:"snapshotDate"`
		Balances     json.RawMessage `json:"balances"`
		Hash         *string         `json:"hash"`
		Level        int             `json:"level"`
		Role         int             `json:"role"`
		UID          *string         `json:"uid"`
	}
	if err := json.Unmarshal(raw, &aux); err != nil {
		return err
	}
	n.SnapshotDate = aux.SnapshotDate
	n.Level = aux.Level
	n.Role = aux.Role
	if aux.Hash != nil {
		n.Hash = *aux.Hash
	}
	if aux.UID != nil {
		n.UID = *aux.UID
	}
	n.Balances = map[string]decimal.Decimal{}
	if len(aux.Balances) > 0 && string(aux.Balances) != "null" {
		var rawMap map[string]json.RawMessage
		if err := json.Unmarshal(aux.Balances, &rawMap); err != nil {
			return err
		}
		for k, v := range rawMap {
			d, err := parseAmount(v)
			if err != nil {
				return fmt.Errorf("balance %s: %w", k, err)
			}
			n.Balances[k] = d
		}
	}
	return nil
}

func parseAmount(raw json.RawMessage) (decimal.Decimal, error) {
	raw = json.RawMessage(strings.TrimSpace(string(raw)))
	if len(raw) == 0 || string(raw) == "null" {
		return decimal.Zero, nil
	}
	if raw[0] == '"' {
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return decimal.Zero, err
		}
		if strings.TrimSpace(s) == "" {
			return decimal.Zero, nil
		}
		return decimal.NewFromString(s)
	}
	// Keep exact decimal text from JSON number tokens.
	return decimal.NewFromString(string(raw))
}

func ValidateProof(raw []byte) (bool, error) {
	var p Proof
	if err := json.Unmarshal(raw, &p); err != nil {
		return false, err
	}
	return validate(p)
}

func RootHash(raw []byte) (string, error) {
	var p Proof
	if err := json.Unmarshal(raw, &p); err != nil {
		return "", err
	}
	if len(p.Path) == 0 {
		return "", fmt.Errorf("proof missing path")
	}
	root := p.Path[len(p.Path)-1]
	if root.Hash == "" {
		return "", fmt.Errorf("proof missing root hash")
	}
	return root.Hash, nil
}

func ProofID(raw []byte) (string, error) {
	var p Proof
	if err := json.Unmarshal(raw, &p); err != nil {
		return "", err
	}
	if p.Self.SnapshotDate == "" {
		return "", nil
	}
	return "PR-" + p.Self.SnapshotDate, nil
}

func validate(p Proof) (bool, error) {
	if p.Self.Hash == "" || len(p.Path) == 0 {
		return false, fmt.Errorf("invalid merkleTree content")
	}
	enabled := enabledCurrencies(p)
	if !validateSelf(p.Self, enabled) {
		return false, nil
	}
	if len(p.Path) == 0 || !validatePathNode(p.Path[0]) {
		return false, nil
	}
	if p.Path[0].Role == p.Self.Role {
		return false, nil
	}
	newRoot := buildMerkelTreeRoot(p.Path, p.Self, enabled)
	oldRoot := p.Path[len(p.Path)-1]
	if newRoot.Hash != oldRoot.Hash {
		return false, nil
	}
	if newRoot.Level != oldRoot.Level {
		return false, nil
	}
	return balancesEqual(newRoot.Balances, oldRoot.Balances), nil
}

func enabledCurrencies(p Proof) []string {
	if len(p.AllCurrencies) > 0 {
		out := append([]string{}, p.AllCurrencies...)
		sort.Strings(out)
		return out
	}
	seen := map[string]struct{}{}
	collect := func(b map[string]decimal.Decimal) {
		for k := range b {
			seen[strings.Split(k, ":")[0]] = struct{}{}
		}
	}
	collect(p.Self.Balances)
	for _, n := range p.Path {
		collect(n.Balances)
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func validateSelf(self TreeNode, enabled []string) bool {
	calc := createMerkelNodeHash(self, enabled)
	return calc == self.Hash
}

func validatePathNode(n TreeNode) bool {
	if n.Hash == "" || n.Role < 0 || n.Role > 3 || n.Level < 0 {
		return false
	}
	if len(n.Balances) == 0 {
		return false
	}
	return true
}

func buildMerkelTreeRoot(path []TreeNode, self TreeNode, enabled []string) TreeNode {
	if len(path) <= 1 {
		return TreeNode{}
	}
	parent := createParentTreeNode(path[0], self, enabled)
	for i := 1; i < len(path)-1; i++ {
		self = parent
		parent = createParentTreeNode(path[i], self, enabled)
	}
	parent.Role = roleRoot
	return parent
}

func createParentTreeNode(friend, self TreeNode, enabled []string) TreeNode {
	if friend.Role == roleLeft {
		return constructInternalNode(friend, self, enabled)
	}
	return constructInternalNode(self, friend, enabled)
}

func constructInternalNode(left, right TreeNode, enabled []string) TreeNode {
	if right.Hash == "" && len(right.Balances) == 0 && right.Role == roleEmpty {
		// keep as-is; empty path placeholders use blank sibling hash
	}
	parent := TreeNode{
		Balances: map[string]decimal.Decimal{},
	}
	mergeAsset(&parent, left)
	mergeAsset(&parent, right)
	parent.leftChildHashId = left.Hash
	parent.rightChildHashId = right.Hash
	parent.Level = left.Level + 1
	parent.SnapshotDate = left.SnapshotDate
	parent.Hash = createMerkelNodeHash(parent, enabled)
	return parent
}

func mergeAsset(dst *TreeNode, src TreeNode) {
	for k, v := range src.Balances {
		dst.Balances[k] = dst.Balances[k].Add(v)
	}
}

func balancesEqual(a, b map[string]decimal.Decimal) bool {
	if len(a) != len(b) {
		// Java compares only keys present on newRoot; tolerate extras on either side if values match where present.
	}
	for k, v := range a {
		ov, ok := b[k]
		if !ok || !v.Equal(ov) {
			return false
		}
	}
	return true
}

func createMerkelNodeHash(node TreeNode, enabled []string) string {
	isLeaf := node.Level == 0
	balMap := balanceMapForCalcHash(node.Balances, isLeaf, enabled)
	balanceSb := balanceString(balMap)
	var content string
	if node.UID != "" {
		content = "PR-" + node.SnapshotDate + "," + node.UID + "," + balanceSb
	} else {
		right := node.rightChildHashId
		content = node.leftChildHashId + "," + right + "," + balanceSb
	}
	return truncatedSHA256(content)
}

func balanceMapForCalcHash(balances map[string]decimal.Decimal, isLeaf bool, enabled []string) map[string]decimal.Decimal {
	out := map[string]decimal.Decimal{}
	for _, c := range enabled {
		out[c] = decimal.Zero
	}
	for k, v := range balances {
		currency := strings.Split(k, ":")[0]
		out[currency] = out[currency].Add(v)
	}
	if isLeaf {
		for k, v := range out {
			if v.IsZero() {
				delete(out, k)
			}
		}
	}
	return out
}

func balanceString(balances map[string]decimal.Decimal) string {
	keys := make([]string, 0, len(balances))
	for k := range balances {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+":"+plainAmount(balances[k]))
	}
	return strings.Join(parts, ",")
}

func plainAmount(d decimal.Decimal) string {
	// Match Java BigDecimal.stripTrailingZeros().toPlainString()
	s := d.String()
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	if s == "" || s == "-0" {
		return "0"
	}
	return s
}

func truncatedSHA256(input string) string {
	sum := sha256.Sum256([]byte(input))
	return hex.EncodeToString(sum[:])[:HashHexLen]
}
