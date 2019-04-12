package api

import (
	"github.com/replicatedcom/support-bundle/pkg/analyze/api/common"
	"github.com/replicatedcom/support-bundle/pkg/analyze/message"
)

type Result struct {
	Message      *message.Message       `json:"message" yaml:"message" hcl:"message"`
	Severity     common.Severity        `json:"severity" yaml:"severity" hcl:"severity"`
	AnalyzerSpec string                 `json:"analyzer_spec" yaml:"analyzer_spec" hcl:"analyzer_spec"`
	Variables    map[string]interface{} `json:"variables,omitempty" yaml:"variables,omitempty" hcl:"variables,omitempty"`
	Error        string                 `json:"error,omitempty" yaml:"error,omitempty" hcl:"error,omitempty"`
}
