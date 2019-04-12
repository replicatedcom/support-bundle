package v1

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/condition"
)

type Condition struct {
	// built-in
	EvalCondition *condition.EvalCondition `json:"eval,omitempty" yaml:"eval,omitempty" hcl:"eval,omitempty"`
	StringCompare *condition.StringCompare `json:"stringCompare,omitempty" yaml:"stringCompare,omitempty" hcl:"stringCompare,omitempty"`

	// predicates
	And *AndPredicate `json:"and,omitempty" yaml:"and,omitempty" hcl:"and,omitempty"`
	Or  *OrPredicate  `json:"or,omitempty" yaml:"or,omitempty" hcl:"or,omitempty"`
	Not *NotPredicate `json:"not,omitempty" yaml:"not,omitempty" hcl:"not,omitempty"`

	// common
	Ref string `json:"ref,omitempty" yaml:"ref,omitempty" hcl:"ref,omitempty"`
}

func (c *Condition) Eval(data map[string]interface{}) (bool, error) {
	val := reflect.Indirect(reflect.ValueOf(c))
	for i := 0; i < val.NumField(); i++ {
		if e, ok := val.Field(i).Interface().(condition.Interface); ok && !reflect.ValueOf(e).IsNil() {
			b, err := condition.Eval(e, c.Ref, data)
			if err != nil {
				return false, errors.Wrapf(err, "condition %q", getTagName(val, i, "yaml"))
			}
			return b, nil
		}
	}
	return false, errors.New("no condition defined")
}

func getTagName(v reflect.Value, i int, key string) string {
	return strings.SplitN(v.Type().Field(i).Tag.Get(key), ",", 2)[0]
}
