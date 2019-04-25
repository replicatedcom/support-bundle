package variable

import (
	"io"
	"regexp"

	"github.com/pkg/errors"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

var (
	_ Interface = new(FileMatch)
)

type FileMatch struct {
	PathRegexps []string `json:"pathRegexps" yaml:"pathRegexps" hcl:"pathRegexps"`

	Distiller `json:",inline" yaml:",inline" hcl:",inline"`
}

func (v *FileMatch) MatchResults(index []collecttypes.Result) ([]collecttypes.Result, error) {
	var rePath []*regexp.Regexp
	for _, path := range v.PathRegexps {
		re, err := regexp.Compile(path)
		if err != nil {
			return nil, errors.Wrapf(err, "compile regexp %q", path)
		}
		rePath = append(rePath, re)
	}
	var results []collecttypes.Result
	for _, result := range index {
		for _, re := range rePath {
			if re.MatchString(result.Path) {
				results = append(results, result)
			}
		}
	}
	return results, nil
}

func (v *FileMatch) DistillReader(r io.Reader, result collecttypes.Result) (interface{}, error) {
	return v.Distiller.Distill(r)
}

func (v *FileMatch) ExtractValue(distilled interface{}) (interface{}, error) {
	return distilled, nil
}
