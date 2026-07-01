// por is the unified Proof-of-Reserves CLI.
package main

import (
	"os"

	"github.com/ardmere/ardmere/internal/porcli"
)

func main() {
	os.Exit(porcli.Main(os.Args))
}
