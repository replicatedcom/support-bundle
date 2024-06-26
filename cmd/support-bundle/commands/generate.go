package commands

import (
	"github.com/replicatedcom/support-bundle/pkg/collect/cli"
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
	cmd.Flags().IntVar(&opts.TimeoutSeconds, "timeout", cli.DefaultGenerateTimeoutSeconds, "The overall support bundle generation timeout")
	cmd.Flags().BoolVar(&opts.EnableCore, "core", true, "Enable Core plugin")
	cmd.Flags().BoolVar(&opts.EnableDocker, "docker", true, "Enable Docker plugin")
	cmd.Flags().BoolVar(&opts.EnableJournald, "journald", true, "Enable Journald plugin")
	cmd.Flags().BoolVar(&opts.RequireJournald, "require-journald", false, "Require Journald plugin")
	cmd.Flags().BoolVar(&opts.EnableKubernetes, "kubernetes", true, "Enable Kubernetes plugin")
	cmd.Flags().BoolVar(&opts.RequireKubernetes, "require-kubernetes", false, "Require Kubernetes plugin")
	cmd.Flags().BoolVar(&opts.EnableRetraced, "retraced", true, "Enable Retraced plugin")
	cmd.Flags().BoolVar(&opts.RequireRetraced, "require-retraced", false, "Require Retraced plugin")
	cmd.Flags().BoolVar(&opts.SkipDefault, "skip-default", false, "If present, skip the default support bundle files")
	cmd.Flags().BoolVarP(&opts.ConfirmUploadPrompt, "yes-upload", "u", false, "If present, auto-confirm any upload prompts")
	cmd.Flags().BoolVar(&opts.DenyUploadPrompt, "no-upload", false, "If present, auto-deny any upload prompts")

	cmd.Flags().StringVar(&opts.Endpoint, "endpoint", cli.DefaultEndpoint, "Customer API Endpoint")

	cmd.Flags().StringVar(&opts.ChannelID, "channel-id", "", "Replicated ChannelID to attempt to get a collector definition from")

	//--out - works, and its totally interactive and everything,
	// and the bundle just gets dumped to stdout.
	// so if you pipe it to a file, you can still read the messages and answer the prompts.
	//
	// This --quiet is to get around a kubectl thing where it combines stderr and stdout,
	// I've tried all kinds of /bin/sh -c wrapping but its a pain and gets really ugly.
	// --quiet is really just a convenience thing for those cases.
	cmd.Flags().BoolVar(&opts.Quiet, "quiet", false, "If set, supress all non-error output and messages. Useful when combined with '--out -' to dump a tarball to stdout.")

	return cmd
}
