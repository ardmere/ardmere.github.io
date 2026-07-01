// Package porcli is the unified PoR command-line entry (por anchor|verify|fetch|probe).
package porcli

import (
	"fmt"
	"os"

	"github.com/ardmere/ardmere/internal/porfetch"
	"github.com/ardmere/ardmere/internal/porprobe"
	"github.com/ardmere/ardmere/internal/porrun"
)

func init() {
	porrun.ValidateAnchorSelector()
}

// Main dispatches por subcommands; args[0] is the program name.
func Main(args []string) int {
	if len(args) < 2 {
		usage()
		return 2
	}
	switch args[1] {
	case "anchor":
		return porrun.RunAnchor(args[2:])
	case "verify":
		return porrun.RunVerify(args[2:])
	case "fetch":
		return porfetch.Main(args[2:])
	case "probe":
		return porprobe.Main(args[2:])
	case "exchanges":
		return runExchanges(args[2:])
	case "batch":
		return runBatch(args[2:])
	case "report":
		return runReport(args[2:])
	case "help", "-h", "--help":
		usage()
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", args[1])
		usage()
		return 2
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `usage: por <command> [flags]

Commands:
  anchor   fetch → verify → write bundles → print anchor calldata
  verify   re-verify cached artifact bundle
  fetch    archive-only download (gateio | htx | bybit | bitget | okx)
  batch    fetch+verify last N snapshots per exchange and write reports
  report   generate assessment + transparency report from cached artifacts
  probe    chain/RPC diagnostics (rpc | tron | stakehub | sonic-sfc)
  exchanges  list official PoR upstream pages and GitHub repos (registry)

Examples:
  go run ./cmd/por exchanges
  go run ./cmd/por exchanges okx
  go run ./cmd/por anchor -exchange okx
  go run ./cmd/por verify -exchange binance -snapshot PR01JUN26
  go run ./cmd/por fetch htx -zk-bundle ./public-data.zip
  go run ./cmd/por probe rpc -network BSC -chainlist
`)
}
