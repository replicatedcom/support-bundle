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
	memoryUsageTotalRegexp     = `(?m)^MemTotal: +(\d+) kB`
	memoryUsageAvailableRegexp = `(?m)^MemAvailable: +(\d+) kB`

	_ = regexp.MustCompile(memoryUsageTotalRegexp)
	_ = regexp.MustCompile(memoryUsageAvailableRegexp)

	_ Interface = new(MemoryUsageTotal)
	_ Interface = new(MemoryUsageAvailable)
)

type MemoryUsageTotal struct {
}

type MemoryUsageAvailable struct {
}

func (v *MemoryUsageTotal) MatchResults(index []collecttypes.Result) ([]collecttypes.Result, error) {
	return matchResultsMemoryUsage(index)
}

func (v *MemoryUsageAvailable) MatchResults(index []collecttypes.Result) ([]collecttypes.Result, error) {
	return matchResultsMemoryUsage(index)
}

func matchResultsMemoryUsage(index []collecttypes.Result) ([]collecttypes.Result, error) {
	return matchAny(
		matchCollector(&collecttypes.CoreLoadavgOptions{}),
		(&CoreReadFilePath{
			Paths: []string{
				"/proc/meminfo",
			},
		}).MatchResults,
	)(index)
}

func (v *MemoryUsageTotal) DistillReader(r io.Reader, result collecttypes.Result) (interface{}, error) {
	return distillReaderMemoryUsage(r, result, memoryUsageTotalRegexp)
}

func (v *MemoryUsageAvailable) DistillReader(r io.Reader, result collecttypes.Result) (interface{}, error) {
	return distillReaderMemoryUsage(r, result, memoryUsageAvailableRegexp)
}

func distillReaderMemoryUsage(r io.Reader, result collecttypes.Result, re string) (interface{}, error) {
	d := &distiller.RegexpCapture{
		Regexp: re,
		Index:  1,
	}
	b, err := distiller.Distill(d, r, false)
	return b, errors.Wrap(err, "distill regexpCapture")
}

func (v *MemoryUsageTotal) ExtractValue(distilled interface{}) (interface{}, error) {
	return extractValueMemoryUsage(distilled)
}

func (v *MemoryUsageAvailable) ExtractValue(distilled interface{}) (interface{}, error) {
	return extractValueMemoryUsage(distilled)
}

func extractValueMemoryUsage(distilled interface{}) (interface{}, error) {
	if distilled == nil {
		return nil, nil
	}
	i, err := strconv.ParseFloat(distilled.(string), 64)
	return i * 1024, errors.Wrap(err, "parse float")
}
