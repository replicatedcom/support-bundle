package condition

import (
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/templates"
)

var _ Interface = new(EvalCondition)

type EvalCondition string

func (c *EvalCondition) Eval(ref interface{}, data map[string]interface{}, err error) (bool, error) {
	if err != nil && errors.Cause(err) != ErrNotFound {
		return false, err
	}
	copy := map[string]interface{}{}
	for key, value := range data {
		copy[key] = value
	}
	copy["Ref"] = ref
	return templates.Bool(string(*c), copy)
}
