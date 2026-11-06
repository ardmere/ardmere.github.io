// Package merkle validates CoinEx user Merkle inclusion proofs.
// Logic ported from github.com/coinexcom/proof-of-reserves (Go; full SHA-256).
package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

type Pos string

const (
	Left  Pos = "left"
	Right Pos = "right"
)

// Balance matches the official CoinEx proof JSON balances object.
type Balance map[string]decimal.Decimal

func (b Balance) Format() string {
	v, _ := json.Marshal(b) // encoding/json sorts map keys; compact, no spaces
	return string(v)
}

func (b Balance) Add(other Balance) Balance {
	r := make(Balance, len(b)+len(other))
	for k, v := range b {
		r[k] = v
	}
	for k, v := range other {
		if vv, ok := r[k]; ok {
			r[k] = v.Add(vv)
		} else {
			r[k] = v
		}
	}
	return r
}

func (b Balance) Equal(other Balance) bool {
	if len(b) != len(other) {
		return false
	}
	for k, v := range b {
		vv, ok := other[k]
		if !ok || !v.Equal(vv) {
			return false
		}
	}
	return true
}

func hash(v string) string {
	sum := sha256.Sum256([]byte(v))
	return hex.EncodeToString(sum[:])
}

// Proof is the login-gated "Copy my audition" JSON from CoinEx.
type Proof struct {
	Root struct {
		Balances Balance `json:"balances"`
		Hash     string  `json:"hash"`
	} `json:"root"`
	Self struct {
		Balances Balance `json:"balances"`
		Nonce    string  `json:"nonce"`
	} `json:"self"`
	Path []struct {
		Balances Balance `json:"balances"`
		Hash     string  `json:"hash"`
		Pos      Pos     `json:"pos"`
	} `json:"path"`
}

// ValidateProof recomputes the leaf→root path per coinexcom/proof-of-reserves.
func ValidateProof(raw []byte) (bool, error) {
	var m Proof
	if err := json.Unmarshal(raw, &m); err != nil {
		return false, err
	}
	if m.Root.Hash == "" || len(m.Path) == 0 {
		return false, fmt.Errorf("empty merkle proof file")
	}
	if strings.TrimSpace(m.Self.Nonce) == "" {
		return false, fmt.Errorf("proof missing self.nonce")
	}
	h := hash(m.Self.Nonce + m.Self.Balances.Format())
	b := m.Self.Balances
	for _, path := range m.Path {
		if path.Hash == "" { // padding node: sibling hash = self hash, zero balances
			h = hash(h + h + b.Format())
			continue
		}
		b = b.Add(path.Balances)
		if path.Pos == Left {
			h = hash(path.Hash + h + b.Format())
		} else {
			// official CLI: any non-left (including "right") concatenates self||sibling
			h = hash(h + path.Hash + b.Format())
		}
	}
	if !strings.EqualFold(h, strings.TrimSpace(m.Root.Hash)) {
		return false, nil
	}
	if !b.Equal(m.Root.Balances) {
		return false, nil
	}
	return true, nil
}

// RootHash returns the published root hash from the proof file.
func RootHash(raw []byte) (string, error) {
	var m Proof
	if err := json.Unmarshal(raw, &m); err != nil {
		return "", err
	}
	return strings.ToLower(strings.TrimPrefix(strings.TrimSpace(m.Root.Hash), "0x")), nil
}
