package main

import (
	"github.com/replicatedcom/support-bundle/pkg/cli"
	"github.com/spf13/cobra"
)

func NewServerCommand(cli *cli.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start a support bundle server",
		Long:  `Run a support bundle server for collection of continuous data`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
