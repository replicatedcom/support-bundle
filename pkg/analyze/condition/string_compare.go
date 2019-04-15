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
	case c.Eq != nil:
		return str == fmt.Sprintf("%v", c.Eq), nil
	case c.Lt != nil:
		return str < fmt.Sprintf("%v", c.Lt), nil
	case c.Lte != nil:
		return str <= fmt.Sprintf("%v", c.Lte), nil
	case c.Gt != nil:
		return str > fmt.Sprintf("%v", c.Gt), nil
	case c.Gte != nil:
		return str >= fmt.Sprintf("%v", c.Gte), nil
	case c.In != nil:
		for _, i := range c.In {
			if str == fmt.Sprintf("%v", i) {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, errors.New("comparison function expected")
	}
}
