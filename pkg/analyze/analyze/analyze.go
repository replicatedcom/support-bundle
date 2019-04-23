package analyze

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/analyzer"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api/common"
	"github.com/replicatedcom/support-bundle/pkg/analyze/resolver"
	collectcli "github.com/replicatedcom/support-bundle/pkg/collect/cli"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
	pkgerrors "github.com/replicatedcom/support-bundle/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/getter"
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

	Resolver *resolver.Resolver
	Getter   getter.Interface
	Analyzer *analyzer.Analyzer

	SpecFiles         []string
	Specs             []string
	SkipDefault       bool
	BundleRootSubpath string
	CustomerID        string // deprecated
	ChannelID         string
	Endpoint          string

	// analyze
	SeverityThreshold string
}

// New gets an instance using viper to pull config
func New(v *viper.Viper, logger log.Logger, resolver *resolver.Resolver, getter *getter.Getter, analyzer *analyzer.Analyzer) *Analyze {
	return &Analyze{
		Logger: logger,

		Resolver: resolver,
		Getter:   getter,
		Analyzer: analyzer,

		SpecFiles:         cast.ToStringSlice(strings.Trim(v.GetString("spec-file"), "[]")),
		Specs:             cast.ToStringSlice(strings.Trim(v.GetString("spec"), "[]")),
		SkipDefault:       v.GetBool("skip-default"),
		BundleRootSubpath: v.GetString("bundle-root-subpath"),
		CustomerID:        v.GetString("customer-id"),
		Endpoint:          v.GetString("endpoint"),

		// analyze
		SeverityThreshold: v.GetString("severity-threshold"),
	}
}

func (a *Analyze) Inspect(ctx context.Context, bundlePath string) (map[string][]collecttypes.Result, error) {
	debug := level.Debug(log.With(a.Logger, "method", "Analyze.Inspect"))

	defer a.deferCleanup()()

	debug.Log("method", "configure", "phase", "initialize",
		"version", version.Version(),
		"gitSHA", version.GitSHA(),
		"buildTime", version.BuildTime(),
		"buildTimeFallback", version.GetBuild().TimeFallback,
	)

	resolvedPath, err := a.Getter.Get(bundlePath)
	debug.Log(
		"phase", "bundle.get",
		"bundlePath", bundlePath,
		"resolvedPath", resolvedPath,
		"error", err)
	if err != nil {
		return nil, errors.Wrap(err, "get bundle")
	}

	debug.Log(
		"phase", "discover")

	bundles, err := a.Analyzer.DiscoverBundles(
		ctx,
		resolvedPath)
	if err != nil {
		debug.Log(
			"phase", "discover",
			"error", err)
		return bundles, errors.Wrap(err, "analyze")
	}

	debug.Log(
		"phase", "discover",
		"status", "complete")

	return bundles, nil
}

func (a *Analyze) Execute(ctx context.Context, bundlePath string) ([]api.Result, error) {
	debug := level.Debug(log.With(a.Logger, "method", "Analyze.Execute"))

	defer a.deferCleanup()()

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

	endpoint := a.Endpoint
	if endpoint == "" {
		endpoint = collectcli.DefaultEndpoint
	}

	input := resolver.Input{
		Files:      a.SpecFiles,
		Inline:     a.Specs,
		CustomerID: a.CustomerID,
		ChannelID:  a.ChannelID,
		Endpoint:   endpoint,
	}
	spec, err := a.Resolver.ResolveSpec(ctx, input, a.SkipDefault)
	debug.Log(
		"phase", "resolve",
		"files", a.SpecFiles,
		"inline", a.Specs,
		"customerID", a.CustomerID,
		"channelID", a.ChannelID,
		"endpoint", endpoint,
		"error", err)
	if err != nil {
		return nil, errors.Wrap(err, "resolve specs")
	}

	if len(spec.Analyze.V1) == 0 {
		err := errors.New("analyze spec empty") // TODO: typed error
		return nil, err
	}

	resolvedPath, err := a.Getter.Get(bundlePath)
	debug.Log(
		"phase", "bundle.get",
		"bundlePath", bundlePath,
		"resolvedPath", resolvedPath,
		"error", err)
	if err != nil {
		return nil, errors.Wrap(err, "get bundle")
	}

	results, err := a.Analyzer.AnalyzeBundle(
		ctx,
		spec.Analyze,
		resolvedPath,
		a.BundleRootSubpath)

	debug.Log(
		"phase", "analyze",
		"resolvedPath", resolvedPath,
		"bundleRootSubpath", a.BundleRootSubpath,
		"error", err)
	if err != nil {
		return results, errors.Wrap(err, "analyze")
	}

	if didResultsFailSeverityThreshold(results, common.Severity(a.SeverityThreshold)) {
		return results, ErrSeverityThreshold
	}
	return results, nil
}

func (a *Analyze) deferCleanup() func() {
	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		a.cleanup()
		os.Exit(1)
	}()

	return func() {
		err := recover() // make sure that we clean up after ourselves no matter what
		a.cleanup()
		if err != nil {
			panic(err)
		}
	}
}

func (a *Analyze) cleanup() error {
	return os.RemoveAll(a.Getter.DstDir())
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
