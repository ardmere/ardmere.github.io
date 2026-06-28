package porprobe

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/ardmere/ardmere/internal/rpc"
	"github.com/ardmere/ardmere/internal/verifier"

	"github.com/shopspring/decimal"
)

func RunStakehub(args []string) {
	fs := flag.NewFlagSet("stakehub", flag.ExitOnError)
	address := fs.String("address", "0x86523c87c8ec98c7539e2c58cd813ee9d1a08d96", "delegator address")
	height := fs.Int64("height", 101590091, "BSC block height")
	claim := fs.String("claim", "172797.6409", "optional CSV claim for comparison")
	_ = fs.Parse(args)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	rpcClient := rpc.New()
	liquid, used, err := rpcClient.GetBalance(ctx, rpc.NetBSC, *address, *height)
	if err != nil {
		log.Fatalf("liquid balance: %v", err)
	}
	liquidDec := decimal.NewFromBigInt(liquid, -18)
	log.Printf("liquid=%s BNB (provider=%s)", liquidDec, used)

	staked, err := verifier.StakeHubAccounted(ctx, rpcClient, *address, *height)
	if err != nil {
		log.Fatalf("stake hub: %v", err)
	}
	total := liquidDec.Add(staked.Staked).Add(staked.Locked)
	log.Printf("staked=%s unbonding=%s total=%s validators=%d provider=%s",
		staked.Staked, staked.Locked, total, staked.Scanned, staked.Used)
	if staked.Note != "" {
		log.Printf("note: %s", staked.Note)
	}
	if *claim != "" {
		claimDec, err := decimal.NewFromString(*claim)
		if err == nil {
			diff := total.Sub(claimDec)
			log.Printf("csv claim=%s diff=%s", claimDec, diff)
		}
	}
}
