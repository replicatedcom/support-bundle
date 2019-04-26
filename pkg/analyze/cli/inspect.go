package cli

import (
	"bytes"
	"context"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/replicatedcom/support-bundle/pkg/analyze/analyze"
	"github.com/replicatedcom/support-bundle/pkg/analyze/render"
	"github.com/replicatedcom/support-bundle/pkg/ui"
	"github.com/replicatedcom/support-bundle/pkg/version"
	"github.com/spf13/cobra"
)

type InspectOptions struct {
	Output string
	Quiet  bool
}

// InspectCmd will inspect a bundle for all its roots
func InspectCmd() *cobra.Command {
	version.Init()

	var opts InspectOptions

	cmd := &cobra.Command{
		Use:   "inspect [BUNDLE]",
		Short: "inspect a troubleshoot bundle archive",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cli := ui.New(nil, cmd.OutOrStdout(), cmd.OutOrStderr())
			return inspectRun(
				context.Background(),
				cli,
				args[0],
				opts)
		},
	}

	// analyze flags
	// cmd.Flags().StringP("collect-bundle-path", "b", "", "path to collect bundle archive") // required
	cmd.Flags().StringVarP(&opts.Output, "output", "o", "json", "output format, one of: json|yaml")
	cmd.Flags().BoolVarP(&opts.Quiet, "quiet", "q", false, "output bundle roots subpaths only")

	return cmd
}

func inspectRun(ctx context.Context, ui cli.Ui, bundlePath string, opts InspectOptions) error {
	bundles, err := analyze.InspectE(ctx, bundlePath)

	if len(bundles) > 0 {
		b := bytes.NewBuffer(nil)
		r := render.New(b, opts.Output)
		if errRender := r.RenderBundles(ctx, bundles, opts.Quiet); errRender != nil {
			err = errRender
		} else {
			ui.Output(strings.TrimSuffix(b.String(), "\n")) // u.Output adds a newline
		}
	}

	return err
}
