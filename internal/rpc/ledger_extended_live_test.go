package rpc

import (
	"context"
	"os"
	"testing"
)

func TestLedgerExtendedLiveProbe(t *testing.T) {
	if os.Getenv("LEDGER_LIVE_PROBE") == "" {
		t.Skip("set LEDGER_LIVE_PROBE=1 to run live ledger probes")
	}
	c := NewLedger()
	ctx := context.Background()

	cases := []struct {
		name string
		run  func() (int64, error)
	}{
		{
			"ALGO",
			func() (int64, error) {
				v, _, err := c.AlgoASABalanceAtRound(ctx, "QYXDGS2XJJT7QNR6EJ2YHNZFONU6ROFM6BKTBNVT63ZXQ5OC6IYSPNDJ4U", 31566704, 61720069)
				return v, err
			},
		},
		{
			"CAKE",
			func() (int64, error) {
				v, _, err := c.AptosCoinBalanceByTypeAtVersion(ctx, "0xae1a6f3d3daccaf77b55044cea133379934bba04a11b9d0bbd643eae5e6e9c70", "0x159df6b7689437016108a019fd5bef736bac692b6d4a1f10c941f6fbb9a74ca6::oft::CakeOFT", 800530436)
				return v, err
			},
		},
		{
			"STATEMINT USDC",
			func() (int64, error) {
				v, _, err := c.SubstrateAssetBalanceAtBlock(ctx, "https://polkadot-asset-hub-rpc.polkadot.io", "13vg3Mrxm3GL9eXxLsGgLYRueiwFCiMbkdHBL4ZN5aob5D4N", 1337, 16481516)
				return v, err
			},
		},
		{
			"XLM USDC",
			func() (int64, error) {
				v, _, err := c.StellarAssetBalanceAtLedger(ctx, "GC5LF63GRVIT5ZXXCXLPI3RX2YXKJQFZVBSAO6AUELN3YIMSWPD6Z6FH", "USDC", "GA5ZSEJYB37JRC5AVCIA5MOP4RHTM335X2KGX3IHOJAPP5RE34K4KZVN", 62825389)
				return v, err
			},
		},
		{
			"XTZ USDT",
			func() (int64, error) {
				v, _, err := c.TezosFABalanceAtLevel(ctx, "tz1Q3jvYU9knekDYJfyvj3GjUy6898MNjvb2", "KT1XnTn74bUtxHfDtBmm2bGZAQfhPbvKWR8o", 13446574)
				return v, err
			},
		},
		{
			"RLUSD",
			func() (int64, error) {
				v, _, err := c.XRPLEmittedBalanceAtLedger(ctx, "rDAE53VfMvftPB4ogpWGWvzkQxfht6JPxr", "rMxCKbEDwqr76QuheSUMdEGf4B9xJ8m5De", "524C555344000000000000000000000000000000", 104615295, 6)
				return v, err
			},
		},
	}
	for _, tc := range cases {
		v, err := tc.run()
		if err != nil {
			t.Errorf("%s: %v", tc.name, err)
			continue
		}
		if v <= 0 {
			t.Errorf("%s: zero balance", tc.name)
		}
		t.Logf("%s raw=%d", tc.name, v)
	}
}
