package rpc

import (
	"context"
	"testing"
)

func TestTronBase58ToHex(t *testing.T) {
	got, err := TronBase58ToHex("TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t")
	if err != nil {
		t.Fatal(err)
	}
	want := "41a614f803b6fd780986a42c78ec9c7f77e6ded13c"
	if got != want {
		t.Fatalf("contract hex=%s want %s", got, want)
	}
}

func TestTRC20BalanceOfLive(t *testing.T) {
	ctx := context.Background()
	c := New()
	// PR01JUN26 sample holder; skip if TronGrid unavailable.
	bal, used, err := c.TRC20BalanceOf(ctx, "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t", "TJ5usJLLwjwn7Pw3TPbdzreG7dvgKzfQ5y", 83201055)
	if err != nil {
		t.Skipf("tron live probe: %v", err)
	}
	if bal.Sign() <= 0 {
		t.Fatalf("expected positive balance, got %s via %s", bal, used)
	}
}
