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
	v := toString(ref)
	switch {
	case c.Eq != nil:
		return v == toString(c.Eq), nil
	case c.Lt != nil:
		return v < toString(c.Lt), nil
	case c.Lte != nil:
		return v <= toString(c.Lte), nil
	case c.Gt != nil:
		return v > toString(c.Gt), nil
	case c.Gte != nil:
		return v >= toString(c.Gte), nil
	case c.In != nil:
		for _, i := range c.In {
			if v == toString(i) {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, errors.New("comparison function expected")
	}
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%s", v)
}
