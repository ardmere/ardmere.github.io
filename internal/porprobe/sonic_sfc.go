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

func RunSonicSFC(args []string) {
	fs := flag.NewFlagSet("sonic-sfc", flag.ExitOnError)
	address := fs.String("address", "0x64de13c46f627d9c86212050d48756fb65c06d8a", "delegator address")
	height := fs.Int64("height", 72001540, "Sonic block height")
	claim := fs.String("claim", "82000048.8", "optional CSV claim for comparison")
	_ = fs.Parse(args)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	rpcClient := rpc.New()
	out, err := verifier.SonicSFCAccounted(ctx, rpcClient, *address, *height)
	if err != nil {
		log.Fatalf("sonic sfc: %v", err)
	}
	total := out.Liquid.Add(out.Staked)
	log.Printf("liquid=%s S staked=%s S total=%s validators_scanned=%d provider=%s",
		out.Liquid, out.Staked, total, out.Scanned, out.Used)
	for k, v := range out.ByValidator {
		log.Printf("  %s=%s S", k, v)
	}
	if *claim != "" {
		claimDec, err := decimal.NewFromString(*claim)
		if err == nil {
			log.Printf("csv claim=%s diff=%s", claimDec, total.Sub(claimDec))
		}
	}
}
