// Package okxwallet parses OKX Proof-of-Reserves wallet ZIP/CSV artifacts.
//
// okx_por_*.csv layout:
//
//	coin,amount                         — exchange reserve totals (top section)
//	(blank line)
//	coin,Network,Snapshot Height,address,amount,message,signature1,...
package okxwallet

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/shopspring/decimal"
)

// AddressRow is one signed wallet address line from the OKX reserves CSV.
type AddressRow struct {
	Coin       string
	Network    string
	Height     int64
	Address    string
	Amount     decimal.Decimal
	Message    string
	Signature1 string
	Signature2 string
	Script     string
	EOA1       string
	EOA2       string
}

// Aggregate mirrors walletzip.Aggregate for OKX wallet zips.
type Aggregate struct {
	Exchange     map[string]decimal.Decimal
	ThirdParty   map[string]decimal.Decimal
	RowsByCoin   map[string]int
	HeightByCoin map[string]int64
	TotalRows    int64
	StakingETH   decimal.Decimal
}

func newAggregate() *Aggregate {
	return &Aggregate{
		Exchange:     map[string]decimal.Decimal{},
		ThirdParty:   map[string]decimal.Decimal{},
		RowsByCoin:   map[string]int{},
		HeightByCoin: map[string]int64{},
	}
}

// MainCSVName returns the primary okx_por_*.csv inside zipPath.
func MainCSVName(zipPath string) (string, error) {
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer zr.Close()
	for _, f := range zr.File {
		base := strings.ToLower(filepath.Base(f.Name))
		if strings.HasPrefix(base, "okx_por_") && strings.HasSuffix(base, ".csv") && !strings.Contains(base, "staking") {
			return f.Name, nil
		}
	}
	return "", fmt.Errorf("okx main por csv not found in %s", zipPath)
}

// ForEachAddressRow streams signed address rows from zipPath.
func ForEachAddressRow(zipPath string, yield func(AddressRow) error) (int64, error) {
	csvName, err := MainCSVName(zipPath)
	if err != nil {
		return 0, err
	}
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return 0, err
	}
	defer zr.Close()
	var pick *zip.File
	for _, f := range zr.File {
		if f.Name == csvName {
			pick = f
			break
		}
	}
	if pick == nil {
		return 0, fmt.Errorf("csv %q missing", csvName)
	}
	rc, err := pick.Open()
	if err != nil {
		return 0, err
	}
	defer rc.Close()

	r := csv.NewReader(rc)
	r.FieldsPerRecord = -1
	r.LazyQuotes = true

	var rows int64
	inAddr := false
	for {
		rec, err := r.Read()
		if err == io.EOF {
			return rows, nil
		}
		if err != nil {
			return rows, err
		}
		if len(rec) == 0 {
			continue
		}
		if !inAddr {
			if len(rec) >= 2 && strings.EqualFold(rec[0], "coin") && strings.Contains(strings.ToLower(rec[1]), "network") {
				inAddr = true
			}
			continue
		}
		row, err := parseAddressRecord(rec)
		if err != nil {
			return rows, fmt.Errorf("row %d: %w", rows+1, err)
		}
		if err := yield(row); err != nil {
			return rows, err
		}
		rows++
	}
}

// AggregateZip sums address rows and ETH staking from an OKX wallet zip.
func AggregateZip(zipPath string) (*Aggregate, error) {
	agg := newAggregate()
	_, err := ForEachAddressRow(zipPath, func(r AddressRow) error {
		agg.Exchange[r.Coin] = agg.Exchange[r.Coin].Add(r.Amount)
		agg.RowsByCoin[r.Coin]++
		if r.Height > agg.HeightByCoin[r.Coin] {
			agg.HeightByCoin[r.Coin] = r.Height
		}
		agg.TotalRows++
		return nil
	})
	if err != nil {
		return nil, err
	}
	staking, err := readStakingETH(zipPath)
	if err != nil {
		return nil, err
	}
	agg.StakingETH = staking
	if staking.IsPositive() {
		agg.Exchange["ETH"] = agg.Exchange["ETH"].Add(staking)
	}
	return agg, nil
}

func parseAddressRecord(rec []string) (AddressRow, error) {
	if len(rec) < 7 {
		return AddressRow{}, fmt.Errorf("short record: %d fields", len(rec))
	}
	amt, err := decimal.NewFromString(rec[4])
	if err != nil {
		return AddressRow{}, fmt.Errorf("amount %q: %w", rec[4], err)
	}
	var h int64
	if _, err := fmt.Sscan(rec[2], &h); err != nil {
		return AddressRow{}, fmt.Errorf("height %q: %w", rec[2], err)
	}
	row := AddressRow{
		Coin:       strings.ToUpper(rec[0]),
		Network:    rec[1],
		Height:     h,
		Address:    rec[3],
		Amount:     amt,
		Message:    rec[5],
		Signature1: rec[6],
	}
	if len(rec) > 7 {
		row.Signature2 = rec[7]
	}
	if len(rec) > 8 {
		row.Script = rec[8]
	}
	if len(rec) > 9 {
		row.EOA1 = rec[9]
	}
	if len(rec) > 10 {
		row.EOA2 = rec[10]
	}
	return row, nil
}

func readStakingETH(zipPath string) (decimal.Decimal, error) {
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return decimal.Zero, err
	}
	defer zr.Close()
	var pick *zip.File
	for _, f := range zr.File {
		if strings.Contains(strings.ToLower(f.Name), "eth_staking") {
			pick = f
			break
		}
	}
	if pick == nil {
		return decimal.Zero, nil
	}
	rc, err := pick.Open()
	if err != nil {
		return decimal.Zero, err
	}
	defer rc.Close()
	r := csv.NewReader(rc)
	r.FieldsPerRecord = -1
	for {
		rec, err := r.Read()
		if err == io.EOF {
			return decimal.Zero, nil
		}
		if err != nil {
			return decimal.Zero, err
		}
		if len(rec) >= 2 && strings.HasPrefix(rec[0], "ETH") {
			return decimal.NewFromString(rec[1])
		}
	}
}
