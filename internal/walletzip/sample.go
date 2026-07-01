package walletzip

import (
	"container/heap"
	"sort"

	"github.com/shopspring/decimal"
)

// SampleTopRowsByPair keeps up to topKPerPair highest-balance rows per coin|network,
// capped at maxTotal rows overall (value-biased subsample for large wallet CSVs).
func SampleTopRowsByPair(rows []Row, topKPerPair, maxTotal int) []Row {
	if topKPerPair <= 0 {
		topKPerPair = 50
	}
	if maxTotal <= 0 {
		maxTotal = 500
	}
	if len(rows) <= maxTotal {
		return rows
	}

	byPair := map[string][]Row{}
	for _, r := range rows {
		key := r.Coin + "|" + r.Network
		byPair[key] = append(byPair[key], r)
	}

	var picked []Row
	for _, group := range byPair {
		h := &rowMaxHeap{}
		heap.Init(h)
		for _, r := range group {
			if h.Len() < topKPerPair {
				heap.Push(h, r)
				continue
			}
			if r.Balance.GreaterThan(h.Peek().Balance) {
				heap.Pop(h)
				heap.Push(h, r)
			}
		}
		for h.Len() > 0 {
			picked = append(picked, heap.Pop(h).(Row))
		}
	}

	sort.Slice(picked, func(i, j int) bool {
		return picked[i].Balance.GreaterThan(picked[j].Balance)
	})
	if len(picked) > maxTotal {
		picked = picked[:maxTotal]
	}
	return picked
}

// ValueCoverageOf returns sampledBalance / totalBalance for the given row sets.
func ValueCoverageOf(all, sample []Row) float64 {
	total := decimal.Zero
	for _, r := range all {
		total = total.Add(r.Balance)
	}
	if total.IsZero() {
		return 0
	}
	sub := decimal.Zero
	for _, r := range sample {
		sub = sub.Add(r.Balance)
	}
	f, _ := sub.Div(total).Float64()
	return f
}

type rowMaxHeap []Row

func (h rowMaxHeap) Len() int           { return len(h) }
func (h rowMaxHeap) Less(i, j int) bool { return h[i].Balance.LessThan(h[j].Balance) }
func (h rowMaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *rowMaxHeap) Push(x any) { *h = append(*h, x.(Row)) }
func (h *rowMaxHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

func (h rowMaxHeap) Peek() Row { return h[0] }
