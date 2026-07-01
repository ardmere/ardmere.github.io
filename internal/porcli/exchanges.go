package porcli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ardmere/ardmere/internal/exchangeregistry"
)

func runExchanges(args []string) int {
	path := exchangeregistry.DefaultPath()
	if env := os.Getenv("ARDMERE_EXCHANGE_REGISTRY"); env != "" {
		path = env
	}
	reg, err := exchangeregistry.Load(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load registry: %v\n", err)
		return 1
	}
	if len(args) == 0 {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		exchangeregistry.PrintTable(reg, w)
		return 0
	}
	ex := reg.Find(args[0])
	if ex == nil {
		fmt.Fprintf(os.Stderr, "unknown exchange %q (see: por exchanges)\n", args[0])
		return 1
	}
	exchangeregistry.PrintDetail(*ex)
	return 0
}
