package analyzer

import (
	"fmt"
	"io"
	"io/ioutil"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/replicatedcom/support-bundle/pkg/analyze/variable"

	"github.com/pkg/errors"
	v1 "github.com/replicatedcom/support-bundle/pkg/analyze/api/v1"
	bundlereader "github.com/replicatedcom/support-bundle/pkg/collect/bundle/reader"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
	"github.com/replicatedcom/support-bundle/pkg/util"
)

func (a *Analyzer) registerVariables(variables []v1.Variable, bundleReader bundlereader.BundleReader) (map[string]interface{}, error) {
	debug := level.Debug(log.With(a.Logger, "method", "Analyzer.registerVariables"))

	variableNamesToDistilled, err := a.distillBundle(variables, bundleReader)
	debug.Log(
		"phase", "analyzer.distill-bundle",
		"distilled", util.SpewJSON(variableNamesToDistilled),
		"error", err)
	if err != nil {
		return nil, errors.Wrap(err, "distill bundle")
	}

	data, err := a.extractValues(variables, variableNamesToDistilled)
	debug.Log(
		"phase", "analyzer.extract-values",
		"values", util.SpewJSON(data),
		"error", err)
	return data, errors.Wrap(err, "extract values")
}

func (a *Analyzer) extractValues(variables []v1.Variable, variableNamesToDistilled map[string][]string) (map[string]interface{}, error) {
	data := map[string]interface{}{}

	for _, v := range variables {
		name := v.GetName()
		i, _, _ := v.GetVariable()
		if distilled, ok := variableNamesToDistilled[name]; ok {
			for _, d := range distilled {
				value, err := i.ExtractValue(d, data)
				if err != nil {
					return data, errors.Wrapf(err, "variable %s", name)
				} else if value != nil {
					// use the first value we find that is not empty
					data[name] = value
					break
				}
			}
		} else {
			value, err := i.ExtractValue(nil, data)
			if err != nil {
				return data, errors.Wrapf(err, "variable %s", name)
			} else if value != nil {
				// use the first value we find that is not empty
				data[name] = value
				break
			}
		}
	}

	return data, nil
}

func (a *Analyzer) distillBundle(variables []v1.Variable, bundleReader bundlereader.BundleReader) (map[string][]string, error) {
	variablesMap := map[string]variable.Interface{}
	resultsToVariables := map[collecttypes.Result][]string{}

	index := bundleReader.GetIndex()

	for i, v := range variables {
		name := v.GetName()
		if name == "" {
			return nil, fmt.Errorf("variable at index %d name empty", i)
		}
		variable, _, ok := v.GetVariable()
		if !ok {
			return nil, fmt.Errorf("no variable defined for %s", v.GetName())
		}
		variablesMap[name] = variable
		for _, result := range variable.MatchResults(index) {
			resultsToVariables[result] = append(resultsToVariables[result], name)
		}
	}

	scanner, err := bundleReader.NewScanner()
	if err != nil {
		return nil, errors.Wrap(err, "new scanner")
	}
	defer scanner.Close()

	values := map[string][]string{}
	for {
		f, err := scanner.Next()
		if err == io.EOF {
			return values, nil
		} else if err != nil {
			return values, errors.Wrap(err, "scanner next")
		}

		for result, variables := range resultsToVariables {
			if !bundlereader.PathsMatch(f.Name, result.Path) {
				continue
			}
			vv, err := a.distillReader(f.Reader, result, variables, variablesMap)
			if err != nil {
				return values, errors.Wrapf(err, "distill result %q", result.Path)
			}
			for name, value := range vv {
				values[name] = append(values[name], value)
			}
		}
	}
}

func (a *Analyzer) distillReader(reader io.Reader, result collecttypes.Result, variableNames []string, variablesMap map[string]variable.Interface) (map[string]string, error) {
	if len(variableNames) == 0 {
		return nil, nil
	}

	values := map[string]string{}
	var valueMu sync.Mutex
	var errGroup util.MultiErrorGroup
	var pws []*io.PipeWriter

	for _, variableName := range variableNames {
		pr, pw := io.Pipe()
		pws = append(pws, pw)
		reader = io.TeeReader(reader, pw)
		func(variableName string) { // new scope
			errGroup.Go(func() error {
				value, err := variablesMap[variableName].DistillReader(pr, result)
				pr.CloseWithError(err)
				valueMu.Lock()
				values[variableName] = value
				valueMu.Unlock()
				return errors.Wrapf(err, "variable %q", variableName)
			})
		}(variableName)
	}

	// Should this error be returned?
	_, err := io.Copy(ioutil.Discard, reader)
	for _, pw := range pws {
		pw.CloseWithError(err)
	}

	return values, errGroup.Wait()
}
