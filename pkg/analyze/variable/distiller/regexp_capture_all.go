package distiller

import (
	"regexp"

	"github.com/pkg/errors"
)

var _ Interface = new(RegexpCaptureAll)

type RegexpCaptureAll struct {
	Regexp string `json:"regexp" yaml:"regexp" hcl:"regexp"`
	Index  int    `json:"index,omitempty" yaml:"index,omitempty" hcl:"index,omitempty"`
}

func (d *RegexpCaptureAll) Distill(input string) (interface{}, bool, error) {
	re, err := regexp.Compile(d.Regexp)
	if err != nil {
		return nil, false, errors.Wrap(err, "regexp compile")
	}
	matches := re.FindAllStringSubmatch(input, -1)
	if len(matches) == 0 || len(matches[0]) < d.Index+1 {
		return nil, false, nil
	}
	var all []string
	for _, submatches := range matches {
		all = append(all, submatches[d.Index])
	}
	return all, true, nil
}
