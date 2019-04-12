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
	"github.com/replicatedcom/support-bundle/pkg/analyze/message"
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
	var multiErr *multierror.Error
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
	return results, multiErr.ErrorOrNil()
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
			result.Severity = analyzerSpec.Messages.PreconditionError.Severity
			result.Message, err = message.ExecuteTemplates(analyzerSpec.Messages.PreconditionError, result.Variables)
			if err != nil {
				return result, errors.Wrap(err, "execute precondition error message template")
			}
		}
		return result, errors.Wrap(err, "eval preconditions")
	} else if !preconditionsOk {
		if analyzerSpec.Messages.PreconditionFalse != nil {
			result.Severity = analyzerSpec.Messages.PreconditionFalse.Severity
			result.Message, err = message.ExecuteTemplates(analyzerSpec.Messages.PreconditionFalse, result.Variables)
			if err != nil {
				return result, errors.Wrap(err, "execute precondition false message template")
			}
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
			result.Severity = analyzerSpec.Messages.ConditionError.Severity
			result.Message, err = message.ExecuteTemplates(analyzerSpec.Messages.ConditionError, result.Variables)
			if err != nil {
				return result, errors.Wrap(err, "execute condition error message template")
			}
		}
		return result, errors.Wrap(err, "eval conditions")
	} else if !conditionsOk {
		if analyzerSpec.Messages.ConditionFalse != nil {
			result.Severity = analyzerSpec.Messages.ConditionFalse.Severity
			result.Message, err = message.ExecuteTemplates(analyzerSpec.Messages.ConditionFalse, result.Variables)
			if err != nil {
				return result, errors.Wrap(err, "execute condition false message template")
			}
		}
		return result, nil
	}

	if analyzerSpec.Messages.ConditionTrue != nil {
		result.Severity = analyzerSpec.Messages.ConditionTrue.Severity
		result.Message, err = message.ExecuteTemplates(analyzerSpec.Messages.ConditionTrue, result.Variables)
		if err != nil {
			fmt.Println("+++", err)
			return result, errors.Wrap(err, "execute condition true message template")
		}
	}

	debug.Log(
		"phase", "analyzer.analyze",
		"status", "complete")

	return result, nil
}
