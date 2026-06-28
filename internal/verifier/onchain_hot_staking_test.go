package verifier

import (
	"encoding/hex"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/rpc"
)

func TestDecodeSecondAddressArray(t *testing.T) {
	first := []string{
		"0x1111111111111111111111111111111111111111",
	}
	second := []string{
		"0x2222222222222222222222222222222222222222",
		"0x3333333333333333333333333333333333333333",
	}

	raw := append(rpc.To32Bytes([]byte{0x40}), rpc.To32Bytes([]byte{0x80})...)
	raw = append(raw, rpc.To32Bytes([]byte{byte(len(first))})...)
	for _, addr := range first {
		raw = append(raw, encodedAddressWord(t, addr)...)
	}
	raw = append(raw, rpc.To32Bytes([]byte{byte(len(second))})...)
	for _, addr := range second {
		raw = append(raw, encodedAddressWord(t, addr)...)
	}

	got, err := decodeSecondAddressArray(raw)
	if err != nil {
		t.Fatalf("decodeSecondAddressArray: %v", err)
	}
	if len(got) != len(second) {
		t.Fatalf("len(got) = %d; want %d", len(got), len(second))
	}
	for i := range second {
		if got[i] != second[i] {
			t.Fatalf("got[%d] = %s; want %s", i, got[i], second[i])
		}
	}
}

func TestLikelyEthDepositGap(t *testing.T) {
	claim := decimal.RequireFromString("80000.017")
	liquid := decimal.RequireFromString("0.017")
	gap := claim.Sub(liquid)
	if !likelyEthDepositGap(claim, liquid, gap) {
		t.Fatal("expected 80k ETH near-empty EOA gap to be deposit-likely")
	}
	if likelyEthDepositGap(decimal.NewFromInt(10), decimal.NewFromInt(0), decimal.NewFromInt(10)) {
		t.Fatal("small non-32 ETH gap should not be deposit-likely")
	}
	// OKX-style rows: sub-1k claim, near-empty EOA
	if !likelyEthDepositGap(decimal.RequireFromString("914.9"), decimal.RequireFromString("0.499"), decimal.RequireFromString("914.4")) {
		t.Fatal("expected 914 ETH / 0.5 liquid to be deposit-likely")
	}
	// Large omnibus label with small liquid share
	if !likelyEthDepositGap(decimal.RequireFromString("311999"), decimal.RequireFromString("1763"), decimal.RequireFromString("310236")) {
		t.Fatal("expected 312k ETH / 1.7k liquid to be deposit-likely")
	}
}

func TestEthHotInternalCustodyLikely(t *testing.T) {
	ok, note := ethHotInternalCustodyLikely(decimal.RequireFromString("832"), decimal.RequireFromString("15.12"))
	if !ok || note == "" {
		t.Fatal("expected omnibus-like ETH hot row")
	}
	if ok, _ := ethHotInternalCustodyLikely(decimal.RequireFromString("50"), decimal.RequireFromString("49")); ok {
		t.Fatal("small claim should not trigger")
	}
}

func encodedAddressWord(t *testing.T, addr string) []byte {
	t.Helper()
	b, err := hex.DecodeString(addr[2:])
	if err != nil {
		t.Fatal(err)
	}
	return rpc.To32Bytes(b)
}
