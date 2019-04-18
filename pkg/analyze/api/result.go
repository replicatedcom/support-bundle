package api

import (
	"github.com/replicatedcom/support-bundle/pkg/analyze/api/common"
	"github.com/replicatedcom/support-bundle/pkg/analyze/message"
	"github.com/replicatedcom/support-bundle/pkg/meta"
)

type Result struct {
	meta.Meta `json:",inline" yaml:",inline" hcl:",inline"`

	Message      *message.Message       `json:"message" yaml:"message" hcl:"message"`
	Severity     common.Severity        `json:"severity" yaml:"severity" hcl:"severity"`
	AnalyzerSpec string                 `json:"analyzerSpec" yaml:"analyzerSpec" hcl:"analyzerSpec"`
	Variables    map[string]interface{} `json:"variables,omitempty" yaml:"variables,omitempty" hcl:"variables,omitempty"`
	Error        string                 `json:"error,omitempty" yaml:"error,omitempty" hcl:"error,omitempty"`
}
