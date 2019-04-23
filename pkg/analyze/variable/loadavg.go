package variable

import (
	"io"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/variable/distiller"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

var (
	loadavgRegexp = `([^ ]+) *([^ ]+) *([^ ]+) *([^ ]+) *([^ ]+)`

	_ = regexp.MustCompile(loadavgRegexp)

	_, _, _ Interface = new(Loadavg1), new(Loadavg5), new(Loadavg15)
)

type Loadavg struct {
}

type Loadavg1 struct {
}

type Loadavg5 struct {
}

type Loadavg15 struct {
}

func (v *Loadavg) MatchResults(index []collecttypes.Result) ([]collecttypes.Result, error) {
	return matchResultsLoadavg(index)
}

func (v *Loadavg1) MatchResults(index []collecttypes.Result) ([]collecttypes.Result, error) {
	return matchResultsLoadavg(index)
}

func (v *Loadavg5) MatchResults(index []collecttypes.Result) ([]collecttypes.Result, error) {
	return matchResultsLoadavg(index)
}

func (v *Loadavg15) MatchResults(index []collecttypes.Result) ([]collecttypes.Result, error) {
	return matchResultsLoadavg(index)
}

func matchResultsLoadavg(index []collecttypes.Result) ([]collecttypes.Result, error) {
	return matchAny(
		matchCollector(&collecttypes.CoreLoadavgOptions{}),
		(&CoreReadFilePath{
			Paths: []string{
				"/proc/loadavg",
			},
		}).MatchResults,
	)(index)
}

func (v *Loadavg) DistillReader(r io.Reader, result collecttypes.Result) (interface{}, error) {
	d := &distiller.Identity{}
	b, _, err := distiller.Distill(d, r, false)
	return b, errors.Wrap(err, "distill identity")
}

func (v *Loadavg1) DistillReader(r io.Reader, result collecttypes.Result) (interface{}, error) {
	return distillReaderLoadavg(r, result, 1)
}

func (v *Loadavg5) DistillReader(r io.Reader, result collecttypes.Result) (interface{}, error) {
	return distillReaderLoadavg(r, result, 2)
}

func (v *Loadavg15) DistillReader(r io.Reader, result collecttypes.Result) (interface{}, error) {
	return distillReaderLoadavg(r, result, 3)
}

func distillReaderLoadavg(r io.Reader, result collecttypes.Result, index int) (interface{}, error) {
	d := &distiller.RegexpCapture{
		Regexp: loadavgRegexp,
		Index:  index,
	}
	b, _, err := distiller.Distill(d, r, false)
	return b, errors.Wrap(err, "distill regexpCapture")
}

func (v *Loadavg) ExtractValue(distilled interface{}, data map[string]interface{}) (interface{}, error) {
	return distilled, nil
}

func (v *Loadavg1) ExtractValue(distilled interface{}, data map[string]interface{}) (interface{}, error) {
	return extractValueLoadavg(distilled, data)
}

func (v *Loadavg5) ExtractValue(distilled interface{}, data map[string]interface{}) (interface{}, error) {
	return extractValueLoadavg(distilled, data)
}

func (v *Loadavg15) ExtractValue(distilled interface{}, data map[string]interface{}) (interface{}, error) {
	return extractValueLoadavg(distilled, data)
}

func extractValueLoadavg(distilled interface{}, data map[string]interface{}) (interface{}, error) {
	if distilled == nil {
		return nil, nil
	}
	i, err := strconv.ParseFloat(distilled.(string), 64)
	return i, errors.Wrap(err, "parse float")
}
