package analyzer

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api/v1alpha1"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
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

func (a *Analyzer) AnalyzeBundle(ctx context.Context, spec api.Analyze, collectSpec api.Collect, collectBundle string) ([]api.Result, error) {
	debug := level.Debug(log.With(a.Logger, "method", "Analyzer.AnalyzeBundle"))

	debug.Log(
		"phase", "analyzer.analyze-bundle")

	r := BundleReader{Fs: a.Fs}
	index, err := r.GetResultsFromIndex(collectBundle, "index.json")
	debug.Log(
		"phase", "analyzer.get-bundle-index",
		"index", spew.Sdump(index),
		"error", err)
	if err != nil {
		return nil, errors.Wrap(err, "get index.json from bundle")
	}

	var results []api.Result
	for _, analyzerSpec := range spec.V1Alpha1 {
		result, err := a.analyze(ctx, analyzerSpec, index, collectBundle)
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

func (a *Analyzer) analyze(ctx context.Context, analyzerSpec v1alpha1.AnalyzerSpec, index []collecttypes.Result, archivePath string) (api.Result, error) {
	debug := level.Debug(log.With(a.Logger, "method", "Analyzer.analyze"))

	debug.Log(
		"phase", "analyzer.analyze",
		"spec", spew.Sdump(analyzerSpec))

	var result api.Result
	result.AnalyzerSpec = api.Analyze{V1Alpha1: []v1alpha1.AnalyzerSpec{analyzerSpec}}

	analyzer := analyzerSpec.GetAnalyzer()
	debug.Log(
		"phase", "analyzer.analyze.get-analyzer",
		"analyzer", spew.Sdump(analyzer))
	if analyzer == nil {
		return result, errors.New("analyzer empty") // TODO: typed error
	}

	rawSpec, err := analyzer.GetSpec()
	debug.Log(
		"phase", "analyzer.analyze.get-spec",
		"rawSpec", spew.Sdump(rawSpec),
		"error", err)
	if err != nil {
		return result, errors.Wrap(err, "get analyzer spec")
	}
	rawSpec.CollectRefs = analyzerSpec.CollectRefs
	rawSpec.CollectRefs[0].Ref = "Ref"
	rawSpec.Meta = analyzerSpec.Meta
	if analyzerSpec.Message != "" {
		rawSpec.Message = analyzerSpec.Message
	}
	result.Requirement = rawSpec.Message

	results := matchCollectResults(rawSpec.CollectRefs, index)
	debug.Log(
		"phase", "analyzer.analyze.match-collect-results",
		"results", spew.Sdump(results))

	r := &BundleReader{Fs: a.Fs}
	data, err := collectRefData(r, rawSpec.CollectRefs, results, archivePath)
	debug.Log(
		"phase", "analyzer.analyze.build-template-data",
		"data", spew.Sdump(data),
		"error", err)
	if err != nil {
		return result, errors.Wrap(err, "collect ref data")
	}

	// TODO: validate

	for _, condition := range rawSpec.Raw.Conditions {
		vars, err := BuildConditionVariables(*condition.Eval, data) // TODO: will eval be the only condition?
		if err != nil {
			return result, errors.Wrap(err, "build variables")
		}
		result.Vars = append(result.Vars, vars)

		ok, err := EvalCondition(*condition.Eval, vars)
		debug.Log(
			"phase", "analyzer.analyze.condition",
			"condition", spew.Sdump(condition),
			"met", ok,
			"err", err)
		if err != nil {
			return result, errors.Wrap(err, "check condition")
		}

		if ok {
			result.Severity = condition.Severity

			message, err := templates.String(condition.Message, vars)
			if err != nil {
				return result, errors.Wrap(err, "execute message template")
			}
			result.Message = message // TODO: build
			break
		}
	}

	debug.Log(
		"phase", "analyzer.analyze",
		"status", "complete")

	// TODO: severity?
	return result, nil
}

func collectRefData(bundleReader *BundleReader, refs []meta.Ref, results map[string]collecttypes.Result, archivePath string) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	for _, ref := range refs {
		result, ok := results[ref.Ref]
		if !ok {
			continue
		}
		r, err := bundleReader.FileReaderFromArchive(archivePath, strings.TrimLeft(result.Path, "/"))
		if err != nil {
			return data, errors.Wrapf(err, "ref %s", ref.Ref)
		}
		b, err := ioutil.ReadAll(r)
		r.Close()
		if err != nil {
			return data, errors.Wrapf(err, "ref %s", ref.Ref)
		}
		data[ref.Ref] = string(b)
	}
	return data, nil
}

func matchCollectResults(refs []meta.Ref, index []collecttypes.Result) map[string]collecttypes.Result {
	matches := map[string]collecttypes.Result{}
	for _, ref := range refs {
		for _, result := range index {
			if meta.RefMatches(ref, result.Spec.Shared().Meta) {
				matches[ref.Ref] = result
			}
		}
	}
	return matches
}
