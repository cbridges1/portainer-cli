package main

import (
	"os"

	"github.com/cbridges1/portainer-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
