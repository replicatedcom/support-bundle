package condition

import (
	"errors"
	"fmt"
)

var _ Interface = new(StringCompare)

type StringCompare struct {
	Compare `json:",inline" yaml:",inline" hcl:",inline"`
}

func (c *StringCompare) Eval(ref interface{}, data map[string]interface{}) (bool, error) {
	var str string
	if ref != nil {
		str = fmt.Sprintf("%v", ref)
	}
	switch {
	case c.Eq != "":
		return str == c.Eq, nil
	case c.Lt != "":
		return str < c.Lt, nil
	case c.Lte != "":
		return str <= c.Lte, nil
	case c.Gt != "":
		return str > c.Gt, nil
	case c.Gte != "":
		return str >= c.Gte, nil
	default:
		return false, errors.New("comparison function expected")
	}
}
