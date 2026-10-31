package kucoinapi_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ardmere/ardmere/internal/exchanges/kucoin/kucoinapi"
)

func fixturePath(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "fixtures", "kucoin", name)
}

func TestParseSummaryBytesSSR(t *testing.T) {
	raw := []byte(`{
		"fetchedAt":"2026-10-31T01:00:00Z",
		"source":"por-page-ssr",
		"info":{
			"latestAuditDate":1790783999000,
			"reserveAsset":[
				{"currency":"BTC","userAsset":"7358.6668991","walletAsset":"7994.36640315","reserveRate":"109"},
				{"currency":"ETH","userAsset":"99131.53844795","walletAsset":"115905.46471533","reserveRate":"117"}
			]
		},
		"coinList":{
			"latestAuditDate":1790783999000,
			"reserveAsset":[
				{"currency":"BTC","userAsset":"7358.6668991","walletAsset":"7994.36640315","reserveRate":"109"},
				{"currency":"ETH","userAsset":"99131.53844795","walletAsset":"115905.46471533","reserveRate":"117"}
			]
		}
	}`)
	got, err := kucoinapi.ParseSummaryBytes(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.AuditID != "20260930" {
		t.Fatalf("audit id: %s", got.AuditID)
	}
	if got.TotalReserveRate != 109 {
		t.Fatalf("total rate (min): %v", got.TotalReserveRate)
	}
	if len(got.CoinRows) != 2 || got.CoinRows[0].Coin != "BTC" || got.CoinRows[0].ReserveRate != 109 {
		t.Fatalf("rows: %+v", got.CoinRows)
	}
}

func TestParsePoRPageHTML(t *testing.T) {
	html := []byte(`<!DOCTYPE html><html><body>
<script id="__NEXT_DATA__" type="application/json">{"props":{"pageProps":{"initialReserve":{"latestAuditDate":1790783999000,"reserveAsset":[{"currency":"USDT","userAsset":"1","walletAsset":"1.1","reserveRate":"110"}]}}}}</script>
</body></html>`)
	raw, got, err := kucoinapi.ParsePoRPageHTML(html)
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) == 0 || got.AuditID != "20260930" || len(got.CoinRows) != 1 {
		t.Fatalf("got %+v raw=%s", got, string(raw))
	}
}

func TestParseSummaryFixture(t *testing.T) {
	raw, err := os.ReadFile(fixturePath("20260930-summary.json"))
	if err != nil {
		t.Fatal(err)
	}
	got, err := kucoinapi.ParseSummaryBytes(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.AuditID != "20260930" || len(got.CoinRows) < 4 {
		t.Fatalf("unexpected: %+v", got)
	}
}
