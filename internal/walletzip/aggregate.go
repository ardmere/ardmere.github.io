package walletzip

import (
	"github.com/shopspring/decimal"
)

// SourceStats is per-CSV-file aggregates inside a wallet zip.
type SourceStats struct {
	TotalRows    int64
	Exchange     map[string]decimal.Decimal
	ThirdParty   map[string]decimal.Decimal
	RowsByCoin   map[string]int
	HeightByCoin map[string]int64
}

func newSourceStats() *SourceStats {
	return &SourceStats{
		Exchange:     map[string]decimal.Decimal{},
		ThirdParty:   map[string]decimal.Decimal{},
		RowsByCoin:   map[string]int{},
		HeightByCoin: map[string]int64{},
	}
}

// Aggregate sums balances per coin, splitting between exchange-owned and
// 3rd-party custodian rows, across the given files of zipPath.
type Aggregate struct {
	Exchange     map[string]decimal.Decimal // coin -> sum
	ThirdParty   map[string]decimal.Decimal // coin -> sum
	RowsByCoin   map[string]int             // coin -> #rows (debug)
	HeightByCoin map[string]int64           // coin -> max(Height) seen
	TotalRows    int64
	BySource     map[File]*SourceStats
}

func newAggregate() *Aggregate {
	return &Aggregate{
		Exchange:     map[string]decimal.Decimal{},
		ThirdParty:   map[string]decimal.Decimal{},
		RowsByCoin:   map[string]int{},
		HeightByCoin: map[string]int64{},
		BySource:     map[File]*SourceStats{},
	}
}

func (agg *Aggregate) addRow(f File, r Row) {
	if r.CustodianName == "" {
		agg.Exchange[r.Coin] = agg.Exchange[r.Coin].Add(r.Balance)
	} else {
		agg.ThirdParty[r.Coin] = agg.ThirdParty[r.Coin].Add(r.Balance)
	}
	agg.RowsByCoin[r.Coin]++
	if r.Height > agg.HeightByCoin[r.Coin] {
		agg.HeightByCoin[r.Coin] = r.Height
	}
	agg.TotalRows++

	src := agg.BySource[f]
	if src == nil {
		src = newSourceStats()
		agg.BySource[f] = src
	}
	if r.CustodianName == "" {
		src.Exchange[r.Coin] = src.Exchange[r.Coin].Add(r.Balance)
	} else {
		src.ThirdParty[r.Coin] = src.ThirdParty[r.Coin].Add(r.Balance)
	}
	src.RowsByCoin[r.Coin]++
	if r.Height > src.HeightByCoin[r.Coin] {
		src.HeightByCoin[r.Coin] = r.Height
	}
	src.TotalRows++
}

// Aggregate sums all rows across the requested files. Streaming, O(coins) memory.
func AggregateFiles(zipPath string, files ...File) (*Aggregate, error) {
	agg := newAggregate()
	for _, f := range files {
		_, err := ForEachRow(zipPath, f, func(r Row) error {
			agg.addRow(f, r)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return agg, nil
}
