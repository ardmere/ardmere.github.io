package walletzip

import (
	"container/heap"
	"testing"

	"github.com/shopspring/decimal"
)

func TestDepositHeapKeepsLargest(t *testing.T) {
	h := &depositHeap{}
	heap.Init(h)
	for _, bal := range []int64{1, 100, 50, 200, 75} {
		r := Row{Coin: "USDT", Balance: decimal.NewFromInt(bal)}
		if h.Len() < 3 {
			heap.Push(h, r)
			continue
		}
		if r.Balance.GreaterThan((*h)[0].Balance) {
			(*h)[0] = r
			heap.Fix(h, 0)
		}
	}
	if h.Len() != 3 {
		t.Fatalf("len=%d", h.Len())
	}
	min := (*h)[0].Balance.IntPart()
	if min < 75 {
		t.Fatalf("min balance %d too low", min)
	}
}

func TestFormatDepositSampleSummary(t *testing.T) {
	s := &DepositSample{Rows: make([]Row, 10), VerifiableRows: 1000, ValueCoverage: 0.991, UnsupportedRows: 5}
	got := FormatDepositSampleSummary(s)
	if got == "" || got == "no sample" {
		t.Fatalf("unexpected %q", got)
	}
}
