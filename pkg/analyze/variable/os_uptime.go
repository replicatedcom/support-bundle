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

func (v *OsUptime) MatchResults(index []collecttypes.Result) ([]collecttypes.Result, error) {
	return matchAny(
		matchCollector(&collecttypes.CoreUptimeOptions{}),
		(&CoreReadFilePath{
			Paths: []string{
				"/proc/uptime",
			},
		}).MatchResults,
	)(index)
}

func (v *OsUptime) DistillReader(r io.Reader, result collecttypes.Result) (interface{}, error) {
	d := &distiller.RegexpCapture{
		Regexp: osUptimeRegexp,
		Index:  1,
	}
	b, err := distiller.Distill(d, r, false)
	return b, errors.Wrap(err, "distill regexpCapture")
}

func (v *OsUptime) ExtractValue(distilled interface{}) (interface{}, error) {
	if distilled == nil {
		return nil, nil
	}
	f, err := strconv.ParseFloat(distilled.(string), 64)
	if err != nil {
		return nil, errors.Wrap(err, "parse float")
	}
	return f, nil
}
