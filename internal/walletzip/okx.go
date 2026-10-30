package walletzip

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/shopspring/decimal"
)

func forEachOKXRow(zipPath string, yield func(Row) error) (int64, error) {
	csvName, err := okxMainCSVName(zipPath)
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
			if isOKXAddressHeader(rec) {
				inAddr = true
			}
			continue
		}
		row, err := parseOKXAddressRecord(rec)
		if err != nil {
			return rows, fmt.Errorf("row %d: %w", rows+1, err)
		}
		if err := yield(row); err != nil {
			return rows, err
		}
		rows++
	}
}

// AggregateOKXZip builds an aggregate for OKX internal consistency.
// Exchange totals come from the CSV top section (coin,amount), which matches
// exchangeReserveBalances in the summary; address rows are counted for coverage.
func AggregateOKXZip(zipPath string) (*Aggregate, error) {
	agg := newAggregate()
	totals, err := readOKXTopTotals(zipPath)
	if err != nil {
		return nil, err
	}
	agg.Exchange = totals
	_, err = forEachOKXRow(zipPath, func(r Row) error {
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
	return agg, nil
}

func readOKXTopTotals(zipPath string) (map[string]decimal.Decimal, error) {
	csvName, err := okxMainCSVName(zipPath)
	if err != nil {
		return nil, err
	}
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("csv %q missing", csvName)
	}
	rc, err := pick.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	r := csv.NewReader(rc)
	r.FieldsPerRecord = -1
	out := map[string]decimal.Decimal{}
	for {
		rec, err := r.Read()
		if err == io.EOF {
			return out, nil
		}
		if err != nil {
			return nil, err
		}
		if len(rec) < 2 || rec[0] == "" {
			return out, nil
		}
		if strings.EqualFold(rec[0], "coin") {
			if isOKXAddressHeader(rec) {
				return out, nil
			}
			continue
		}
		if len(rec) != 2 {
			return out, nil
		}
		amt, err := decimal.NewFromString(rec[1])
		if err != nil {
			return nil, fmt.Errorf("total %s: %w", rec[0], err)
		}
		out[strings.ToUpper(rec[0])] = amt
	}
}

func okxMainCSVName(zipPath string) (string, error) {
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

// isOKXAddressHeader detects V3 (`coin,network,...`) and V4
// (`coin,Type,Network,Snapshot Height,address,amount,...`) section headers.
func isOKXAddressHeader(rec []string) bool {
	if len(rec) < 2 || !strings.EqualFold(strings.TrimSpace(rec[0]), "coin") {
		return false
	}
	for _, col := range rec[1:] {
		if strings.Contains(strings.ToLower(col), "network") {
			return true
		}
	}
	return false
}

func parseOKXAddressRecord(rec []string) (Row, error) {
	// V4: coin, Type, Network, Snapshot Height, address, amount, message, sig...
	if len(rec) >= 8 && looksLikeOKXV4Address(rec) {
		amt, err := decimal.NewFromString(rec[5])
		if err != nil {
			return Row{}, fmt.Errorf("amount %q: %w", rec[5], err)
		}
		var h int64
		if _, err := fmt.Sscan(rec[3], &h); err != nil {
			return Row{}, fmt.Errorf("height %q: %w", rec[3], err)
		}
		return Row{
			Coin:    strings.ToUpper(rec[0]),
			Network: normalizeOKXNetwork(rec[2]),
			Address: rec[4],
			Balance: amt,
			Height:  h,
		}, nil
	}
	// V3 / earlier: coin, network, height, address, amount, message, sig...
	if len(rec) < 7 {
		return Row{}, fmt.Errorf("short record: %d fields", len(rec))
	}
	amt, err := decimal.NewFromString(rec[4])
	if err != nil {
		return Row{}, fmt.Errorf("amount %q: %w", rec[4], err)
	}
	var h int64
	if _, err := fmt.Sscan(rec[2], &h); err != nil {
		return Row{}, fmt.Errorf("height %q: %w", rec[2], err)
	}
	return Row{
		Coin:    strings.ToUpper(rec[0]),
		Network: normalizeOKXNetwork(rec[1]),
		Address: rec[3],
		Balance: amt,
		Height:  h,
	}, nil
}

func looksLikeOKXV4Address(rec []string) bool {
	// Column 1 is staking type (e.g. "Non Staking" / "Staking"), not a network code.
	typ := strings.ToLower(strings.TrimSpace(rec[1]))
	return strings.Contains(typ, "staking")
}

func normalizeOKXNetwork(net string) string {
	switch strings.ToUpper(strings.TrimSpace(net)) {
	case "TRON":
		return "TRX"
	case "APTOS":
		return "APT"
	case "POLYGON":
		return "MATIC"
	case "ETH_LINEA":
		return "LINEA"
	case "ASSET_HUB":
		return "DOT"
	default:
		return strings.ToUpper(strings.TrimSpace(net))
	}
}

func readOKXStakingETH(zipPath string) (decimal.Decimal, error) {
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
