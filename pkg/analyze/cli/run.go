package cli

import (
	"bytes"
	"context"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/replicatedcom/support-bundle/pkg/analyze/analyze"
	"github.com/replicatedcom/support-bundle/pkg/analyze/render"
	collectcli "github.com/replicatedcom/support-bundle/pkg/collect/cli"
	pkgerrors "github.com/replicatedcom/support-bundle/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/logger"
	"github.com/replicatedcom/support-bundle/pkg/ui"
	"github.com/replicatedcom/support-bundle/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runExample = `
$ analyze run ~/Downloads/supportbundle.tar.gz

$ cat ~/Downloads/supportbundle.tar.gz | analyze run -`

type RunOptions struct {
	Output string
	Quiet  bool
}

// RunCmd will collect and analyze a troubleshoot spec
func RunCmd() *cobra.Command {
	version.Init()

	var opts RunOptions

	cmd := &cobra.Command{
		Use:     "run [BUNDLE]",
		Short:   "Analyze a troubleshoot bundle archive",
		Long:    "Analyze a troubleshoot bundle archive. Arg \"-\" denotes read bundle from stdin.",
		Example: runExample,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			v := viper.GetViper()
			logLevel := logger.GetLevel(v)
			cli := ui.New(nil, cmd.OutOrStdout(), cmd.OutOrStderr())
			if opts.Output == "human" && !v.GetBool("no-color") {
				cli = ui.Colored(cli, v.GetBool("force-color"))
			}
			return analyzeRun(
				context.Background(),
				cli,
				args[0],
				opts,
				logLevel)
		},
	}

	cmd.Flags().StringSliceP("spec-file", "f", nil, "spec file")
	cmd.Flags().StringSliceP("spec", "s", nil, "spec doc")
	cmd.Flags().Bool("skip-default", false, "Skip the default analyze spec")
	cmd.Flags().String("bundle-root-subpath", "", "The subpath within the archive at which the bundle root resides")

	cmd.Flags().String("customer-id", "", "Replicated Customer ID")
	cmd.Flags().MarkDeprecated("customer-id", "This argument is no longer supported. Consider using \"channel-id\"")

	cmd.Flags().String("channel-id", "", "Replicated ChannelID to attempt to get a collector definition from")
	cmd.Flags().String("endpoint", collectcli.DefaultEndpoint, "Endpoint to fetch collector definitions fom")

	// analyze flags
	// cmd.Flags().StringP("collect-bundle-path", "b", "", "path to collect bundle archive") // required
	cmd.Flags().StringVarP(&opts.Output, "output", "o", "human", "output format, one of: human|json|yaml")
	cmd.Flags().BoolVarP(&opts.Quiet, "quiet", "q", false, "suppress normal output")
	cmd.Flags().String("severity-threshold", "error", "the severity threshold at which to exit with an error")

	viper.BindPFlags(cmd.Flags())
	viper.BindPFlags(cmd.PersistentFlags())

	return cmd
}

func analyzeRun(ctx context.Context, ui cli.Ui, bundlePath string, opts RunOptions, logLevel string) error {
	// "-" denotes stdin
	if bundlePath == "-" {
		var err error
		bundlePath, err = bundleFromStdin()
		defer os.Remove(bundlePath)
		if err != nil {
			return err
		}
	}

	results, err := analyze.RunE(ctx, bundlePath)

	if !opts.Quiet && (len(results) > 0 || err == nil) {
		b := bytes.NewBuffer(nil)
		r := render.New(b, opts.Output)
		if errRender := r.RenderResults(ctx, results); errRender != nil {
			err = errRender
		} else {
			ui.Output(strings.TrimSuffix(b.String(), "\n")) // u.Output adds a newline
		}
	}

	if err != nil && err != analyze.ErrSeverityThreshold {
		outputCmdError(err, ui, logLevel)
		err = pkgerrors.CmdError{Err: err}
	}
	return err
}
