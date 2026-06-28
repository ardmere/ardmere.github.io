package verifier

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type okxProofEnvelope struct {
	Proof        json.RawMessage `json:"proof"`
	PublicInputs []json.Number   `json:"public_inputs"`
}

type okxZKAggregates struct {
	Equity    uint64
	Debt      uint64
	Liability uint64
	MerkleHex string
}

// okxZKAggregatesFromProof parses Plonky2 public inputs from sum_proof_data.json.
// Layout: [equity, debt, root_hash x4 u64 limbs LE].
func okxZKAggregatesFromProof(env okxProofEnvelope) (okxZKAggregates, error) {
	if len(env.PublicInputs) < 6 {
		return okxZKAggregates{}, fmt.Errorf("public_inputs: need 6 elements, got %d", len(env.PublicInputs))
	}
	equity, err := parseOKXPublicU64(env.PublicInputs[0])
	if err != nil {
		return okxZKAggregates{}, fmt.Errorf("equity: %w", err)
	}
	debt, err := parseOKXPublicU64(env.PublicInputs[1])
	if err != nil {
		return okxZKAggregates{}, fmt.Errorf("debt: %w", err)
	}
	rootHex, err := okxMerkleRootHex(env.PublicInputs[2:6])
	if err != nil {
		return okxZKAggregates{}, err
	}
	liability := equity - debt
	return okxZKAggregates{
		Equity:    equity,
		Debt:      debt,
		Liability: liability,
		MerkleHex: rootHex,
	}, nil
}

func parseOKXPublicU64(n json.Number) (uint64, error) {
	return strconv.ParseUint(n.String(), 10, 64)
}

func okxMerkleRootHex(limbs []json.Number) (string, error) {
	if len(limbs) != 4 {
		return "", fmt.Errorf("merkle root: need 4 limbs, got %d", len(limbs))
	}
	var buf [32]byte
	for i, limb := range limbs {
		u, err := parseOKXPublicU64(limb)
		if err != nil {
			return "", fmt.Errorf("merkle limb %d: %w", i, err)
		}
		binary.LittleEndian.PutUint64(buf[i*8:], u)
	}
	return fmt.Sprintf("%x", buf[:]), nil
}

func normalizeMerkleHex(s string) string {
	return strings.TrimPrefix(strings.ToLower(strings.TrimSpace(s)), "0x")
}

func merkleHexEqual(a, b string) bool {
	return normalizeMerkleHex(a) == normalizeMerkleHex(b)
}
