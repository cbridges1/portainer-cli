package main

import (
	"os"

	"github.com/jalenbridges/portainer-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
