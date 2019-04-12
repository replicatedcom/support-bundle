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

func (d *RegexpCapture) Distill(input string) (interface{}, error) {
	re, err := regexp.Compile(d.Regexp)
	if err != nil {
		return nil, errors.Wrap(err, "regexp compile")
	}
	matches := re.FindStringSubmatch(input)
	if len(matches) < d.Index+1 {
		return nil, nil
	}
	return matches[d.Index], nil
}
