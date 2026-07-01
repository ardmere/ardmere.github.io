package htx

import (
	"testing"

	"github.com/ardmere/ardmere/internal/exchanges/htx/htxapi"
)

func TestNormalizeZkDerived(t *testing.T) {
	raw := htxapi.ParsedSummary{
		AuditID:      "20230910",
		AuditTimeRaw: "2023-09-10T03:42:30Z",
		MerkleRoot:   "abc",
		TotalReserveRate: 105.5,
		CoinRows: []htxapi.CoinRow{
			{Coin: "BTC", ReserveRate: 110},
		},
	}
	got := Normalize(raw, Meta{Exchange: "htx", ZkDerived: true})
	if got.Exchange != "htx" || got.ID != "20230910" {
		t.Fatalf("snapshot: %+v", got)
	}
	if got.Extra["summarySource"] != "zk-bundle-derived" {
		t.Fatalf("extra: %+v", got.Extra)
	}
}
