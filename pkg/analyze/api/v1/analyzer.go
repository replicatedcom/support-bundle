package v1

import "github.com/replicatedcom/support-bundle/pkg/analyze/message"

type Analyzer struct {
	Name              string     `json:"name" yaml:"name" hcl:"name"`
	RegisterVariables []Variable `json:"registerVariables" yaml:"registerVariables" hcl:"registerVariables"`
	Precondition      *Condition `json:"precondition,omitempty" yaml:"precondition,omitempty" hcl:"precondition,omitempty"`
	Condition         Condition  `json:"condition" yaml:"condition" hcl:"condition"`
	Messages          Messages   `json:"messages" yaml:"messages" hcl:"messages"`
}

type Messages struct {
	PreconditionError *message.Message `json:"preconditionError,omitempty" yaml:"preconditionError,omitempty" hcl:"preconditionError,omitempty"`
	PreconditionFalse *message.Message `json:"preconditionFalse,omitempty" yaml:"preconditionFalse,omitempty" hcl:"preconditionFalse,omitempty"`
	ConditionError    *message.Message `json:"conditionError,omitempty" yaml:"conditionError,omitempty" hcl:"conditionError,omitempty"`
	ConditionFalse    *message.Message `json:"conditionFalse,omitempty" yaml:"conditionFalse,omitempty" hcl:"conditionFalse,omitempty"`
	ConditionTrue     *message.Message `json:"conditionTrue,omitempty" yaml:"conditionTrue,omitempty" hcl:"conditionTrue,omitempty"`
}
