package analyze

import (
	"context"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/analyzer"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api/common"
	"github.com/replicatedcom/support-bundle/pkg/analyze/resolver"
	bundleresolver "github.com/replicatedcom/support-bundle/pkg/collect/bundle/resolver"
	collectcli "github.com/replicatedcom/support-bundle/pkg/collect/cli"
	pkgerrors "github.com/replicatedcom/support-bundle/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/spew"
	"github.com/replicatedcom/support-bundle/pkg/version"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var (
	ErrSeverityThreshold = pkgerrors.CmdError{
		ExitCode: 2,
		Err:      errors.New("results did not meet severity threshold"),
	}
)

type Analyze struct {
	Logger log.Logger

	Resolver       *resolver.Resolver
	BundleResolver *bundleresolver.Factory
	Analyzer       *analyzer.Analyzer

	SpecFiles  []string
	Specs      []string
	CustomerID string // deprecated
	ChannelID  string
	Endpoint   string

	// analyze
	SeverityThreshold string
}

// New gets an instance using viper to pull config
func New(v *viper.Viper, logger log.Logger, resolver *resolver.Resolver, bundleResolver *bundleresolver.Factory, analyzer *analyzer.Analyzer) *Analyze {
	return &Analyze{
		Logger: logger,

		Resolver:       resolver,
		BundleResolver: bundleResolver,
		Analyzer:       analyzer,

		SpecFiles:  cast.ToStringSlice(strings.Trim(v.GetString("spec-file"), "[]")),
		Specs:      cast.ToStringSlice(strings.Trim(v.GetString("spec"), "[]")),
		CustomerID: v.GetString("customer-id"),
		Endpoint:   v.GetString("endpoint"),

		// analyze
		SeverityThreshold: v.GetString("severity-threshold"),
	}
}

func (a *Analyze) Execute(ctx context.Context, bundlePath string) ([]api.Result, error) {
	debug := level.Debug(log.With(a.Logger, "method", "Analyze.Execute"))

	debug.Log("method", "configure", "phase", "initialize",
		"version", version.Version(),
		"gitSHA", version.GitSHA(),
		"buildTime", version.BuildTime(),
		"buildTimeFallback", version.GetBuild().TimeFallback,
	)

	if bundlePath == "" {
		// stdin?
		return nil, errors.New("bundle path not specified")
	}

	debug.Log(
		"phase", "resolve")

	endpoint := a.Endpoint
	if endpoint == "" {
		endpoint = collectcli.DefaultEndpoint
	}

	input := resolver.Input{
		Files:      a.SpecFiles,
		Inline:     a.Specs,
		CustomerID: a.CustomerID,
		ChannelID:  a.ChannelID,
		Endpoint:   a.Endpoint,
	}

	spec, err := a.Resolver.ResolveSpec(ctx, input)
	if err != nil {
		debug.Log(
			"phase", "resolve",
			"error", err)
		return nil, errors.Wrap(err, "resolve specs")
	}

	if len(spec.Analyze.V1) == 0 {
		err := errors.New("analyze spec empty") // TODO: typed error
		debug.Log(
			"phase", "resolve",
			"error", err)
		return nil, err
	}

	debug.Log(
		"phase", "resolve",
		"status", "complete")

	debug.Log(
		"phase", "bundle.resolve",
		"bundlePath", bundlePath)

	fs, resolvedPath, err := a.BundleResolver.Fs(bundlePath)
	if err != nil {
		debug.Log(
			"phase", "bundle.resolve",
			"error", err)
		return nil, errors.Wrap(err, "resolve bundle")
	}

	debug.Log(
		"phase", "bundle.resolve",
		"status", "complete")

	debug.Log(
		"phase", "analyze",
		"spec", spew.Sdump(spec.Analyze))

	results, err := a.Analyzer.AnalyzeBundle(
		ctx,
		spec.Analyze,
		fs,
		resolvedPath)
	if err != nil {
		debug.Log(
			"phase", "analyze",
			"error", err)
		return results, errors.Wrap(err, "analyze")
	}

	debug.Log(
		"phase", "analyze",
		"status", "complete")

	if didResultsFailSeverityThreshold(results, common.Severity(a.SeverityThreshold)) {
		return results, ErrSeverityThreshold
	}
	return results, nil
}

func didResultsFailSeverityThreshold(results []api.Result, threshold common.Severity) bool {
	for _, result := range results {
		if common.SeverityCompare(
			result.Severity,
			threshold,
		) >= 0 {
			return true
		}
	}
	return false
}
