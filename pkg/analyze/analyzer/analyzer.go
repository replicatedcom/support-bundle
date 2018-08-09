package analyzer

import (
	"context"
	"io/ioutil"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api/v1"
	bundlereader "github.com/replicatedcom/support-bundle/pkg/collect/bundle/reader"
	"github.com/replicatedcom/support-bundle/pkg/meta"
	"github.com/replicatedcom/support-bundle/pkg/spew"
	"github.com/replicatedcom/support-bundle/pkg/templates"
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
	for _, analyzerSpec := range spec.V1 {
		result, err := a.analyze(ctx, bundleReader, analyzerSpec)
		if err != nil {
			return results, errors.Wrapf(err, "analyze spec")
		}
		results = append(results, result)
	}

	debug.Log(
		"phase", "analyzer.analyze-bundle",
		"status", "complete")
	return results, nil
}

func (a *Analyzer) analyze(ctx context.Context, bundleReader bundlereader.BundleReader, analyzerSpec v1.AnalyzerSpec) (api.Result, error) {
	debug := level.Debug(log.With(a.Logger, "method", "Analyzer.analyze"))

	debug.Log(
		"phase", "analyzer.analyze",
		"spec", spew.Sdump(analyzerSpec))

	var result api.Result
	result.AnalyzerSpec = api.Analyze{V1: []v1.AnalyzerSpec{analyzerSpec}}

	requirement := analyzerSpec.GetRequirement()
	debug.Log(
		"phase", "analyzer.analyze.get-requirement",
		"analyzer", spew.Sdump(requirement))
	if requirement == nil {
		return result, errors.New("analyzer empty") // TODO: typed error
	}

	// TODO: analyzer.Validate(analyzerSpec)

	rawSpec, err := getRawSpec(analyzerSpec, requirement)
	debug.Log(
		"phase", "analyzer.analyze.get-spec",
		"rawSpec", spew.Sdump(rawSpec),
		"error", err)
	if err != nil {
		return result, errors.Wrap(err, "get raw spec")
	}

	result.Requirement = rawSpec.Message

	data, err := collectRefData(bundleReader, rawSpec.CollectRefs)
	debug.Log(
		"phase", "analyzer.analyze.build-template-data",
		"data", spew.Sdump(data),
		"error", err)
	if err != nil {
		return result, errors.Wrap(err, "collect ref data")
	}

	for _, condition := range rawSpec.Raw.Conditions {
		debug.Log(
			"phase", "analyzer.analyze.condition",
			"condition", spew.Sdump(condition))

		vars, err := BuildConditionVariables(condition, data)
		debug.Log(
			"phase", "analyzer.analyze.condition.build",
			"err", err)
		if err != nil {
			return result, errors.Wrap(err, "build variables")
		}
		result.Vars = append(result.Vars, vars)

		ok, err := EvalCondition(condition, vars)
		debug.Log(
			"phase", "analyzer.analyze.condition.eval",
			"met", ok,
			"err", err)
		if err != nil {
			return result, errors.Wrap(err, "check condition")
		}

		if ok {
			// severity override
			if rawSpec.Severity != "" {
				result.Severity = rawSpec.Severity
			} else {
				result.Severity = condition.Severity
			}

			message, err := templates.String(condition.Message, vars)
			if err != nil {
				return result, errors.Wrap(err, "execute message template")
			}
			result.Message = message
			break
		}
	}

	debug.Log(
		"phase", "analyzer.analyze",
		"status", "complete")

	return result, nil
}

func getRawSpec(analyzerSpec v1.AnalyzerSpec, requirement v1.Requirement) (v1.RawSpec, error) {
	rawSpec, err := requirement.GetRawSpec()
	if err != nil {
		return rawSpec, err
	}
	rawSpec.CollectRefs = analyzerSpec.CollectRefs
	rawSpec.CollectRefs[0].Ref = "_Ref"
	rawSpec.Meta = analyzerSpec.Meta
	if analyzerSpec.Message != "" {
		rawSpec.Message = analyzerSpec.Message
	}
	if analyzerSpec.Severity != "" {
		rawSpec.Severity = analyzerSpec.Severity
	}
	return rawSpec, nil
}

func collectRefData(bundleReader bundlereader.BundleReader, refs []meta.Ref) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	for _, ref := range refs {
		r, err := bundleReader.ReaderFromRef(ref)
		if err != nil {
			return data, errors.Wrapf(err, "ref %s", ref.Ref)
		}
		if r != nil {
			b, err := ioutil.ReadAll(r)
			r.Close()
			if err != nil {
				return data, errors.Wrapf(err, "ref %s", ref.Ref)
			}
			data[ref.Ref] = string(b)
		}
	}
	return data, nil
}
