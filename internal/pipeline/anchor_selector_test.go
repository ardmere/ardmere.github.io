package pipeline

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ardmere/ardmere/internal/keccak"
)

func TestAnchorSelector_matchesV6Signature(t *testing.T) {
	sig := "anchorSnapshot(string,string,bytes32,uint32,uint64,uint32,bytes32,bytes32,bytes32,uint8,uint16,uint8,bool,uint16,uint16,uint8)"
	got := keccak.Sum256([]byte(sig))
	want, _ := hex.DecodeString(strings.TrimPrefix(anchorSelector, "0x"))
	for i := 0; i < 4; i++ {
		if got[i] != want[i] {
			t.Fatalf("selector mismatch at %d: got %x want %x", i, got[:4], want)
		}
	}
}
