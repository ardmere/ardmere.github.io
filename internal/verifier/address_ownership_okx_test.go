package verifier

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ardmere/ardmere/internal/exchanges/okx/okxwallet"
	"github.com/shopspring/decimal"
)

func TestOKXFieldSet(t *testing.T) {
	cases := map[string]bool{
		"":       false,
		"  ":     false,
		"null":   false,
		"NULL":   false,
		`\n`:     false,
		`\N`:     false,
		"0xabc":  true,
		"sighex": true,
	}
	for in, want := range cases {
		if got := okxFieldSet(in); got != want {
			t.Fatalf("okxFieldSet(%q)=%v want %v", in, got, want)
		}
	}
}

func TestNormalizeAptosHexAddr(t *testing.T) {
	got := normalizeAptosHexAddr("0xabc")
	if !strings.HasPrefix(got, "0x") || len(got) != 66 {
		t.Fatalf("got %q", got)
	}
	if normalizeAptosHexAddr("0x"+strings.Repeat("0", 63)+"1") != normalizeAptosHexAddr("0x1") {
		t.Fatal("padding mismatch")
	}
}

func TestVerifyOKXAddressRowDualOwnerWithoutSig2(t *testing.T) {
	row := okxwallet.AddressRow{
		Coin:       "ETH",
		Network:    "ETH",
		Address:    "0x80299387783b6b8c0645a17ea72fdc3cd32e0c9367445fdbe7344742af521fc39afc1fbc3917efc3a4bd72c48814e9b0",
		Amount:     decimal.RequireFromString("32"),
		Message:    "I am an OKX address",
		Signature1: "0x57ca924e1c6a87ff5eb86c5e8a9483cd05f5865e7c54fcc4478e5e481c31b6f3678c9139b04c7f9ce379f624373956c10b676006197e7e64a10a5f0d8505d9701b",
		Signature2: "",
		EOA1:       "0xaa0cbae2dd290f8aed1b672ebe2e197fd969628b",
		EOA2:       "0xe6a64adc9270ffd033d9f06686c4617afef7cfa3",
	}
	_, err := verifyOKXAddressRow(&okxOwnershipState{authCache: map[string]string{}}, row)
	if err != nil && strings.Contains(err.Error(), "missing signature2") {
		t.Fatalf("unexpected dual-sig hard fail: %v", err)
	}
	if err != nil {
		t.Fatalf("eigenlayer single-owner path should pass for published sample: %v", err)
	}
}

func TestAptosResourceAccountAuthKeyFallback(t *testing.T) {
	resource := "0x117d2d4e8bffe853addefb9aa380eb9c4fdfbcb73ea8cd81b74138c12fcd2f82"
	recoverEOA := "0x4c1ef44079cb31851349fba50d385f708a10ec7ac612859fdbf28888d1f7b572"
	zip := filepath.Join("..", "..", "artifacts", "okx", "502260945", "raw",
		"5230bd14fe64f688c5fd480006f714e746a12b6065ac9d683ca02b89852e9e07.zip")
	if _, err := os.Stat(zip); err != nil {
		t.Skip(err)
	}

	var row *okxwallet.AddressRow
	_, err := okxwallet.ForEachAddressRow(zip, func(r okxwallet.AddressRow) error {
		if strings.EqualFold(r.Address, resource) && strings.EqualFold(r.Coin, "APTOS") {
			cp := r
			row = &cp
			return errStopIteration
		}
		return nil
	})
	if err != nil && err != errStopIteration {
		t.Fatal(err)
	}
	if row == nil {
		t.Fatal("sample Aptos resource row not found")
	}

	prev := fetchAptosAuthKey
	t.Cleanup(func() { fetchAptosAuthKey = prev })
	fetchAptosAuthKey = func(base, addr string, _ *http.Client) (string, error) {
		if normalizeAptosHexAddr(addr) != normalizeAptosHexAddr(resource) {
			return "", fmt.Errorf("unexpected addr %s", addr)
		}
		return recoverEOA, nil
	}

	path, err := verifyOKXAddressRow(&okxOwnershipState{authCache: map[string]string{}}, *row)
	if err != nil {
		t.Fatalf("expected Aptos auth_key fallback PASS, got %v", err)
	}
	if path != ownershipViaAptosAuthKey {
		t.Fatalf("path=%v want aptos auth_key", path)
	}
}

var errStopIteration = fmt.Errorf("stop")
