package verifier

import (
	"context"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/rpc"
)

const (
	polERC20Ethereum  = "0x455e53CBB86018Ac2B8092FdCd39d8444aFFC3F6"
	maticERC20Ethereum = "0x7D1AfA7B718fb893dF2331a6195D492317827Aa0"
)

// polEthAccounted sums POL + legacy MATIC ERC20 on Ethereum (post-migration holdings).
func polEthAccounted(ctx context.Context, c *rpc.Client, holder string, height int64) (decimal.Decimal, string, map[string]string, error) {
	var total decimal.Decimal
	used := ""
	components := map[string]string{"mode": "pol_plus_matic_erc20"}
	var errs []string

	for label, token := range map[string]string{"pol": polERC20Ethereum, "matic": maticERC20Ethereum} {
		bal, u, err := erc20BalanceOf(ctx, c, rpc.NetEthereum, token, holder, height)
		if err != nil {
			if label == "matic" && strings.Contains(err.Error(), "0 bytes") {
				components[label] = "0"
				components["matic_note"] = "legacy MATIC ERC20 has no code at snapshot block (post-migration)"
				continue
			}
			errs = append(errs, fmt.Sprintf("%s: %v", label, err))
			continue
		}
		used = joinProviders(used, u)
		part := decimal.NewFromBigInt(bal, -18)
		components[label] = part.String()
		total = total.Add(part)
	}
	if used == "" {
		return decimal.Zero, "", nil, fmt.Errorf("pol|eth balance: %s", strings.Join(errs, "; "))
	}
	return total, used, components, nil
}
