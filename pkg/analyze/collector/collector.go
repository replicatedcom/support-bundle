package collector

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/collect/bundle"
	"github.com/replicatedcom/support-bundle/pkg/collect/cli"
)

type Collector struct {
	Logger log.Logger
}

type CollectorInput struct {
	CustomerID string
	ChannelID  string
	Specs      []string
	SpecFiles  []string
	Dest       string
	Opts       Options
}

func New(logger log.Logger) Interface {
	return &Collector{
		Logger: logger,
	}
}

func (c *Collector) CollectBundle(ctx context.Context, input CollectorInput) error {
	debug := level.Debug(log.With(c.Logger, "method", "Collector.Collect"))

	cmd := cli.NewCli()
	err := cmd.Generate(cli.GenerateOptions{
		CfgFiles:         input.SpecFiles,
		CfgDocs:          input.Specs,
		BundlePath:       input.Dest,
		SkipDefault:      true,
		TimeoutSeconds:   int(input.Opts.Timeout.Seconds()),
		CustomerID:       input.CustomerID,
		ChannelID:        input.ChannelID,
		Endpoint:         input.Opts.Endpoint,
		DenyUploadPrompt: true,
		Quiet:            true,
		PlannerOptions: bundle.PlannerOptions{
			EnableCore:       input.Opts.EnableCore,
			EnableDocker:     input.Opts.EnableDocker,
			EnableJournald:   input.Opts.EnableJournald,
			EnableKubernetes: input.Opts.EnableKubernetes,
			EnableRetraced:   input.Opts.EnableRetraced,

			RequireJournald:   input.Opts.EnableJournald,
			RequireKubernetes: input.Opts.EnableKubernetes,
			RequireRetraced:   input.Opts.EnableRetraced,
		},
	})

	debug.Log(
		"phase", "bundle.generate",
		"timeout", input.Opts.Timeout,
		"dest", input.Dest,
		"error", err)

	return errors.Wrap(err, "generate")
}
