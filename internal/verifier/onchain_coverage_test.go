package verifier

import (
	"testing"

	"github.com/ardmere/ardmere/internal/rpc"
	"github.com/ardmere/ardmere/internal/verifier/onchainconfig"
)

func TestShouldReportUnsupportedDedup(t *testing.T) {
	onchainconfig.ResetForTest()
	nativeMu.Lock()
	nativeCache = map[string]map[string]rpc.Network{}
	nativeMu.Unlock()
	ledgerMu.Lock()
	ledgerCache = map[string]map[string]onchainconfig.LedgerSpec{}
	ledgerMu.Unlock()

	ex := "binance"
	if shouldReportUnsupportedHot(ex, "ETH", "ETH") {
		t.Fatal("ETH|ETH covered by hot")
	}
	if shouldReportUnsupportedHot(ex, "USDT", "ETH") {
		t.Fatal("USDT|ETH covered by token")
	}
	if shouldReportUnsupportedHot(ex, "BTC", "BTC") {
		t.Fatal("BTC|BTC covered by ledger")
	}
	if !shouldReportUnsupportedHot(ex, "UNKNOWN", "ETH") {
		t.Fatal("UNKNOWN|ETH should be unsupported in hot")
	}
	if shouldReportUnsupportedToken(ex, "BNB", "BSC") {
		t.Fatal("BNB|BSC owned by hot, not token unsupported")
	}
	if shouldReportUnsupportedLedger(ex, "FDUSD", "BSC") {
		t.Fatal("FDUSD|BSC owned by token, not ledger unsupported")
	}
}
