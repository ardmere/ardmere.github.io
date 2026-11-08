package pipeline

import (
	"fmt"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/walletzip"
)

// LoadPublishedWalletSet loads the exchange-published address set from a wallet zip.
// Returns publishedUnknown=true when path is empty (UNVERIFIABLE denominator).
func LoadPublishedWalletSet(exchange, walletZipPath string) (count int, balances map[string]map[string]decimal.Decimal, publishedUnknown bool, err error) {
	if walletZipPath == "" {
		return 0, nil, true, nil
	}
	var rows []walletzip.Row
	file := walletzip.WalletFileForExchange(exchange)
	_, err = walletzip.ForEachRow(walletZipPath, file, func(r walletzip.Row) error {
		rows = append(rows, r)
		return nil
	})
	if err != nil {
		return 0, nil, false, fmt.Errorf("read published wallet set: %w", err)
	}
	count, balances = CollectPublishedAddresses(rows)
	return count, balances, false, nil
}
