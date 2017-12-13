package main

import (
	"os"

	"github.com/replicatedcom/support-bundle/cmd/support-bundle/commands"
	"github.com/replicatedcom/support-bundle/pkg/cli"
)

func main() {
	c := commands.NewSupportBundleCommand(cli.NewCli())
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
