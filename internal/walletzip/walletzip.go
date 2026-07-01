// Package walletzip downloads, parses and aggregates Binance's per-snapshot
// wallet address ZIP. The ZIP contains two CSVs:
//
//	PR<DDMMMYY>_HotCold.csv   ~10^3 rows  (Binance owned cold/hot wallets)
//	PR<DDMMMYY>_Deposit.csv   ~10^7 rows  (per-user deposit addresses)
//
// Schema (both files):
//
//	coin,network,address,balance,Height,Third party custodian name
//
// See docs/binance-por-data-guide.md §6.
package walletzip

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/shopspring/decimal"
)

// Row is one parsed CSV row.
type Row struct {
	Coin             string
	Network          string
	Address          string
	Balance          decimal.Decimal
	Height           int64
	CustodianName    string // empty => exchange-owned, non-empty => 3rd-party custodian
}

// File identifies one of the two CSVs inside the ZIP.
type File int

const (
	HotCold File = iota
	Deposit
	OKXAddress
)

func (f File) String() string {
	switch f {
	case HotCold:
		return "HotCold"
	case Deposit:
		return "Deposit"
	case OKXAddress:
		return "OKXAddress"
	default:
		return "unknown"
	}
}

// Download fetches the ZIP at url to outDir and returns the local path,
// the sha256 hex, and the size in bytes. The local file is content-addressed.
func Download(ctx context.Context, url, outDir string) (path, sumHex string, size int64, err error) {
	req, e := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if e != nil {
		return "", "", 0, e
	}
	req.Header.Set("User-Agent", "ardmere/0.1 (+https://ardmere.org)")

	resp, e := http.DefaultClient.Do(req)
	if e != nil {
		return "", "", 0, e
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", "", 0, fmt.Errorf("GET %s: HTTP %d", url, resp.StatusCode)
	}

	tmp, e := os.CreateTemp(outDir, "walletzip-*.tmp")
	if e != nil {
		return "", "", 0, e
	}
	tmpName := tmp.Name()
	defer func() {
		if err != nil {
			os.Remove(tmpName)
		}
	}()

	hasher := sha256.New()
	n, e := io.Copy(io.MultiWriter(tmp, hasher), resp.Body)
	tmp.Close()
	if e != nil {
		return "", "", 0, e
	}

	sumHex = hex.EncodeToString(hasher.Sum(nil))
	final := filepath.Join(outDir, sumHex+".zip")
	if e := os.Rename(tmpName, final); e != nil {
		return "", "", 0, e
	}
	return final, sumHex, n, nil
}

// OpenCSV returns a streaming csv.Reader over the requested file inside the ZIP,
// plus a closer that releases both the inner Reader and the outer ZIP.
func OpenCSV(zipPath string, which File) (*csv.Reader, io.Closer, error) {
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, nil, err
	}

	var pick *zip.File
	for _, f := range zr.File {
		name := strings.ToLower(f.Name)
		if which == HotCold && strings.Contains(name, "hotcold") {
			pick = f
			break
		}
		if which == Deposit && strings.Contains(name, "deposit") {
			pick = f
			break
		}
		if which == OKXAddress && strings.HasPrefix(name, "okx_por_") && strings.HasSuffix(strings.ToLower(name), ".csv") && !strings.Contains(strings.ToLower(name), "staking") {
			pick = f
			break
		}
	}
	if pick == nil {
		zr.Close()
		return nil, nil, fmt.Errorf("could not find %s csv inside %s", which, zipPath)
	}
	rc, err := pick.Open()
	if err != nil {
		zr.Close()
		return nil, nil, err
	}
	r := csv.NewReader(rc)
	r.FieldsPerRecord = -1 // Binance occasionally pads quotes; be lenient
	r.LazyQuotes = true

	return r, &combinedCloser{rc: rc, zr: zr}, nil
}

type combinedCloser struct {
	rc io.Closer
	zr *zip.ReadCloser
}

// WalletFileForExchange selects the wallet CSV inside a zip for an exchange.
func WalletFileForExchange(exchange string) File {
	if exchange == "okx" {
		return OKXAddress
	}
	return HotCold
}

func (c *combinedCloser) Close() error {
	_ = c.rc.Close()
	return c.zr.Close()
}

// ForEachRow streams rows from `which` CSV inside zipPath, invoking yield for
// each parsed row. Header row is consumed automatically. Stops on first non-nil
// error from yield.
func ForEachRow(zipPath string, which File, yield func(Row) error) (rowsRead int64, err error) {
	if which == OKXAddress {
		return forEachOKXRow(zipPath, yield)
	}
	r, closer, err := OpenCSV(zipPath, which)
	if err != nil {
		return 0, err
	}
	defer closer.Close()

	header, err := r.Read()
	if err != nil {
		return 0, fmt.Errorf("read header: %w", err)
	}
	if len(header) < 6 || strings.ToLower(header[0]) != "coin" {
		return 0, fmt.Errorf("unexpected header: %v", header)
	}

	for {
		rec, e := r.Read()
		if e == io.EOF {
			return rowsRead, nil
		}
		if e != nil {
			return rowsRead, fmt.Errorf("row %d: %w", rowsRead+1, e)
		}
		if len(rec) < 6 {
			continue
		}
		bal, e := decimal.NewFromString(rec[3])
		if e != nil {
			return rowsRead, fmt.Errorf("row %d balance %q: %w", rowsRead+1, rec[3], e)
		}
		var h int64
		_, e = fmt.Sscanf(rec[4], "%d", &h)
		if e != nil {
			return rowsRead, fmt.Errorf("row %d height %q: %w", rowsRead+1, rec[4], e)
		}
		row := Row{
			Coin:          rec[0],
			Network:       rec[1],
			Address:       rec[2],
			Balance:       bal,
			Height:        h,
			CustodianName: rec[5],
		}
		if err := yield(row); err != nil {
			return rowsRead, err
		}
		rowsRead++
	}
}
