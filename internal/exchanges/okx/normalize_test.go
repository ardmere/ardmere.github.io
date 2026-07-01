package okx

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ardmere/ardmere/internal/exchanges/okx/okxapi"
	"github.com/ardmere/ardmere/internal/verifier"
)

func TestNormalize(t *testing.T) {
	raw := okxapi.AuditRootInfo{
		AuditID:    "506872725",
		MerkleHash: "abc123",
		CreateTime: 1778083200000,
		CapitalRatio: map[string]string{
			"BTC": "106%",
		},
		LiabilityBalances:       map[string]string{"BTC": "111188"},
		ExchangeReserveBalances: map[string]string{"BTC": "112393"},
		CustodyReserveBalances:  map[string]string{"BTC": "5541"},
		ReserveBalances:         map[string]string{"BTC": "117934"},
		ReserveCurrencies:       []string{"BTC"},
	}
	snap := Normalize(raw, Meta{Exchange: ID})
	if snap.ID != "506872725" || snap.Exchange != ID {
		t.Fatalf("snapshot: %+v", snap)
	}
	if len(snap.CoinSummaries) != 1 || snap.CoinSummaries[0].Extra["capitalRatio"] != "106%" {
		t.Fatalf("coins: %+v", snap.CoinSummaries)
	}
}

func TestSolvencyClaimFromSummary(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "..", "fixtures", "okx", "506872725", "summary.json"))
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := okxapi.ParseSummaryBytes(raw)
	if err != nil {
		t.Fatal(err)
	}
	snap := Normalize(parsed.Audit, Meta{Exchange: ID})
	v := verifier.SolvencyClaim{Snapshot: snap, SummarySha256: "abc"}.Run()
	if v.Verdict != verifier.VerdictPass {
		t.Fatalf("verdict=%s reason=%s", v.Verdict, v.Reason)
	}
}
