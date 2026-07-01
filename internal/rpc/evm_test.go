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
