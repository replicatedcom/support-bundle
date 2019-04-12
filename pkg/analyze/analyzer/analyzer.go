package analyzer

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api/common"
	v1 "github.com/replicatedcom/support-bundle/pkg/analyze/api/v1"
	bundlereader "github.com/replicatedcom/support-bundle/pkg/collect/bundle/reader"
	"github.com/replicatedcom/support-bundle/pkg/spew"
	"github.com/spf13/afero"
)

type Analyzer struct {
	Logger log.Logger
	Fs     afero.Fs
}

func New(logger log.Logger, fs afero.Fs) *Analyzer {
	return &Analyzer{
		Logger: logger,
		Fs:     fs,
	}
}

func (a *Analyzer) AnalyzeBundle(ctx context.Context, spec api.Analyze, archivePath string) ([]api.Result, error) {
	debug := level.Debug(log.With(a.Logger, "method", "Analyzer.AnalyzeBundle"))

	debug.Log(
		"phase", "analyzer.analyze-bundle")

	bundleReader, err := bundlereader.NewBundle(a.Fs, archivePath)
	debug.Log(
		"phase", "analyzer.get-bundle-index",
		"index", spew.Sdump(bundleReader.GetIndex()),
		"error", err)
	if err != nil {
		return nil, errors.Wrapf(err, "new bundle from %s", archivePath)
	}

	var results []api.Result
	var multiErr error
	for _, analyzerSpec := range spec.V1 {
		result, err := a.analyze(ctx, bundleReader, analyzerSpec)
		if err != nil {
			result.Error = fmt.Sprintf("%v", err)
			if result.Severity == "" {
				result.Severity = common.SeverityError
			}
			multiErr = multierror.Append(multiErr, err)
		} else {
			if result.Severity == "" {
				result.Severity = common.SeverityInfo
			}
		}
		results = append(results, result)
	}

	debug.Log(
		"phase", "analyzer.analyze-bundle",
		"status", "complete")
	return results, nil
}

func (a *Analyzer) analyze(ctx context.Context, bundleReader bundlereader.BundleReader, analyzerSpec v1.Analyzer) (api.Result, error) {
	debug := level.Debug(log.With(a.Logger, "method", "Analyzer.analyze"))

	debug.Log(
		"phase", "analyzer.analyze",
		"spec", spew.Sdump(analyzerSpec))

	var result api.Result
	result.Variables = map[string]interface{}{}
	result.AnalyzerSpec = api.Analyze{V1: []v1.Analyzer{analyzerSpec}}

	for _, variable := range analyzerSpec.RegisterVariables {
		reg, err := variable.Register(bundleReader)
		for key, val := range reg {
			result.Variables[key] = val
		}
		debug.Log(
			"phase", "analyzer.analyze.register-variable",
			"variable", spew.Sdump(variable),
			"register", reg,
			"error", err)
		if err != nil {
			return result, errors.Wrapf(err, "register variable %s", variable.Name)
		}
	}

	preconditionsOk, err := analyzerSpec.Preconditions.Eval(nil, result.Variables)
	debug.Log(
		"phase", "analyzer.analyze.eval-preconditions",
		"preconditions", spew.Sdump(analyzerSpec.Preconditions),
		"variables", result.Variables,
		"ok", preconditionsOk,
		"error", err)
	if err != nil {
		if analyzerSpec.Messages.PreconditionError != nil {
			result.Message = analyzerSpec.Messages.PreconditionError
			result.Severity = analyzerSpec.Messages.PreconditionError.Severity
		}
		return result, errors.Wrap(err, "eval preconditions")
	} else if !preconditionsOk {
		if analyzerSpec.Messages.PreconditionFalse != nil {
			result.Message = analyzerSpec.Messages.PreconditionFalse
			result.Severity = analyzerSpec.Messages.PreconditionFalse.Severity
		}
		return result, nil
	}

	conditionsOk, err := analyzerSpec.Conditions.Eval(nil, result.Variables)
	debug.Log(
		"phase", "analyzer.analyze.eval-conditions",
		"conditions", spew.Sdump(analyzerSpec.Conditions),
		"variables", result.Variables,
		"ok", conditionsOk,
		"error", err)
	if err != nil {
		if analyzerSpec.Messages.ConditionError != nil {
			result.Message = analyzerSpec.Messages.ConditionError
			result.Severity = analyzerSpec.Messages.ConditionError.Severity
		}
		return result, errors.Wrap(err, "eval conditions")
	} else if !conditionsOk {
		if analyzerSpec.Messages.ConditionFalse != nil {
			result.Message = analyzerSpec.Messages.ConditionFalse
			result.Severity = analyzerSpec.Messages.ConditionFalse.Severity
		}
		return result, nil
	}

	if analyzerSpec.Messages.ConditionTrue != nil {
		result.Message = analyzerSpec.Messages.ConditionTrue
		result.Severity = analyzerSpec.Messages.ConditionTrue.Severity
	}

	debug.Log(
		"phase", "analyzer.analyze",
		"status", "complete")

	return result, nil
}
