package rpc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBTCBlockTimeAtHeight(t *testing.T) {
	t.Setenv("RPC_CACHE_DIR", t.TempDir())
	const (
		height = int64(951913)
		hash   = "00000000000000000000d40f10282220ab5904b995daf315010f1c47f8eb32a2"
	)
	want := time.Date(2026, 5, 31, 23, 37, 57, 0, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/block-height/951913":
			_, _ = w.Write([]byte(hash))
		case "/block/" + hash:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":        hash,
				"height":    height,
				"timestamp": want.Unix(),
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	old := EsploraBases
	EsploraBases = []string{srv.URL}
	defer func() { EsploraBases = old }()

	got, used, err := NewLedger().BTCBlockTimeAtHeight(context.Background(), height)
	if err != nil {
		t.Fatal(err)
	}
	if used != srv.URL {
		t.Fatalf("used=%q", used)
	}
	if !got.Equal(want) {
		t.Fatalf("got=%s want=%s", got, want)
	}

	got2, used2, err := NewLedger().BTCBlockTimeAtHeight(context.Background(), height)
	if err != nil {
		t.Fatal(err)
	}
	if used2 != "cache" {
		t.Fatalf("cache miss: used=%q", used2)
	}
	if !got2.Equal(want) {
		t.Fatalf("cached=%s want=%s", got2, want)
	}
}
