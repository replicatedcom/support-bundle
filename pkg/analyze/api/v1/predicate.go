package v1

import "github.com/pkg/errors"

type AndPredicate []Condition

type OrPredicate []Condition

type NotPredicate struct {
	Condition `json:",inline" yaml:",inline" hcl:",inline"`
}

func (c *AndPredicate) Eval(ref interface{}, data map[string]interface{}, err error) (bool, error) {
	for i, condition := range *c {
		b, err := condition.Eval(data)
		if err != nil || !b {
			return b, errors.Wrapf(err, "index %d", i)
		}
	}
	return true, nil
}

func (c *OrPredicate) Eval(ref interface{}, data map[string]interface{}, err error) (bool, error) {
	for i, condition := range *c {
		b, err := condition.Eval(data)
		if err != nil || b {
			return b, errors.Wrapf(err, "index %d", i)
		}
	}
	return false, nil
}

func (c *NotPredicate) Eval(ref interface{}, data map[string]interface{}, err error) (bool, error) {
	b, err := c.Condition.Eval(data)
	if err != nil {
		return b, err
	}
	return !b, nil
}
