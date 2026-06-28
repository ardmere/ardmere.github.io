package verifier

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ardmere/ardmere/internal/rpc"
)

func TestPolEthProbe0xa64b(t *testing.T) {
	addr := "0xa64b436964e7415c0e70b9989a53e1fb9a90e726"
	ethH := int64(25218797)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	c := rpc.New()

	for label, token := range map[string]string{
		"pol":   polERC20Ethereum,
		"matic": maticERC20Ethereum,
	} {
		bal, used, err := erc20BalanceOf(ctx, c, rpc.NetEthereum, token, addr, ethH)
		if err != nil {
			t.Logf("%s ERR: %v", label, err)
			continue
		}
		t.Logf("%s=%s (%s)", label, formatWei18(bal), used)
	}

	ethWei, used, err := c.GetBalance(ctx, rpc.NetEthereum, addr, ethH)
	if err != nil {
		t.Logf("eth native ERR: %v", err)
	} else {
		t.Logf("eth native=%s wei (%s)", ethWei, used)
	}

	for _, h := range []int64{87732730, ethH} {
		polWei, used, err := c.GetBalance(ctx, rpc.NetPolygon, addr, h)
		if err != nil {
			t.Logf("pol native MATIC#%d ERR: %v", h, err)
			continue
		}
		t.Logf("pol native MATIC#%d=%s wei (%s)", h, polWei, used)
	}

	total, used, comps, err := polEthAccounted(ctx, c, addr, ethH)
	if err != nil {
		t.Fatalf("polEthAccounted: %v", err)
	}
	t.Logf("polEthAccounted=%s (%s) comps=%v", total, used, comps)
}

func formatWei18(bi *big.Int) string {
	if bi == nil {
		return "0"
	}
	r := new(big.Rat).SetInt(bi)
	r.Quo(r, big.NewRat(1e18, 1))
	return r.FloatString(18)
}
