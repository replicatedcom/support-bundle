package variable

import (
	"io"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/variable/distiller"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

var (
	osReleaseRegexp     = `(?m)^ID="?([^"\n]+)"?`
	systemReleaseRegexp = `(CentOS|Red Hat Enterprise Linux Server|Amazon Linux AMI) release `

	_ = regexp.MustCompile(osReleaseRegexp)
	_ = regexp.MustCompile(systemReleaseRegexp)

	_ Interface = new(Os)
)

type Os struct {
}

func (v *Os) MatchResults(index []collecttypes.Result) ([]collecttypes.Result, error) {
	m := &CoreReadFilePath{
		Paths: []string{
			"/etc/os-release",
			"/usr/lib/os-release",
			"/etc/system-release",
		},
	}
	return m.MatchResults(index)
}

func (v *Os) DistillReader(r io.Reader, result collecttypes.Result) (interface{}, error) {
	parts := strings.Split(result.Spec.CoreReadFile.Filepath, "/")
	switch parts[len(parts)-1] {

	case "os-release":
		d := &distiller.RegexpCapture{
			Regexp: osReleaseRegexp,
			Index:  1,
		}
		str, err := distiller.Distill(d, r, false)
		return str, errors.Wrap(err, "distill regexpCapture")

	case "system-release":
		// special cases for:
		// centos 6
		// rhel 6
		// amzn 2014.03
		d := &distiller.RegexpCapture{
			Regexp: systemReleaseRegexp,
			Index:  1,
		}
		i, err := distiller.Distill(d, r, false)
		if err != nil {
			return i, errors.Wrap(err, "distill regexpCapture")
		}
		str, _ := i.(string)
		switch strings.ToLower(str) {
		case "centos":
			return "centos", nil
		case "red hat enterprise linux server":
			return "rhel", nil
		case "amazon linux ami":
			return "amzn", nil
		}
	}
	return nil, nil
}

func (v *Os) ExtractValue(distilled interface{}) (interface{}, error) {
	return distilled, nil
}
