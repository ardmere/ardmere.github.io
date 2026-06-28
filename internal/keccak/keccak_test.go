package keccak

import (
	"encoding/hex"
	"testing"
)

// Cross-checks against `cast keccak ...` outputs to detect any regression.
func TestSum256_KnownVectors(t *testing.T) {
	cases := []struct {
		in   string
		want string // hex (no 0x)
	}{
		{"", "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"},
		{"abc", "4e03657aea45a94fc7d47ba826c8d667c0d1e6e33a64a036ec44f58fa12d6c45"},
		{"artifact-bundle", "36c9980257d10f782b645c08d76eeeb3914561ba1f4e112c08a716ba17970554"},
		{"PR01APR26", "60965eb19c93667e7d40a3f7641303fa44e52e28a7e77eec2663a1232f028857"},
		{"anchor(bytes32,bytes32,bytes32)", "6ed37a1d6b8e5ad9bb48fa948ccc6bdffe993ad459fbe175d30a654afef93292"},
	}
	for _, c := range cases {
		got := Sum256([]byte(c.in))
		gotHex := hex.EncodeToString(got[:])
		if gotHex != c.want {
			t.Errorf("Sum256(%q) = %s; want %s", c.in, gotHex, c.want)
		}
	}
}
