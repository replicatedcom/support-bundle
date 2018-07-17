package cli

import (
	"bytes"
	"context"
	"os"
	"time"

	"github.com/mitchellh/cli"
	"github.com/replicatedcom/support-bundle/pkg/analyze/analyze"
	"github.com/replicatedcom/support-bundle/pkg/analyze/render"
	"github.com/replicatedcom/support-bundle/pkg/ui"
	"github.com/replicatedcom/support-bundle/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RunCmd will collect and analyze a troubleshoot spec
func RunCmd() *cobra.Command {
	version.Init()
	cmd := &cobra.Command{
		Use:     "run",
		Short:   "collect and analyze troubleshoot spec",
		Version: getVersionString(),
		Run: func(cmd *cobra.Command, args []string) {
			v := viper.GetViper()
			analyzeRun(
				context.Background(),
				ui.FromViper(v),
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
	cmd.Flags().Duration("collect-timeout", 120*time.Second, "collect step timeout")
	cmd.Flags().String("collect-tmp-dir", os.TempDir(), "collect step temporary directory")

	viper.BindPFlags(cmd.Flags())
	viper.BindPFlags(cmd.PersistentFlags())

	return cmd
}

func analyzeRun(
	ctx context.Context,
	ui cli.Ui,
	outputFormat string,
	quiet bool,
	logLevel string,
) {
	results, err := analyze.RunE(ctx)

	if !quiet && len(results) > 0 {
		r := render.New(ui, outputFormat)
		var b bytes.Buffer
		if errRender := r.RenderResults(ctx, &b, results); errRender != nil {
			err = errRender
		} else {
			ui.Output(b.String())
		}
	}

	switch {
	case err == analyze.ErrAnalysisFailed:
		os.Exit(1)
	case err != nil:
		outputCmdError(err, ui, logLevel)
		os.Exit(2)
	}
}
