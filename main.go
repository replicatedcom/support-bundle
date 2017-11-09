package main

import (
	"os"

	"github.com/replicatedcom/support-bundle/cmd"
	"github.com/replicatedcom/support-bundle/pkg/cli"
)

func main() {
	c := cmd.NewSupportBundleCommand(cli.NewCli())
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
