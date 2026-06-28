package okxwallet

import (
	"path/filepath"
	"testing"
)

func TestForEachAddressRowFixture(t *testing.T) {
	zipPath := filepath.Join("..", "..", "..", "..", "fixtures", "okx", "506872725", "wallet.zip")
	var rows int64
	n, err := ForEachAddressRow(zipPath, func(r AddressRow) error {
		rows++
		if r.Message != "I am an OKX address" {
			t.Fatalf("message: %q", r.Message)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 || rows != 1 {
		t.Fatalf("rows=%d n=%d", rows, n)
	}
}
