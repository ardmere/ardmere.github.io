package verifier

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/rpc"
	"github.com/ardmere/ardmere/internal/walletzip"
)

// tryTokenHeightOffsets re-queries balance at H-1 and H+1 when the claim does not
// match at the CSV height (snapshot block boundary ambiguity).
func (v OnchainBalanceToken) tryTokenHeightOffsets(
	ctx context.Context,
	r walletzip.Row,
	decCache *decimalsCache,
	sonicCache *sonicSFCCache,
) (decimal.Decimal, string, int64, string, bool) {
	for _, delta := range []int64{-1, 1} {
		h := r.Height + delta
		if h < 0 {
			continue
		}
		row := r
		row.Height = h
		actual, used, _, err := v.tokenBalanceAtHeight(ctx, row, decCache, sonicCache)
		if err != nil {
			continue
		}
		if balanceWithinTolerance(actual, r.Balance) {
			note := fmt.Sprintf("height-boundary: csv block %d matches at block %d", r.Height, h)
			return actual, used, h, note, true
		}
	}
	return decimal.Zero, "", 0, "", false
}

func nativeBalanceDecimal(ctx context.Context, c *rpc.Client, net rpc.Network, address string, height int64, decimals int) (decimal.Decimal, string, error) {
	wei, used, err := c.GetBalance(ctx, net, address, height)
	if err != nil {
		return decimal.Zero, used, err
	}
	return decimal.NewFromBigInt(wei, -int32(decimals)), used, nil
}

func sonicBalanceNote(claim, actual decimal.Decimal) string {
	if claim.IsZero() || actual.IsZero() {
		return ""
	}
	if actual.Div(claim).GreaterThan(decimal.NewFromFloat(0.01)) {
		return ""
	}
	return "liquid+SFC stake still << CSV; residual may be rewards/unbonding or internal ledger"
}

func polEthBalanceNote(claim, actual decimal.Decimal) string {
	if claim.IsZero() || actual.IsZero() {
		return ""
	}
	if actual.Div(claim).GreaterThan(decimal.NewFromFloat(0.01)) {
		return ""
	}
	return "POL+MATIC ERC20 on ETH still << CSV; residual may be Polygon-native POL or internal ledger"
}

func tronSnapshotMismatchNote(coin, network string) string {
	if network != "TRX" && network != "TRON" {
		return ""
	}
	if coin == "TRX" {
		return "live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes"
	}
	return "live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective)"
}

func fevmSnapshotNote(coin, network string) string {
	if network != "FEVM" {
		return ""
	}
	return "live FEVM balance vs snapshot CSV; public Filecoin EVM archive lookback limited"
}

// tronTrc20SnapshotNote is deprecated; use tronSnapshotMismatchNote.
func tronTrc20SnapshotNote(coin, network string) string {
	return tronSnapshotMismatchNote(coin, network)
}

// btcEthNativeCustodyUnverifiable marks BTC|ETH rows with zero WBTC but a
// material CSV claim — likely native/off-chain BTC custody, not WBTC on Ethereum.
func btcEthNativeCustodyUnverifiable(r walletzip.Row, actual decimal.Decimal) (bool, string) {
	if r.Coin != "BTC" || r.Network != "ETH" {
		return false, ""
	}
	if !actual.IsZero() {
		return false, ""
	}
	if r.Balance.LessThanOrEqual(decimal.NewFromFloat(0.01)) {
		return false, ""
	}
	return true, "no WBTC at snapshot height; BTC likely native/off-chain custody — EVM WBTC balanceOf cannot verify this row"
}

// stablecoinEthOmnibusMismatch marks mega-stablecoin ETH rows where on-chain ERC20
// balance is far below CSV — exchange omnibus allocation labels, not address custody.
func stablecoinEthOmnibusMismatch(r walletzip.Row, actual decimal.Decimal) (bool, string) {
	if r.Network != "ETH" {
		return false, ""
	}
	switch r.Coin {
	case "USDT", "USDC", "DAI", "BUSD", "TUSD", "FDUSD", "USDP", "PYUSD":
	default:
		return false, ""
	}
	if actual.GreaterThan(r.Balance) {
		return false, ""
	}
	if r.Balance.LessThan(decimal.NewFromInt(1_000_000)) {
		return false, ""
	}
	if actual.IsZero() || actual.Div(r.Balance).LessThan(decimal.NewFromFloat(0.001)) {
		return true, "on-chain ERC20 balance << CSV mega-stablecoin allocation; likely omnibus/internal ledger label — not single-address custody"
	}
	return false, ""
}

// classifyBalanceMismatch decides PASS / WARN (chain surplus) / FAIL (chain short).
func classifyBalanceMismatch(actual, claim decimal.Decimal) (within bool, surplus bool) {
	if balanceWithinTolerance(actual, claim) {
		return true, false
	}
	return false, actual.GreaterThan(claim)
}
