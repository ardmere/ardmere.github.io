package rpc

import (
	"encoding/hex"
	"encoding/json"
	"testing"
)

func TestDecodeHexBigInt(t *testing.T) {
	raw := json.RawMessage(`"0x2a"`)
	bi, err := decodeHexBigInt(raw)
	if err != nil {
		t.Fatal(err)
	}
	if bi.Int64() != 42 {
		t.Fatalf("got %d want 42", bi.Int64())
	}
}

func TestDecodeHexBytes(t *testing.T) {
	raw := json.RawMessage(`"0xdeadbeef"`)
	out, err := decodeHexBytes(raw)
	if err != nil {
		t.Fatal(err)
	}
	want, _ := hex.DecodeString("deadbeef")
	if string(out) != string(want) {
		t.Fatalf("got %x want %x", out, want)
	}
}

func TestNormalizeEthCallResultEmpty(t *testing.T) {
	out := normalizeEthCallResult(nil)
	if len(out) != 32 {
		t.Fatalf("got len %d want 32", len(out))
	}
	for _, b := range out {
		if b != 0 {
			t.Fatal("expected zero padding")
		}
	}
}

func TestProvidersForHistoricalPrefersArchive(t *testing.T) {
	all := map[Network][]Provider{
		NetBSC: {
			{URL: "https://full.example", Archive: false, Weight: 100},
			{URL: "https://archive.example", Archive: true, Weight: 10},
		},
	}
	got := ProvidersFor(all, NetBSC, true)
	if len(got) != 2 || got[0].URL != "https://archive.example" {
		t.Fatalf("archive-first order: %+v", got)
	}
}

func TestResolveProvidersExpandsEnv(t *testing.T) {
	t.Setenv("INFURA_KEY", "test-infura-key")
	got := resolveProviders([]Provider{
		{URL: "https://mainnet.infura.io/v3/${INFURA_KEY}", Archive: true, Weight: 1},
		{URL: "https://mainnet.infura.io/v3/${MISSING_KEY}", Archive: true, Weight: 1},
	})
	if len(got) != 1 || got[0].URL != "https://mainnet.infura.io/v3/test-infura-key" {
		t.Fatalf("resolveProviders: %+v", got)
	}
}

func TestResultCacheRoundTrip(t *testing.T) {
	cache := NewResultCache(t.TempDir())
	raw := json.RawMessage(`"0x01"`)
	cache.Put(NetBSC, "eth_call", "0xabc", "0x01", 123, raw)
	got, ok := cache.Get(NetBSC, "eth_call", "0xabc", "0x01", 123)
	if !ok {
		t.Fatal("cache miss")
	}
	if string(got) != string(raw) {
		t.Fatalf("got %s want %s", got, raw)
	}
}
