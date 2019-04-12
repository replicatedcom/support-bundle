package condition

import (
	"github.com/replicatedcom/support-bundle/pkg/templates"
)

var _ Interface = new(EvalCondition)

type EvalCondition string

func (c *EvalCondition) Eval(ref interface{}, data map[string]interface{}) (bool, error) {
	return templates.Bool(string(*c), data)
}
