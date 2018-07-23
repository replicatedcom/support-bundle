package cli

import (
	"bytes"
	"context"
	"os"
	"time"

	"github.com/mitchellh/cli"
	"github.com/replicatedcom/support-bundle/pkg/analyze/analyze"
	"github.com/replicatedcom/support-bundle/pkg/analyze/render"
	collectcli "github.com/replicatedcom/support-bundle/pkg/collect/cli"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			v := viper.GetViper()
			return analyzeRun(
				context.Background(),
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
	cmd.Flags().String("customer-id", "", "Replicated Customer ID")
	cmd.Flags().String("customer-endpoint", collectcli.DefaultCustomerEndpoint, "Replicated customer API endpoint")

	// analyze flags
	cmd.Flags().StringP("output", "o", "human", "output format, one of: human|json|yaml")
	cmd.Flags().BoolP("quiet", "q", false, "suppress normal output")
	cmd.Flags().String("severity-threshold", "error", "the severity threshold at which to exit with an error")
	cmd.Flags().StringP("collect-bundle-path", "b", "", "collect bundle path (will override any collect spec)")

	// generate flags
	cmd.Flags().Bool("collect-core", true, "enable Core plugin")
	cmd.Flags().Bool("collect-docker", false, "enable Docker plugin")
	cmd.Flags().Bool("collect-journald", false, "enable Journald plugin")
	cmd.Flags().Bool("collect-kubernetes", false, "enable Kubernetes plugin")
	cmd.Flags().Bool("collect-retraced", false, "enable Retraced plugin")
	cmd.Flags().Duration("collect-timeout", collectcli.DefaultGenerateTimeoutSeconds*time.Second, "collect step timeout")
	cmd.Flags().String("collect-temporary-directory", os.TempDir(), "collect step temporary directory")

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
) error {

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

	if err != nil && err != analyze.ErrSeverityThreshold {
		outputCmdError(err, ui, logLevel)
		err = pkgerrors.CmdError{Err: err}
	}
	return err
}
