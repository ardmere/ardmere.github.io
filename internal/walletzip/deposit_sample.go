package walletzip

import (
	"container/heap"
	"fmt"
	"sort"

	"github.com/shopspring/decimal"
)

// DepositSampleOpts controls value-weighted sampling of Deposit.csv rows.
type DepositSampleOpts struct {
	TopKPerCoin         int
	MaxTotal            int
	ValueCoverageTarget float64 // e.g. 0.99
}

// DepositSample is a value-biased subset of verifiable deposit rows.
type DepositSample struct {
	Rows            []Row
	VerifiableRows  int64
	TotalBalance    decimal.Decimal
	SampledBalance  decimal.Decimal
	ValueCoverage   float64
	UnsupportedRows int64
}

type depositHeap []Row

func (h depositHeap) Len() int           { return len(h) }
func (h depositHeap) Less(i, j int) bool { return h[i].Balance.LessThan(h[j].Balance) }
func (h depositHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *depositHeap) Push(x any) { *h = append(*h, x.(Row)) }
func (h *depositHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

// SampleDepositRows streams Deposit.csv twice and returns high-balance rows
// covering up to ValueCoverageTarget of verifiable exchange-owned balance.
func SampleDepositRows(zipPath string, routable func(Row) bool, opts DepositSampleOpts) (*DepositSample, error) {
	if opts.TopKPerCoin <= 0 {
		opts.TopKPerCoin = 5000
	}
	if opts.MaxTotal <= 0 {
		opts.MaxTotal = 50000
	}
	if opts.ValueCoverageTarget <= 0 || opts.ValueCoverageTarget > 1 {
		opts.ValueCoverageTarget = 0.99
	}

	out := &DepositSample{}
	err := forEachDepositRow(zipPath, func(r Row) error {
		if r.CustodianName != "" {
			return nil
		}
		if !routable(r) {
			out.UnsupportedRows++
			return nil
		}
		out.VerifiableRows++
		out.TotalBalance = out.TotalBalance.Add(r.Balance)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if out.VerifiableRows == 0 {
		return out, nil
	}

	heaps := map[string]*depositHeap{}
	err = forEachDepositRow(zipPath, func(r Row) error {
		if r.CustodianName != "" || !routable(r) {
			return nil
		}
		h, ok := heaps[r.Coin]
		if !ok {
			init := depositHeap{}
			h = &init
			heaps[r.Coin] = h
			heap.Init(h)
		}
		if h.Len() < opts.TopKPerCoin {
			heap.Push(h, r)
			return nil
		}
		if r.Balance.GreaterThan((*h)[0].Balance) {
			(*h)[0] = r
			heap.Fix(h, 0)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var candidates []Row
	for _, h := range heaps {
		for _, r := range *h {
			candidates = append(candidates, r)
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Balance.GreaterThan(candidates[j].Balance)
	})

	target := out.TotalBalance.Mul(decimal.NewFromFloat(opts.ValueCoverageTarget))
	for _, r := range candidates {
		if len(out.Rows) >= opts.MaxTotal {
			break
		}
		if out.SampledBalance.GreaterThanOrEqual(target) && len(out.Rows) > 0 {
			break
		}
		out.Rows = append(out.Rows, r)
		out.SampledBalance = out.SampledBalance.Add(r.Balance)
	}
	if out.TotalBalance.IsZero() {
		out.ValueCoverage = 0
	} else {
		out.ValueCoverage, _ = out.SampledBalance.Div(out.TotalBalance).Float64()
	}
	return out, nil
}

func forEachDepositRow(zipPath string, fn func(Row) error) error {
	_, err := ForEachRow(zipPath, Deposit, fn)
	return err
}

// FormatDepositSampleSummary returns a short human-readable sample summary.
func FormatDepositSampleSummary(s *DepositSample) string {
	if s == nil {
		return "no sample"
	}
	return fmt.Sprintf("%d/%d verifiable rows sampled, %.2f%% value coverage (unsupported=%d)",
		len(s.Rows), s.VerifiableRows, s.ValueCoverage*100, s.UnsupportedRows)
}
