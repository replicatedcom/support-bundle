package condition

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
)

var _ Interface = new(RegexpMatch)

type RegexpMatch struct {
	Regexp string `json:"regexp" yaml:"regexp" hcl:"regexp"`
}

func (c *RegexpMatch) Eval(ref interface{}, data map[string]interface{}, err error) (bool, error) {
	if err != nil {
		return false, err
	}
	var str string
	if ref != nil {
		str = fmt.Sprintf("%v", ref)
	}
	re, err := regexp.Compile(c.Regexp)
	if err != nil {
		return false, errors.Wrap(err, "compile regexp")
	}
	return re.MatchString(str), nil
}
