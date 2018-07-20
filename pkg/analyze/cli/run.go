package cli

import (
	"bytes"
	"context"

	"github.com/mitchellh/cli"
	"github.com/replicatedcom/support-bundle/pkg/analyze/analyze"
	"github.com/replicatedcom/support-bundle/pkg/analyze/render"
	pkgerrors "github.com/replicatedcom/support-bundle/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/ui"
	"github.com/replicatedcom/support-bundle/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RunCmd will collect and analyze a troubleshoot spec
func RunCmd() *cobra.Command {
	version.Init()
	cmd := &cobra.Command{
		Use:   "run [BUNDLE]",
		Short: "collect and analyze troubleshoot spec",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			v := viper.GetViper()
			bundle := args[0]
			return analyzeRun(
				context.Background(),
				bundle,
				ui.New(
					cmd.OutOrStdout(), cmd.OutOrStderr(),
					v.GetBool("force-color"), v.GetBool("no-color"),
				),
				v.GetString("output"),
				v.GetBool("quiet"),
				v.GetString("log-level"))
		},
	}

	cmd.Flags().StringArrayP("spec-file", "f", nil, "spec file")
	cmd.Flags().StringArrayP("spec", "s", nil, "spec doc")

	cmd.Flags().StringP("output", "o", "human", "output format, one of: human|json|yaml")
	cmd.Flags().BoolP("quiet", "q", false, "suppress normal output")
	cmd.Flags().String("severity-threshold", "error", "the severity threshold at which to exit with an error")

	viper.BindPFlags(cmd.Flags())
	viper.BindPFlags(cmd.PersistentFlags())

	return cmd
}

func analyzeRun(
	ctx context.Context,
	bundle string,
	ui cli.Ui,
	outputFormat string,
	quiet bool,
	logLevel string,
) error {
	results, err := analyze.RunE(ctx, bundle)

	if !quiet && len(results) > 0 {
		r := render.New(ui, outputFormat)
		var b bytes.Buffer
		if errRender := r.RenderResults(ctx, &b, results); errRender != nil {
			err = errRender
		} else {
			ui.Output(b.String())
		}
	}

	if err != nil && err != analyze.ErrSeverityThreshold {
		outputCmdError(err, ui, logLevel)
		err = pkgerrors.CmdError{Err: err}
	}
	return err
}
