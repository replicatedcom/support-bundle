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

	Resolver *resolver.Resolver
	Analyzer *analyzer.Analyzer

	// required
	SeverityThreshold string
	SpecFiles         []string
	Specs             []string

	// optional
	CustomerID string
}

// New gets an instance using viper to pull config
func New(
	v *viper.Viper,
	logger log.Logger,
	resolver *resolver.Resolver,
	analyzer *analyzer.Analyzer,
) *Analyze {
	return &Analyze{
		Logger: logger,

		Resolver: resolver,
		Analyzer: analyzer,

		SeverityThreshold: v.GetString("severity-threshold"),
		SpecFiles:         cast.ToStringSlice(strings.Trim(v.GetString("spec-file"), "[]")),
		Specs:             cast.ToStringSlice(strings.Trim(v.GetString("spec"), "[]")),
		CustomerID:        v.GetString("customer-id"),
	}
}

func (a *Analyze) Execute(ctx context.Context, bundle string) ([]api.Result, error) {
	debug := level.Debug(log.With(a.Logger, "method", "Execute"))

	debug.Log("method", "configure", "phase", "initialize",
		"version", version.Version(),
		"gitSHA", version.GitSHA(),
		"buildTime", version.BuildTime(),
		"buildTimeFallback", version.GetBuild().TimeFallback,
	)

	debug.Log(
		"phase", "resolve")

	spec, err := a.Resolver.ResolveSpec(
		ctx,
		a.SpecFiles,
		a.Specs)
	if err != nil {
		debug.Log(
			"phase", "resolve",
			"error", err)
		return nil, errors.Wrap(err, "resolve specs")
	}

	debug.Log(
		"phase", "resolve",
		"status", "complete")

	debug.Log(
		"phase", "analyze",
		"spec", spew.Sdump(spec.Analyze))

	results, err := a.Analyzer.AnalyzeBundle(
		ctx,
		spec.Analyze,
		bundle)
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
