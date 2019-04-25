package condition

import (
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/templates"
	"github.com/replicatedcom/support-bundle/pkg/util"
)

var _ Interface = new(EvalCondition)

type EvalCondition string

func (c *EvalCondition) Eval(ref interface{}, data map[string]interface{}, err error) (bool, error) {
	if err != nil && errors.Cause(err) != ErrNotFound {
		return false, err
	}
	copy := util.CopyMap(data)
	copy["Ref"] = ref
	return templates.Bool(string(*c), copy)
}
