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

func New(
	logger log.Logger,
) Interface {
	return &Collector{
		Logger: logger,
	}
}

func (c *Collector) CollectBundle(
	ctx context.Context,
	customerID string,
	specs []string,
	specFiles []string,
	dest string,
	opts Options,
) error {

	debug := level.Debug(log.With(c.Logger, "method", "Collector.Collect"))

	cmd := cli.NewCli()
	err := cmd.Generate(cli.GenerateOptions{
		CfgFiles:         specFiles,
		CfgDocs:          specs,
		BundlePath:       dest,
		SkipDefault:      true,
		TimeoutSeconds:   int(opts.Timeout.Seconds()),
		CustomerID:       customerID,
		CustomerEndpoint: opts.CustomerEndpoint,
		DenyUploadPrompt: true,
		Quiet:            true,
		PlannerOptions: bundle.PlannerOptions{
			EnableCore:       opts.EnableCore,
			EnableDocker:     opts.EnableDocker,
			EnableJournald:   opts.EnableJournald,
			EnableKubernetes: opts.EnableKubernetes,
			EnableRetraced:   opts.EnableRetraced,

			RequireJournald:   opts.EnableJournald,
			RequireKubernetes: opts.EnableKubernetes,
			RequireRetraced:   opts.EnableRetraced,
		},
	})

	debug.Log(
		"phase", "bundle.generate",
		"timeout", opts.Timeout,
		"dest", dest,
		"error", err)

	return errors.Wrap(err, "generate")
}
