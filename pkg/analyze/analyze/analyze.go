package analyze

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/analyzer"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api/common"
	"github.com/replicatedcom/support-bundle/pkg/analyze/collector"
	"github.com/replicatedcom/support-bundle/pkg/analyze/resolver"
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

	Resolver  *resolver.Resolver
	Collector collector.Interface
	Analyzer  *analyzer.Analyzer

	SpecFiles  []string
	Specs      []string
	CustomerID string // deprecated
	ChannelID  string
	Endpoint   string

	// analyze
	SeverityThreshold string
	CollectBundlePath string

	// collect
	CollectEnableCore       bool
	CollectEnableDocker     bool
	CollectEnableJournald   bool
	CollectEnableKubernetes bool
	CollectEnableRetraced   bool
	CollectTimeout          time.Duration
	CollectTmpDir           string
}

// New gets an instance using viper to pull config
func New(v *viper.Viper, logger log.Logger, resolver *resolver.Resolver, collector collector.Interface, analyzer *analyzer.Analyzer) *Analyze {
	return &Analyze{
		Logger: logger,

		Resolver:  resolver,
		Collector: collector,
		Analyzer:  analyzer,

		SpecFiles:  cast.ToStringSlice(strings.Trim(v.GetString("spec-file"), "[]")),
		Specs:      cast.ToStringSlice(strings.Trim(v.GetString("spec"), "[]")),
		CustomerID: v.GetString("customer-id"),
		Endpoint:   v.GetString("endpoint"),

		// analyze
		SeverityThreshold: v.GetString("severity-threshold"),
		CollectBundlePath: v.GetString("collect-bundle-path"),

		// collect
		CollectEnableCore:       v.GetBool("collect-core"),
		CollectEnableDocker:     v.GetBool("collect-docker"),
		CollectEnableJournald:   v.GetBool("collect-journald"),
		CollectEnableKubernetes: v.GetBool("collect-kubernetes"),
		CollectEnableRetraced:   v.GetBool("collect-retraced"),
		CollectTimeout:          v.GetDuration("collect-timeout"),
		CollectTmpDir:           v.GetString("collect-temporary-directory"),
	}
}

func (a *Analyze) Execute(ctx context.Context) ([]api.Result, error) {
	debug := level.Debug(log.With(a.Logger, "method", "Analyze.Execute"))

	debug.Log("method", "configure", "phase", "initialize",
		"version", version.Version(),
		"gitSHA", version.GitSHA(),
		"buildTime", version.BuildTime(),
		"buildTimeFallback", version.GetBuild().TimeFallback,
	)

	debug.Log(
		"phase", "resolve")

	endpoint := a.Endpoint
	if endpoint == "" {
		endpoint = collectcli.DefaultEndpoint
	}

	input := resolver.ResolverInput{
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

	bundlePath := a.CollectBundlePath
	if bundlePath == "" {
		bundlePath = filepath.Join(a.CollectTmpDir, "bundle.tgz")
		defer os.RemoveAll(bundlePath)

		err := a.collectBundle(ctx, bundlePath)
		if err != nil {
			return nil, err
		}
	}

	debug.Log(
		"phase", "analyze",
		"spec", spew.Sdump(spec.Analyze))

	results, err := a.Analyzer.AnalyzeBundle(
		ctx,
		spec.Analyze,
		bundlePath)
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

func (a *Analyze) collectBundle(ctx context.Context, dest string) error {
	debug := level.Debug(log.With(a.Logger, "method", "Analyze.collect"))

	debug.Log(
		"phase", "bundle.generate",
		"timeout", a.CollectTimeout,
		"dest", dest)

	input := collector.CollectorInput{
		CustomerID: a.CustomerID,
		ChannelID:  a.ChannelID,
		Specs:      a.Specs,
		SpecFiles:  a.SpecFiles,
		Dest:       dest,
		Opts: collector.Options{
			EnableCore:       a.CollectEnableCore,
			EnableDocker:     a.CollectEnableDocker,
			EnableJournald:   a.CollectEnableJournald,
			EnableKubernetes: a.CollectEnableKubernetes,
			EnableRetraced:   a.CollectEnableRetraced,
			Timeout:          a.CollectTimeout,
			Endpoint:         a.Endpoint,
		},
	}

	err := a.Collector.CollectBundle(ctx, input)
	if err != nil {
		debug.Log(
			"phase", "bundle.generate",
			"error", err)
	}

	debug.Log(
		"phase", "bundle.generate",
		"status", "complete")

	return errors.Wrap(err, "generate bundle")
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
