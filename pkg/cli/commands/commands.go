package commands

import (
	"github.com/replicatedcom/support-bundle/pkg/cli"
	"github.com/spf13/cobra"
)

func NewSupportBundleCommand(cli *cli.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "support-bundle",
		Short: "Generate and manage support bundles",
		Long: `A support bundle is an archive of files, output, metrics and state 
	from a server that can be used to assist when troubleshooting a server.
	
	The support-bundle utility can generate human readable archives, and can also 
	be used to generate input to the support.io service.`,
		SilenceUsage: true,
	}

	addCommands(cmd, cli)

	return cmd
}

func addCommands(cmd *cobra.Command, cli *cli.Cli) {
	cmd.AddCommand(NewGenerateCommand(cli))
	cmd.AddCommand(NewUploadCommand(cli))
	cmd.AddCommand(NewIndexCommand(cli))
}
