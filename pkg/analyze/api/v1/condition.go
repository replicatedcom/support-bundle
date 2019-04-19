package v1

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/condition"
	"github.com/replicatedcom/support-bundle/pkg/meta"
)

type Condition struct {
	meta.Meta   `json:",inline" yaml:",inline" hcl:",inline"`
	VariableRef string `json:"variableRef,omitempty" yaml:"variableRef,omitempty" hcl:"variableRef,omitempty"` // Should we use a meta.Ref here?

	// built-in
	Empty         *condition.Empty         `json:"empty,omitempty" yaml:"empty,omitempty" hcl:"empty,omitempty"`
	EvalCondition *condition.EvalCondition `json:"eval,omitempty" yaml:"eval,omitempty" hcl:"eval,omitempty"`
	RegexpMatch   *condition.RegexpMatch   `json:"regexpMatch,omitempty" yaml:"regexpMatch,omitempty" hcl:"regexpMatch,omitempty"`
	StringCompare *condition.StringCompare `json:"stringCompare,omitempty" yaml:"stringCompare,omitempty" hcl:"stringCompare,omitempty"`

	// predicates
	And *AndPredicate `json:"and,omitempty" yaml:"and,omitempty" hcl:"and,omitempty"`
	Or  *OrPredicate  `json:"or,omitempty" yaml:"or,omitempty" hcl:"or,omitempty"`
	Not *NotPredicate `json:"not,omitempty" yaml:"not,omitempty" hcl:"not,omitempty"`
}

func (c *Condition) Eval(data map[string]interface{}) (bool, error) {
	e, tagName, ok := c.GetCondition()
	if !ok {
		return false, errors.New("condition not defined")
	}
	b, err := condition.Eval(e, c.VariableRef, data)
	if err != nil {
		return false, errors.Wrapf(err, "condition %q", tagName)
	}
	return b, nil
}

func (c *Condition) GetCondition() (condition.Interface, string, bool) {
	val := reflect.Indirect(reflect.ValueOf(c))
	for i := 0; i < val.NumField(); i++ {
		if e, ok := val.Field(i).Interface().(condition.Interface); ok && !reflect.ValueOf(e).IsNil() {
			return e, getTagName(val, i, "yaml"), true
		}
	}
	return nil, "", false
}
