package commands

import (
	"github.com/replicatedcom/support-bundle/pkg/cli"
	"github.com/spf13/cobra"
)

func NewGenerateCommand(supportBundle *cli.Cli) *cobra.Command {
	var opts cli.GenerateOptions

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a new support bundle",
		Long:  `Collect data and generate a new support bundle`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return supportBundle.Generate(opts)
		},
	}

	cmd.Flags().StringArrayVarP(&opts.CfgFiles, "spec-file", "f", nil, "spec file (default is to run core tasks only)")
	cmd.Flags().StringArrayVarP(&opts.CfgDocs, "spec", "s", nil, "spec doc (default is to run core tasks only)")
	cmd.Flags().StringVarP(&opts.BundlePath, "out", "o", "supportbundle.tar.gz", "Path where the generated bundle should be stored")
	cmd.Flags().IntVar(&opts.TimeoutSeconds, "timeout", 60, "The overall support bundle generation timeout")
	cmd.Flags().BoolVar(&opts.EnableCore, "core", true, "Enable Core plugin")
	cmd.Flags().BoolVar(&opts.EnableDocker, "docker", true, "Enable Docker plugin")
	cmd.Flags().BoolVar(&opts.EnableJournald, "journald", false, "Enable Journald plugin")
	cmd.Flags().BoolVar(&opts.EnableKubernetes, "kubernetes", false, "Enable Kubernetes plugin")
	cmd.Flags().BoolVar(&opts.EnableRetraced, "retraced", false, "Enable Retraced plugin")
	cmd.Flags().StringVar(&opts.CustomerID, "customer-id", "", "Replicated Customer ID")

	return cmd
}
