// Package merkle validates MEXC user Merkle proofs.
// Logic ported from github.com/mexcmerkle/mexc-merkle-proof-of-reserves (SHA-256 truncated to 32 hex chars).
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

	typeLeft  = 1
	typeRight = 2
	typeRoot  = 3
)

// Proof is the myProof.json shape from MEXC (self + pathList).
type Proof struct {
	Self     SelfNode   `json:"self"`
	PathList []PathNode `json:"pathList"`
}

type SelfNode struct {
	Level     int               `json:"level"`
	Type      int               `json:"type"`
	Hash      string            `json:"hash"`
	Balances  map[string]string `json:"balances"`
	ProofID   string            `json:"proofId"`
	AccountID string            `json:"accountId"`
}

type PathNode struct {
	Level    int               `json:"level"`
	Type     int               `json:"type"`
	Hash     string            `json:"hash"`
	Balances map[string]string `json:"balances"`
}

type node struct {
	level    int
	typ      int
	hash     string
	balances map[string]string
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
	if len(p.PathList) == 0 {
		return "", fmt.Errorf("proof missing pathList")
	}
	root := p.PathList[len(p.PathList)-1]
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
	return p.Self.ProofID, nil
}

func validate(p Proof) (bool, error) {
	if len(p.PathList) == 0 || p.Self.Hash == "" || p.Self.ProofID == "" || p.Self.AccountID == "" {
		return false, fmt.Errorf("invalid merkleTree content")
	}
	if !validateSelf(p.Self) {
		return false, nil
	}
	return pathValidate(p.PathList, p.Self), nil
}

func validateSelf(self SelfNode) bool {
	if self.Balances == nil {
		return false
	}
	content := self.ProofID + "," + self.AccountID + "," + balanceStringRaw(self.Balances)
	return truncatedSHA256(content) == self.Hash
}

func pathValidate(path []PathNode, self SelfNode) bool {
	root := path[len(path)-1]
	depth := root.Level
	levelMap := make(map[int]PathNode, len(path))
	for _, n := range path {
		levelMap[n.Level] = n
	}

	current := node{
		level:    self.Level,
		typ:      self.Type,
		hash:     self.Hash,
		balances: copyBalances(self.Balances),
	}
	var left, right node
	havePair := false
	for i := 0; i < depth; i++ {
		sibling, ok := levelMap[i]
		var sib node
		if !ok {
			sibTyp := typeLeft
			if current.typ == typeLeft {
				sibTyp = typeRight
			}
			sib = emptyNode(sibTyp, i)
		} else {
			sib = pathToNode(sibling)
		}
		left, right = decidePair(current, sib)
		havePair = true
		parentLevel := i + 1
		parentTyp := parentType(levelMap, parentLevel)
		current = createInner(left, right, parentTyp, parentLevel)
	}
	if !havePair {
		return false
	}
	clcRoot := createInner(left, right, typeRoot, root.Level)
	return clcRoot.hash == root.Hash
}

func parentType(levelMap map[int]PathNode, parentLevel int) int {
	parentBrother, ok := levelMap[parentLevel]
	if !ok {
		return typeLeft
	}
	if parentBrother.Type == typeLeft {
		return typeRight
	}
	return typeLeft
}

func decidePair(current, brother node) (left, right node) {
	if current.typ == typeLeft {
		return current, brother
	}
	return brother, current
}

func emptyNode(typ, level int) node {
	return node{typ: typ, level: level, hash: "", balances: map[string]string{}}
}

func pathToNode(p PathNode) node {
	return node{
		level:    p.Level,
		typ:      p.Type,
		hash:     p.Hash,
		balances: copyBalances(p.Balances),
	}
}

func createInner(left, right node, typ, level int) node {
	out := node{
		level:    level,
		typ:      typ,
		balances: map[string]string{},
	}
	mergeBalance(&out, left)
	if strings.TrimSpace(right.hash) != "" {
		initBalanceData(&out)
		mergeBalance(&out, right)
	}
	out.hash = truncatedSHA256(left.hash + "," + right.hash + "," + balanceStringPlain(out.balances))
	return out
}

func initBalanceData(n *node) {
	for _, asset := range []string{"BTC", "USDT", "ETH", "USDC"} {
		if _, ok := n.balances[asset]; !ok {
			n.balances[asset] = "0"
		}
	}
}

func mergeBalance(dst *node, src node) {
	for k, v := range src.balances {
		sum := parseDec(dst.balances[k]).Add(parseDec(v))
		dst.balances[k] = sum.String()
	}
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

func balanceStringRaw(balances map[string]string) string {
	keys := sortedKeys(balances)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+":"+balances[k])
	}
	return strings.Join(parts, ",")
}

func balanceStringPlain(balances map[string]string) string {
	keys := sortedKeys(balances)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+":"+plainAmount(balances[k]))
	}
	return strings.Join(parts, ",")
}

func plainAmount(s string) string {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return s
	}
	return d.String() // shopspring strips trailing zeros by default in String()
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func copyBalances(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func truncatedSHA256(input string) string {
	sum := sha256.Sum256([]byte(input))
	return hex.EncodeToString(sum[:])[:HashHexLen]
}
