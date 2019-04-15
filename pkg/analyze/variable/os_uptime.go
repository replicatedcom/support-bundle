package variable

import (
	"io"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/variable/distiller"
	bundlereader "github.com/replicatedcom/support-bundle/pkg/collect/bundle/reader"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

var (
	osUptimeRegexp = `(?m)([\d\.]+) (?:[\d\.]+)`

	_ = regexp.MustCompile(osUptimeRegexp)

	_ Interface = new(OsUptime)
)

type OsUptime struct {
}

func (v *OsUptime) MatchResults(bundleReader bundlereader.BundleReader) (results []collecttypes.Result) {
	for _, result := range bundleReader.GetIndex() {
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

func (v *OsUptime) ExtractValue(r io.Reader, result collecttypes.Result, data interface{}) (interface{}, error) {
	d := &distiller.RegexpCapture{
		Regexp: osUptimeRegexp,
		Index:  1,
	}
	b, err := distiller.Distill(d, r, false)
	if err != nil {
		return nil, errors.Wrap(err, "distill regexpCapture")
	}
	i, err := strconv.ParseFloat(b.(string), 64)
	if err != nil {
		return nil, errors.Wrap(err, "parse float")
	}
	return int(i), nil
}
