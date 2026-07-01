package porprobe

import (
	"fmt"
	"os"
)

// Main runs probe subcommands; args[0] is rpc, tron, stakehub, or sonic-sfc.
func Main(args []string) int {
	if len(args) < 1 {
		usage()
		return 2
	}
	switch args[0] {
	case "rpc":
		RunRPC(args[1:])
	case "tron":
		RunTron(args[1:])
	case "stakehub":
		RunStakehub(args[1:])
	case "sonic-sfc":
		RunSonicSFC(args[1:])
	case "help", "-h", "--help":
		usage()
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown probe %q\n\n", args[0])
		usage()
		return 2
	}
	return 0
}

func usage() {
	fmt.Fprintf(os.Stderr, `usage: por probe <kind> [flags]

Kinds:
  rpc        probe EVM archive RPC providers at a historical block
  tron       probe Tron TRC20 balance at a historical block
  stakehub   aggregate BSC Stake Hub delegations for one address
  sonic-sfc  aggregate Sonic SFC delegations for one address

Examples:
  por probe rpc -network BSC -height 101590091 -chainlist
  por probe tron -holder TJ5usJLLwjwn7Pw3TPbdzreG7dvgKzfQ5y
`)
}
