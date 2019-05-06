package v1

import (
	"github.com/replicatedcom/support-bundle/pkg/analyze/insight"
	"github.com/replicatedcom/support-bundle/pkg/meta"
)

type Analyzer struct {
	meta.Meta `json:",inline" yaml:",inline" hcl:",inline"`

	Insight *insight.Insight `json:"insight" yaml:"insight" hcl:"insight"`

	RegisterVariables  []Variable          `json:"registerVariables" yaml:"registerVariables" hcl:"registerVariables"`
	EvaluateConditions []EvaluateCondition `json:"evaluateConditions" yaml:"evaluateConditions" hcl:"evaluateConditions"`
}
