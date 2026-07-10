package rpc

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAlchemyBalanceAtHeight(t *testing.T) {
	const wantBal = int64(4225100)
	const height = int64(900000)
	addr := fmt.Sprintf("bc1qalchemytest%x", height)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/test-key/api/v2/address/"+addr {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		if r.URL.Query().Get("to") != "900000" {
			t.Fatalf("unexpected to=%q", r.URL.Query().Get("to"))
		}
		if r.URL.Query().Get("details") != "basic" {
			t.Fatalf("unexpected details=%q", r.URL.Query().Get("details"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"balance":"4225100"}`))
	}))
	defer srv.Close()

	oldHost := alchemyUTXOHosts
	alchemyUTXOHosts = map[string]string{"bitcoin": srv.URL}
	defer func() { alchemyUTXOHosts = oldHost }()

	c := NewLedger()
	c.alchemyKey = "test-key"
	bal, used, err := c.AlchemyBalanceAtHeight(context.Background(), "bitcoin", addr, height)
	if err != nil {
		t.Fatal(err)
	}
	if bal != wantBal {
		t.Fatalf("balance=%d want %d", bal, wantBal)
	}
	if used == "" || used == "cache" {
		t.Fatalf("expected used URL, got %q", used)
	}
}

func TestAlchemyBalanceAtHeightMissingKey(t *testing.T) {
	c := NewLedger()
	c.alchemyKey = ""
	t.Setenv("ALCHEMY_KEY", "")
	_, _, err := c.AlchemyBalanceAtHeight(context.Background(), "bitcoin", "bc1qtest", 1)
	if err == nil {
		t.Fatal("expected error without ALCHEMY_KEY")
	}
}

func TestAlchemyBalanceAtHeightUnsupportedChain(t *testing.T) {
	c := NewLedger()
	c.alchemyKey = "test-key"
	_, _, err := c.AlchemyBalanceAtHeight(context.Background(), "dogecoin", "addr", 1)
	if err == nil {
		t.Fatal("expected error for unsupported chain")
	}
}

func TestRedactAlchemyKey(t *testing.T) {
	u := "https://bitcoin-mainnet.g.alchemy.com/v2/secret123/api/v2/address/x"
	got := redactAlchemyKey(u, "secret123")
	if got != "https://bitcoin-mainnet.g.alchemy.com/v2/${ALCHEMY_KEY}/api/v2/address/x" {
		t.Fatalf("redact: %q", got)
	}
}
