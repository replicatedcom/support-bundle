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
	"github.com/replicatedcom/support-bundle/pkg/logger"
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
			logLevel := logger.GetLevel(v)
			cli := ui.New(nil, cmd.OutOrStdout(), cmd.OutOrStderr())
			if v.GetString("output") == "human" && !v.GetBool("no-color") {
				cli = ui.Colored(cli, v.GetBool("force-color"))
			}
			return analyzeRun(
				context.Background(),
				cli,
				v.GetString("output"),
				v.GetBool("quiet"),
				logLevel)
		},
	}

	cmd.Flags().StringArrayP("spec-file", "f", nil, "spec file")
	cmd.Flags().StringArrayP("spec", "s", nil, "spec doc")

	cmd.Flags().String("customer-id", "", "Replicated Customer ID")
	cmd.Flags().MarkDeprecated("customer-id", "This argument is no longer supported. Consider using \"channel-id\"")

	cmd.Flags().String("channel-id", "", "Replicated ChannelID to attempt to get a collector definition from")
	cmd.Flags().String("endpoint", collectcli.DefaultEndpoint, "Endpoint to fetch collector definitions fom")

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

func analyzeRun(ctx context.Context, ui cli.Ui, outputFormat string, quiet bool, logLevel string) error {
	results, err := analyze.RunE(ctx)

	if !quiet && len(results) > 0 {
		b := bytes.NewBuffer(nil)
		r := render.New(b, outputFormat)
		if errRender := r.RenderResults(ctx, results); errRender != nil {
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
