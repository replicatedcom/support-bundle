package v1alpha1

import "github.com/replicatedcom/support-bundle/pkg/analyze/api/common"

type Condition struct {
	ConditionShared `json:",inline" yaml:",inline" hcl:",inline"`

	Eval *EvalCondition `json:"eval,omitempty" yaml:"eval,omitempty" hcl:"eval,omitempty"`
}

type ConditionShared struct {
	Severity common.Severity `json:"severity" yaml:"severity" hcl:"severity"`
	Message  string          `json:"message,omitempty" yaml:"message,omitempty" hcl:"message,omitempty"`
}

type EvalCondition struct {
	Operator   EvalOperator        `json:"operator,omitempty" yaml:"operator,omitempty" hcl:"operator,omitempty"`
	Statements []string            `json:"statements" yaml:"statements" hcl:"statements"`
	Variables  []map[string]string `json:"variables,omitempty" yaml:"variables,omitempty" hcl:"variables,omitempty"`
}

type EvalOperator string

const (
	AndOperator EvalOperator = "and"
	OrOperator  EvalOperator = "or"
)
