package commands

import (
	"github.com/replicatedcom/support-bundle/pkg/cli"
	"github.com/spf13/cobra"
)

type indexOptions struct {
	specFiles   []string
	specDocs    []string
	skipDefault bool
	format      string
}

func NewIndexCommand(cli *cli.Cli) *cobra.Command {
	opts := indexOptions{}

	cmd := &cobra.Command{
		Use:   "index",
		Short: "Generate support bundle index document",
		Long:  `Generate a support bundle index json document comprising of paths and descriptions of all support bundle files.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.Index(opts.specFiles, opts.specDocs, opts.skipDefault, opts.format)
		},
	}

	cmd.Flags().StringArrayVarP(&opts.specFiles, "spec-file", "f", nil, "spec file (default is to run core tasks only)")
	cmd.Flags().StringArrayVarP(&opts.specDocs, "spec", "s", nil, "spec doc (default is to run core tasks only)")
	cmd.Flags().BoolVar(&opts.skipDefault, "skip-default", false, "If present, skip the default support bundle files")
	cmd.Flags().StringVar(&opts.format, "format", "json", `Index format (one of "json", "yaml")`)

	return cmd
}
