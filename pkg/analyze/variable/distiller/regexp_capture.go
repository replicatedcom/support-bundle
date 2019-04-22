package distiller

import (
	"regexp"

	"github.com/pkg/errors"
)

var _ Interface = new(RegexpCapture)

type RegexpCapture struct {
	Regexp string `json:"regexp" yaml:"regexp" hcl:"regexp"`
	Index  int    `json:"index,omitempty" yaml:"index,omitempty" hcl:"index,omitempty"`
}

func (d *RegexpCapture) Distill(input string) (interface{}, bool, error) {
	re, err := regexp.Compile(d.Regexp)
	if err != nil {
		return nil, false, errors.Wrap(err, "regexp compile")
	}
	matches := re.FindStringSubmatch(input)
	if len(matches) < d.Index+1 {
		return nil, false, nil
	}
	return matches[d.Index], true, nil
}
