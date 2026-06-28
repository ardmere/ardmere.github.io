package verifier

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEtherscanETHDepositIndexerDepositedETH(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("action") != "txlist" {
			t.Fatalf("action = %s; want txlist", r.URL.Query().Get("action"))
		}
		if r.URL.Query().Get("endblock") != "12345" {
			t.Fatalf("endblock = %s; want 12345", r.URL.Query().Get("endblock"))
		}
		fmt.Fprintf(w, `{
			"status":"1",
			"message":"OK",
			"result":[
				{"hash":"0x1","to":"%s","value":"32000000000000000000","isError":"0"},
				{"hash":"0x2","to":"0x0000000000000000000000000000000000000000","value":"99000000000000000000","isError":"0"},
				{"hash":"0x3","to":"%s","value":"32000000000000000000","isError":"1"}
			]
		}`, ethDepositContract, ethDepositContract)
	}))
	defer srv.Close()

	idx := NewEtherscanETHDepositIndexer("test-key")
	idx.BaseURL = srv.URL
	got, err := idx.DepositedETH(context.Background(), "0x32e11a20337ebc79abd0eeab2d91bafbd9591149", 12345)
	if err != nil {
		t.Fatalf("DepositedETH: %v", err)
	}
	if got.Deposited.String() != "32" {
		t.Fatalf("Deposited = %s; want 32", got.Deposited)
	}
	if got.TxCount != 1 {
		t.Fatalf("TxCount = %d; want 1", got.TxCount)
	}
}

func TestEtherscanETHDepositIndexerRequiresAPIKey(t *testing.T) {
	idx := NewEtherscanETHDepositIndexer("")
	if _, err := idx.DepositedETH(context.Background(), "0x32e11a20337ebc79abd0eeab2d91bafbd9591149", 12345); err == nil {
		t.Fatal("expected missing API key error")
	}
}
