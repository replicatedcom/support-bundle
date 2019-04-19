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
	osUptimeRegexp = `(?m)([\d\.]+) (?:[\d\.]+)`

	_ = regexp.MustCompile(osUptimeRegexp)

	_ Interface = new(OsUptime)
)

type OsUptime struct {
}

func (v *OsUptime) MatchResults(index []collecttypes.Result) (results []collecttypes.Result) {
	for _, result := range index {
		if matchAny(
			result,
			matcherCoreReadFileFilepath("/proc/uptime"),
			matcherCollector(&collecttypes.CoreUptimeOptions{}),
		) {
			results = append(results, result)
		}
	}
	return
}

func (v *OsUptime) DistillReader(r io.Reader, result collecttypes.Result) (string, error) {
	d := &distiller.RegexpCapture{
		Regexp: osUptimeRegexp,
		Index:  1,
	}
	b, err := distiller.Distill(d, r, false)
	return string(b), errors.Wrap(err, "distill regexpCapture")
}

func (v *OsUptime) ExtractValue(distilled interface{}, data interface{}) (interface{}, error) {
	if distilled == nil {
		return nil, nil
	}
	i, err := strconv.ParseFloat(distilled.(string), 64)
	if err != nil {
		return nil, errors.Wrap(err, "parse float")
	}
	return int(i), nil
}
