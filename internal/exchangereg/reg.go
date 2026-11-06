package exchangereg

import (
	"fmt"
	"strings"

	"github.com/ardmere/ardmere/internal/exchange"
	"github.com/ardmere/ardmere/internal/exchanges/binance"
	"github.com/ardmere/ardmere/internal/exchanges/bingx"
	"github.com/ardmere/ardmere/internal/exchanges/bitget"
	"github.com/ardmere/ardmere/internal/exchanges/bitunix"
	"github.com/ardmere/ardmere/internal/exchanges/btcc"
	"github.com/ardmere/ardmere/internal/exchanges/bybit"
	"github.com/ardmere/ardmere/internal/exchanges/coinex"
	"github.com/ardmere/ardmere/internal/exchanges/gateio"
	"github.com/ardmere/ardmere/internal/exchanges/htx"
	"github.com/ardmere/ardmere/internal/exchanges/kraken"
	"github.com/ardmere/ardmere/internal/exchanges/kucoin"
	"github.com/ardmere/ardmere/internal/exchanges/mexc"
	"github.com/ardmere/ardmere/internal/exchanges/okx"
	"github.com/ardmere/ardmere/internal/exchanges/ourbit"
	"github.com/ardmere/ardmere/internal/exchanges/phemex"
	"github.com/ardmere/ardmere/internal/exchanges/poloniex"
	"github.com/ardmere/ardmere/internal/exchanges/toobit"
	"github.com/ardmere/ardmere/internal/exchanges/woox"
)

var adapters = map[string]exchange.Adapter{
	binance.ID:  binance.New(),
	bingx.ID:    bingx.New(),
	bitget.ID:   bitget.New(),
	bybit.ID:    bybit.New(),
	gateio.ID:   gateio.New(),
	htx.ID:      htx.New(),
	kucoin.ID:   kucoin.New(),
	mexc.ID:     mexc.New(),
	okx.ID:      okx.New(),
	poloniex.ID: poloniex.New(),
	ourbit.ID:   ourbit.New(),
	coinex.ID:   coinex.New(),
	bitunix.ID:  bitunix.New(),
	woox.ID:     woox.New(),
	toobit.ID:   toobit.New(),
	phemex.ID:   phemex.New(),
	btcc.ID:     btcc.New(),
	kraken.ID:   kraken.New(),
}

// Get returns the adapter for id (e.g. "binance").
func Get(id string) (exchange.Adapter, error) {
	a, ok := adapters[strings.ToLower(strings.TrimSpace(id))]
	if !ok {
		return nil, fmt.Errorf("unknown exchange %q (supported: %v)", id, Supported())
	}
	return a, nil
}

// Default returns the primary MVP adapter (Binance).
func Default() exchange.Adapter { return adapters[binance.ID] }

// Supported lists registered exchange ids.
func Supported() []string {
	out := make([]string, 0, len(adapters))
	for id := range adapters {
		out = append(out, id)
	}
	return out
}
