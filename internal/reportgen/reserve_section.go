package reportgen

import (
	"fmt"
	"strings"

	"github.com/ardmere/ardmere/internal/pipeline"
)

const reserveSectionHeading = "## Self-reported liabilities, reserves, and major-coin ratios"

// writeReserveSection renders the dedicated totals / major-coin block.
// Figures are exchange self-reports priced with CoinGecko. They are not a solvency opinion.
func writeReserveSection(b *strings.Builder, fig pipeline.ReserveFigures) {
	b.WriteString(reserveSectionHeading)
	b.WriteString("\n\n")
	b.WriteString("Label: **self-reported**. These rows repeat the exchange snapshot. ")
	b.WriteString("USD totals sum per-coin customer liability and exchange reserve (self-custody + third-party custody) using CoinGecko prices at `snapshotTime`, the same price path as `reserveCoverageBps`. ")
	b.WriteString("A coin with no price is left out of the USD total and listed as unpriced. ")
	b.WriteString("`walletCoverageBps` and `reserveCoverageBps` are the separate **checked** coverage ratios. ")
	b.WriteString("This section is not solvency advice.\n\n")

	b.WriteString("| Figure | Value | Label |\n| --- | --- | --- |\n")
	fmt.Fprintf(b, "| Total liabilities (USD) | %s | self-reported |\n", usdCell(fig, fig.TotalLiabilityUsd))
	fmt.Fprintf(b, "| Total reserves (USD) | %s | self-reported |\n", usdCell(fig, fig.TotalReserveUsd))
	fmt.Fprintf(b, "| Unpriced coins | %s | — |\n\n", unpricedCell(fig))

	if fig.SourceMissing {
		b.WriteString("The summary artifact was not available when this report was generated, so USD totals were not computed.\n\n")
	} else if !fig.AmountsPublished {
		b.WriteString("The snapshot did not publish amounts, so USD totals are null. Major-coin rows below use the exchange-reported reserve rate when that coin is in the snapshot.\n\n")
	} else if fig.TotalLiabilityUsd == nil || fig.TotalReserveUsd == nil {
		b.WriteString("Amount-bearing coins were all unpriced, so USD totals are null.\n\n")
	}

	b.WriteString("| Coin | Reserve ratio | Basis | Label |\n| --- | --- | --- | --- |\n")
	if len(fig.MajorCoinReserveRates) == 0 {
		b.WriteString("| — | — | — | — |\n")
	}
	for _, row := range fig.MajorCoinReserveRates {
		label := "self-reported"
		fmt.Fprintf(b, "| `%s` | `%s` | `%s` | %s |\n", row.Coin, row.ReserveRate.String(), row.Basis, label)
	}
	b.WriteString("\n")
	b.WriteString("Reserve ratio is exchange reserve ÷ customer liability when both amounts exist (`basis=amounts`); otherwise it is the exchange-reported rate already parsed into the coin summary (`basis=reported`, percentage ÷ 100 so 1 = 100%). ")
	b.WriteString("Major coins are BTC, ETH, USDT, and USDC. A coin is omitted when the snapshot has no row for it. ")
	b.WriteString("On-chain schema v7 stores the same USD totals as uint128 whole dollars (`0` = unavailable) and the four ratios as uint16 basis points (`10000` = 100%, `65535` = unavailable).\n\n")
}

func usdCell(fig pipeline.ReserveFigures, v *float64) string {
	if v == nil {
		if fig.SourceMissing {
			return "`null` (summary artifact not loaded)"
		}
		if !fig.AmountsPublished {
			return "`null` (snapshot did not publish amounts)"
		}
		return "`null`"
	}
	return fmt.Sprintf("`%.2f`", *v)
}

func unpricedCell(fig pipeline.ReserveFigures) string {
	if len(fig.UnpricedCoins) == 0 {
		return "—"
	}
	parts := make([]string, len(fig.UnpricedCoins))
	for i, c := range fig.UnpricedCoins {
		parts[i] = "`" + strings.ReplaceAll(c, "`", "") + "`"
	}
	return strings.Join(parts, ", ")
}
