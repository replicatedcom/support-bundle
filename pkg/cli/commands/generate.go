package commands

import (
	"github.com/replicatedcom/support-bundle/pkg/cli"
	"github.com/spf13/cobra"
)

type generateOptions struct {
	specFiles      []string
	specDocs       []string
	bundlePath     string
	skipDefault    bool
	timeoutSeconds int
}

func NewGenerateCommand(cli *cli.Cli) *cobra.Command {
	opts := generateOptions{}

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a new support bundle",
		Long:  `Collect data and generate a new support bundle`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.Generate(opts.specFiles, opts.specDocs, opts.bundlePath, opts.skipDefault, opts.timeoutSeconds)
		},
	}

	cmd.Flags().StringArrayVarP(&opts.specFiles, "spec-file", "f", nil, "spec file (default is to run core tasks only)")
	cmd.Flags().StringArrayVarP(&opts.specDocs, "spec", "s", nil, "spec doc (default is to run core tasks only)")
	cmd.Flags().StringVarP(&opts.bundlePath, "out", "o", "supportbundle.tar.gz", "Path where the generated bundle should be stored")
	cmd.Flags().BoolVar(&opts.skipDefault, "skip-default", false, "If present, skip the default support bundle files")
	cmd.Flags().IntVar(&opts.timeoutSeconds, "timeout", 60, "The overall support bundle generation timeout")

	return cmd
}
