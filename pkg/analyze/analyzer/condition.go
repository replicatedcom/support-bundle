package analyzer

import (
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api/v1"
	"github.com/replicatedcom/support-bundle/pkg/templates"
)

func EvalCondition(condition v1.Condition, vars map[string]interface{}) (bool, error) {
	// TODO: will eval be the only condition?
	if condition.Eval == nil {
		// TODO: what should we return here?
		return false, nil
	}

	eval := true
	for _, statement := range condition.Eval.Statements {
		b, err := templates.Bool(statement, vars)
		if err != nil {
			return false, errors.Wrapf(err, "execute statement %q", statement)
		}
		if condition.Eval.Operator == v1.OrOperator {
			if b == true {
				// short circuit
				return true, nil
			}
		} else { // default to and operator
			eval = eval && b
		}
	}
	return eval, nil
}

func BuildConditionVariables(condition v1.Condition, data map[string]interface{}) (map[string]interface{}, error) {
	copy := map[string]interface{}{}
	for k := range data {
		copy[k] = data[k]
	}

	// TODO: will eval be the only condition?
	if condition.Eval == nil {
		// TODO: what should we return here?
		return copy, nil
	}

	vars := map[string]interface{}{}
	for _, step := range condition.Eval.Variables {
		for key, tmpl := range step {
			built, err := templates.Execute(tmpl, copy)
			if err != nil {
				return vars, errors.Wrapf(err, "var %s", key)
			}
			vars[key] = built
			copy[key] = built // vars in vars
		}
	}
	return vars, nil
}
