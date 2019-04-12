package variable

import (
	"io"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/variable/distiller"
	bundlereader "github.com/replicatedcom/support-bundle/pkg/collect/bundle/reader"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

var (
	osReleaseRegexp     = `(?m)^ID="?([^"\n]+)"?`
	systemReleaseRegexp = `(CentOS|Red Hat Enterprise Linux Server|Amazon Linux AMI) release `

	_ Interface = new(Os)

	_ = regexp.MustCompile(osReleaseRegexp)
	_ = regexp.MustCompile(systemReleaseRegexp)
)

type Os struct {
}

func (v *Os) MatchResults(bundleReader bundlereader.BundleReader) (results []collecttypes.Result) {
	for _, result := range bundleReader.GetIndex() {
		switch {
		case result.Spec.CoreReadFile != nil:
			if result.Spec.CoreReadFile.Filepath == "/etc/os-release" ||
				result.Spec.CoreReadFile.Filepath == "/usr/lib/os-release" ||
				result.Spec.CoreReadFile.Filepath == "/etc/system-release" {

				if result.Size > 0 {
					results = append(results, result)
				}
			}
		}
	}
	return
}

func (v *Os) ExtractValue(r io.Reader, result collecttypes.Result, data interface{}) (interface{}, error) {
	parts := strings.Split(result.Spec.CoreReadFile.Filepath, "/")
	switch parts[len(parts)-1] {

	case "os-release":
		d := &distiller.RegexpCapture{
			Regexp: osReleaseRegexp,
			Index:  1,
		}
		b, err := distiller.Distill(d, r, false)
		return b, errors.Wrap(err, "distill regexpCapture")

	case "system-release":
		// special cases for:
		// centos 6
		// rhel 6
		// amzn 2014.03
		d := &distiller.RegexpCapture{
			Regexp: systemReleaseRegexp,
			Index:  1,
		}
		b, err := distiller.Distill(d, r, false)
		if err != nil {
			return nil, errors.Wrap(err, "distill regexpCapture")
		} else if b != nil {
			switch strings.ToLower(b.(string)) {
			case "centos":
				return "centos", nil
			case "red hat enterprise linux server":
				return "rhel", nil
			case "amazon linux ami":
				return "amzn", nil
			}
		}
	}

	return nil, nil
}
