package variable

import (
	"io"
	"regexp"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/variable/distiller"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

var (
	cpuCoresRegexp = `processor\s+:\s+(\d+)`

	_ = regexp.MustCompile(cpuCoresRegexp)

	_ Interface = new(CPUCores)
)

type CPUCores struct {
}

func (v *CPUCores) MatchResults(index []collecttypes.Result) (results []collecttypes.Result) {
	m := &CoreReadFilePath{
		Paths: []string{
			"/proc/cpuinfo",
		},
	}
	return m.MatchResults(index)
}

func (v *CPUCores) DistillReader(r io.Reader, result collecttypes.Result) (interface{}, error) {
	d := &distiller.RegexpCaptureAll{
		Regexp: cpuCoresRegexp,
		Index:  1,
	}
	b, _, err := distiller.Distill(d, r, false)
	return b, errors.Wrap(err, "distill regexpCaptureAll")
}

func (v *CPUCores) ExtractValue(distilled interface{}, data interface{}) (interface{}, error) {
	switch value := distilled.(type) {
	case []string:
		numproc := len(value)
		if numproc < 1 {
			return nil, nil
		}
		return numproc, nil
	}
	return nil, nil
}
