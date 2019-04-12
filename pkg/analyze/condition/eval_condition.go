package condition

import (
	"github.com/replicatedcom/support-bundle/pkg/templates"
)

var _ Interface = new(EvalCondition)

type EvalCondition struct {
	Value string `json:",inline" yaml:",inline" hcl:",inline"`
}

func (c *EvalCondition) Eval(ref interface{}, data map[string]interface{}) (bool, error) {
	return templates.Bool(c.Value, data)
}
