package pipeline

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ardmere/ardmere/internal/keccak"
)

func TestAnchorSelector_v7IsFourBytes(t *testing.T) {
	ValidateAnchorSelector()
	got := keccak.Sum256([]byte(anchorSelectorSig))
	want := "0x" + hex.EncodeToString(got[:4])
	if anchorSelector != want {
		t.Fatalf("selector %s want %s", anchorSelector, want)
	}
	if strings.Contains(anchorSelectorSig, "uint128,uint128,uint16,uint16,uint16,uint16") == false {
		t.Fatal("v7 signature missing reserve total args")
	}
}
