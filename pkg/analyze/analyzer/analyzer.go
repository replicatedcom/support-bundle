package analyzer

import (
	"context"
	"encoding/json"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api/common"
	v1 "github.com/replicatedcom/support-bundle/pkg/analyze/api/v1"
	"github.com/replicatedcom/support-bundle/pkg/analyze/message"
	bundlereader "github.com/replicatedcom/support-bundle/pkg/collect/bundle/reader"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
	"github.com/replicatedcom/support-bundle/pkg/util"
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

func (a *Analyzer) DiscoverBundles(ctx context.Context, archivePath string) (map[string][]collecttypes.Result, error) {
	debug := level.Debug(log.With(a.Logger, "method", "Analyzer.DiscoverBundles"))

	debug.Log(
		"phase", "analyzer.discover-bundles")

	bundleReader, err := bundlereader.NewMultiBundle(a.Fs, archivePath)
	debug.Log(
		"phase", "analyzer.discover-bundles",
		"len", len(bundleReader.GetBundles()),
		"error", err)
	if err != nil {
		return nil, errors.Wrapf(err, "discover bundles from %s", archivePath)
	}

	bundles := map[string][]collecttypes.Result{}
	for prefix, bundle := range bundleReader.GetBundles() {
		bundles[prefix] = bundle.GetIndex()
	}
	return bundles, nil
}

func (a *Analyzer) AnalyzeBundle(ctx context.Context, spec api.Analyze, archivePath, bundleRootSubpath string) ([]api.Result, error) {
	debug := level.Debug(log.With(a.Logger, "method", "Analyzer.AnalyzeBundle"))

	debug.Log(
		"phase", "analyzer.analyze-bundle")

	bundleReader, err := bundlereader.NewBundle(a.Fs, archivePath, bundleRootSubpath)
	debug.Log(
		"phase", "analyzer.get-bundle-index",
		// "index", util.SpewJSON(bundleReader.GetIndex()), // TOO NOISY
		"error", err)
	if err != nil {
		return nil, errors.Wrapf(err, "new bundle from %s", archivePath)
	}

	var results []api.Result
	var multiErr *multierror.Error
	for _, analyzerSpec := range spec.V1 {
		result, err := a.analyze(ctx, bundleReader, analyzerSpec)
		if err != nil {
			multiErr = multierror.Append(multiErr, err)
		}
		if result != nil {
			results = append(results, *result)
		}
	}

	debug.Log(
		"phase", "analyzer.analyze-bundle",
		"status", "complete")
	return results, multiErr.ErrorOrNil()
}

func (a *Analyzer) analyze(ctx context.Context, bundleReader bundlereader.BundleReader, analyzerSpec v1.Analyzer) (*api.Result, error) {
	debug := level.Debug(log.With(a.Logger, "method", "Analyzer.analyze"))

	debug.Log(
		"phase", "analyzer.analyze",
		"spec", util.SpewJSON(analyzerSpec))

	data, err := a.registerVariables(analyzerSpec, bundleReader)
	if err != nil {
		return resultFromAnalysis(nil, err, analyzerSpec, data)
	}

	message, err := a.evalConditions(analyzerSpec, data)

	debug.Log(
		"phase", "analyzer.analyze",
		"status", "complete")

	return resultFromAnalysis(message, err, analyzerSpec, data)
}

func (a *Analyzer) registerVariables(analyzerSpec v1.Analyzer, bundleReader bundlereader.BundleReader) (map[string]interface{}, error) {
	debug := level.Debug(log.With(a.Logger, "method", "Analyzer.registerVariables"))

	data := map[string]interface{}{}
	for _, variable := range analyzerSpec.RegisterVariables {
		reg, err := variable.Register(bundleReader, data)
		for key, val := range reg {
			data[key] = val
		}
		debug.Log(
			"phase", "analyzer.analyze.register-variable",
			"variable", util.SpewJSON(variable),
			"register", util.SpewJSON(reg),
			"error", err)
		if err != nil {
			return data, errors.Wrapf(err, "register variable %s", variable.Name)
		}
	}
	return data, nil
}

func (a *Analyzer) evalConditions(analyzerSpec v1.Analyzer, data map[string]interface{}) (*message.Message, error) {
	debug := level.Debug(log.With(a.Logger, "method", "Analyzer.evalConditions"))

	if analyzerSpec.Precondition != nil {
		preconditionsOk, err := analyzerSpec.Precondition.Eval(data)
		debug.Log(
			"phase", "analyzer.analyze.preconditions-eval",
			"preconditions", util.SpewJSON(analyzerSpec.Precondition),
			"variables", util.SpewJSON(data),
			"ok", preconditionsOk,
			"error", err)
		if err != nil {
			return analyzerSpec.Messages.PreconditionError, errors.Wrap(err, "eval preconditions")
		} else if !preconditionsOk {
			return analyzerSpec.Messages.PreconditionFalse, nil
		}
	}

	conditionsOk, err := analyzerSpec.Condition.Eval(data)
	debug.Log(
		"phase", "analyzer.analyze.conditions-eval",
		"conditions", util.SpewJSON(analyzerSpec.Condition),
		"variables", util.SpewJSON(data),
		"ok", conditionsOk,
		"error", err)
	if err != nil {
		return analyzerSpec.Messages.ConditionError, errors.Wrap(err, "eval conditions")
	} else if !conditionsOk {
		return analyzerSpec.Messages.ConditionFalse, nil
	}

	return analyzerSpec.Messages.ConditionTrue, nil
}

func resultFromAnalysis(msg *message.Message, analysisErr error, analyzerSpec v1.Analyzer, data map[string]interface{}) (result *api.Result, err error) {
	if msg == nil && analysisErr == nil {
		return // off
	}

	result = &api.Result{
		Name:      analyzerSpec.Name,
		Variables: data,
	}

	var marshalledSpec []byte
	marshalledSpec, err = json.Marshal(api.Analyze{V1: []v1.Analyzer{analyzerSpec}})
	if err != nil {
		result.Severity = common.SeverityError
		err = errors.Wrap(err, "marshal spec")
		return
	}
	result.AnalyzerSpec = string(marshalledSpec)

	if analysisErr != nil {
		result.Error = analysisErr.Error()

		if msg == nil {
			result.Severity = common.SeverityError
			err = analysisErr
			return
		}
	}

	result.Message, err = msg.Render(data)
	if err != nil {
		result.Severity = common.SeverityError
		err = errors.Wrap(err, "render message")
		return
	}
	result.Severity = result.Message.Severity

	return
}
