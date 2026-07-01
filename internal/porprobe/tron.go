package porprobe

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ardmere/ardmere/internal/rpc"
)

func RunTron(args []string) {
	fs := flag.NewFlagSet("tron", flag.ExitOnError)
	contract := fs.String("contract", "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t", "TRC20 contract base58")
	holder := fs.String("holder", "TJ5usJLLwjwn7Pw3TPbdzreG7dvgKzfQ5y", "holder base58")
	height := fs.Int64("height", 83201055, "Tron block number")
	_ = fs.Parse(args)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	c := rpc.New()
	bal, used, err := c.TRC20BalanceOf(ctx, *contract, *holder, *height)
	if err != nil {
		log.Fatalf("tron: %v", err)
	}
	fmt.Printf("balance=%s (raw wei) provider=%s\n", bal, used)
}
