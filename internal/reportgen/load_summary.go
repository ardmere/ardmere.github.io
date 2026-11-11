package reportgen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ardmere/ardmere/internal/exchanges/binance"
	"github.com/ardmere/ardmere/internal/exchanges/binance/bapi"
	"github.com/ardmere/ardmere/internal/exchanges/bingx"
	bingxapi "github.com/ardmere/ardmere/internal/exchanges/bingx/bingxapi"
	"github.com/ardmere/ardmere/internal/exchanges/bitget"
	"github.com/ardmere/ardmere/internal/exchanges/bitget/bitgetapi"
	"github.com/ardmere/ardmere/internal/exchanges/bitunix"
	"github.com/ardmere/ardmere/internal/exchanges/bitunix/bitunixapi"
	"github.com/ardmere/ardmere/internal/exchanges/btcc"
	"github.com/ardmere/ardmere/internal/exchanges/btcc/btccapi"
	"github.com/ardmere/ardmere/internal/exchanges/bybit"
	"github.com/ardmere/ardmere/internal/exchanges/bybit/bybitapi"
	"github.com/ardmere/ardmere/internal/exchanges/coinex"
	"github.com/ardmere/ardmere/internal/exchanges/coinex/coinexapi"
	"github.com/ardmere/ardmere/internal/exchanges/coinw"
	"github.com/ardmere/ardmere/internal/exchanges/coinw/coinwapi"
	"github.com/ardmere/ardmere/internal/exchanges/gateio"
	"github.com/ardmere/ardmere/internal/exchanges/gateio/gateapi"
	"github.com/ardmere/ardmere/internal/exchanges/htx"
	"github.com/ardmere/ardmere/internal/exchanges/htx/htxapi"
	"github.com/ardmere/ardmere/internal/exchanges/kraken"
	"github.com/ardmere/ardmere/internal/exchanges/kraken/krakenapi"
	"github.com/ardmere/ardmere/internal/exchanges/kucoin"
	"github.com/ardmere/ardmere/internal/exchanges/kucoin/kucoinapi"
	"github.com/ardmere/ardmere/internal/exchanges/mexc"
	"github.com/ardmere/ardmere/internal/exchanges/mexc/mexcapi"
	"github.com/ardmere/ardmere/internal/exchanges/okx"
	"github.com/ardmere/ardmere/internal/exchanges/okx/okxapi"
	"github.com/ardmere/ardmere/internal/exchanges/ourbit"
	"github.com/ardmere/ardmere/internal/exchanges/ourbit/ourbitapi"
	"github.com/ardmere/ardmere/internal/exchanges/phemex"
	"github.com/ardmere/ardmere/internal/exchanges/phemex/phemexapi"
	"github.com/ardmere/ardmere/internal/exchanges/poloniex"
	"github.com/ardmere/ardmere/internal/exchanges/poloniex/poloniexapi"
	"github.com/ardmere/ardmere/internal/exchanges/toobit"
	"github.com/ardmere/ardmere/internal/exchanges/toobit/toobitapi"
	"github.com/ardmere/ardmere/internal/exchanges/woox"
	"github.com/ardmere/ardmere/internal/exchanges/woox/wooxapi"
	"github.com/ardmere/ardmere/internal/por"
)

// LoadSummaryFile normalizes a saved summary artifact into por.Snapshot.
func LoadSummaryFile(exchange, path string) (por.Snapshot, error) {
	if exchange == "gateio" {
		if snap, err := loadGateDir(filepath.Dir(path)); err == nil && strings.Contains(path, "getProofOfReservesInfo") {
			return snap, nil
		}
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return por.Snapshot{}, err
	}
	return LoadSummaryBytes(exchange, raw)
}

// LoadSummaryBytes normalizes summary JSON for a known exchange adapter.
func LoadSummaryBytes(exchange string, raw []byte) (por.Snapshot, error) {
	switch exchange {
	case "binance":
		snap, err := bapi.ParseSnapshotBytes(raw)
		if err != nil {
			return por.Snapshot{}, err
		}
		return binance.Normalize(snap, binance.Meta{Exchange: exchange}), nil
	case "gateio":
		parsed, err := gateapi.ParseSummaryBytes(raw)
		if err != nil {
			return por.Snapshot{}, err
		}
		return gateio.Normalize(parsed, gateio.Meta{Exchange: exchange}), nil
	case "htx":
		parsed, err := htxapi.ParseSummaryBytes(raw)
		if err != nil {
			return por.Snapshot{}, err
		}
		return htx.Normalize(parsed, htx.Meta{Exchange: exchange}), nil
	case "okx":
		bundle, err := okxapi.ParseSummaryBytes(raw)
		if err != nil {
			return por.Snapshot{}, err
		}
		return okx.Normalize(bundle.Audit, okx.Meta{Exchange: exchange}), nil
	case "bybit":
		parsed, err := bybitapi.ParseSummaryBytes(raw)
		if err != nil {
			return por.Snapshot{}, err
		}
		return bybit.Normalize(parsed, bybit.Meta{Exchange: exchange}), nil
	case "bitget":
		parsed, err := bitgetapi.ParseSummaryBytes(raw)
		if err != nil {
			return por.Snapshot{}, err
		}
		return bitget.Normalize(parsed, bitget.Meta{Exchange: exchange}), nil
	case "mexc":
		parsed, err := mexcapi.ParseSummaryBytes(raw)
		if err != nil {
			return por.Snapshot{}, err
		}
		return mexc.Normalize(parsed, mexc.Meta{Exchange: exchange}), nil
	case "kucoin":
		parsed, err := kucoinapi.ParseSummaryBytes(raw)
		if err != nil {
			return por.Snapshot{}, err
		}
		return kucoin.Normalize(parsed, kucoin.Meta{Exchange: exchange}), nil
	case "bingx":
		parsed, err := bingxapi.ParseSummaryBytes(raw)
		if err != nil {
			return por.Snapshot{}, err
		}
		return bingx.Normalize(parsed, bingx.Meta{Exchange: exchange}), nil
	case "poloniex":
		parsed, err := poloniexapi.ParseSummaryBytes(raw)
		if err != nil {
			return por.Snapshot{}, err
		}
		return poloniex.Normalize(parsed, poloniex.Meta{Exchange: exchange}), nil
	case "ourbit":
		parsed, err := ourbitapi.ParseSummaryBytes(raw)
		if err != nil {
			return por.Snapshot{}, err
		}
		return ourbit.Normalize(parsed, ourbit.Meta{Exchange: exchange}), nil
	case "coinex":
		parsed, err := coinexapi.ParseSummaryBytes(raw)
		if err != nil {
			return por.Snapshot{}, err
		}
		return coinex.Normalize(parsed, coinex.Meta{Exchange: exchange}), nil
	case "bitunix":
		parsed, err := bitunixapi.ParseSummaryBytes(raw)
		if err != nil {
			return por.Snapshot{}, err
		}
		return bitunix.Normalize(parsed, bitunix.Meta{Exchange: exchange}), nil
	case "woox":
		parsed, err := wooxapi.ParseSummaryBytes(raw)
		if err != nil {
			return por.Snapshot{}, err
		}
		return woox.Normalize(parsed, woox.Meta{Exchange: exchange}), nil
	case "toobit":
		parsed, err := toobitapi.ParseSummaryBytes(raw)
		if err != nil {
			return por.Snapshot{}, err
		}
		return toobit.Normalize(parsed, toobit.Meta{Exchange: exchange}), nil
	case "phemex":
		parsed, err := phemexapi.ParseSummaryBytes(raw)
		if err != nil {
			return por.Snapshot{}, err
		}
		return phemex.Normalize(parsed, phemex.Meta{Exchange: exchange}), nil
	case "btcc":
		parsed, err := btccapi.ParseSummaryBytes(raw)
		if err != nil {
			return por.Snapshot{}, err
		}
		return btcc.Normalize(parsed, btcc.Meta{Exchange: exchange}), nil
	case "kraken":
		parsed, err := krakenapi.ParseSummaryBytes(raw)
		if err != nil {
			return por.Snapshot{}, err
		}
		return kraken.Normalize(parsed, kraken.Meta{Exchange: exchange}), nil
	case "coinw":
		parsed, err := coinwapi.ParseSummaryBytes(raw)
		if err != nil {
			return por.Snapshot{}, err
		}
		return coinw.Normalize(parsed, coinw.Meta{Exchange: exchange}), nil
	default:
		return por.Snapshot{}, fmt.Errorf("no summary loader for %s", exchange)
	}
}

func loadGateDir(dir string) (por.Snapshot, error) {
	info, err := os.ReadFile(filepath.Join(dir, "getProofOfReservesInfo.json"))
	if err != nil {
		return por.Snapshot{}, err
	}
	coins, err := os.ReadFile(filepath.Join(dir, "getProofOfReservesCoinList.json"))
	if err != nil {
		return por.Snapshot{}, err
	}
	parsed, err := gateapi.ParseSummaryBundle(info, coins)
	if err != nil {
		return por.Snapshot{}, err
	}
	return gateio.Normalize(parsed, gateio.Meta{Exchange: "gateio"}), nil
}

// FindSummaryPath looks for a committed fixture for this exchange snapshot.
func FindSummaryPath(fixturesDir, exchange, snapshotID string) (string, error) {
	candidates := []string{
		filepath.Join(fixturesDir, exchange, snapshotID+"-summary.json"),
		filepath.Join(fixturesDir, exchange, snapshotID, "summary.json"),
		filepath.Join(fixturesDir, exchange, snapshotID, "getProofOfReservesInfo.json"),
	}
	for _, p := range candidates {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p, nil
		}
	}
	return "", fmt.Errorf("no fixture summary for %s %s", exchange, snapshotID)
}
