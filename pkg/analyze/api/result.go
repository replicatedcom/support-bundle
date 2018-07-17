package api

import "github.com/replicatedcom/support-bundle/pkg/analyze/api/common"

type Result struct {
	Message      string                   `json:"message" yaml:"message" hcl:"message"`
	Requirement  string                   `json:"requirement,omitempty" yaml:"requirement,omitempty" hcl:"requirement,omitempty"`
	Severity     common.Severity          `json:"severity" yaml:"severity" hcl:"severity"`
	AnalyzerSpec Analyze                  `json:"analyzer_spec" yaml:"analyzer_spec" hcl:"analyzer_spec"`
	Vars         []map[string]interface{} `json:"variables,omitempty" yaml:"variables,omitempty" hcl:"variables,omitempty"`
	Error        string                   `json:"error,omitempty" yaml:"error,omitempty" hcl:"error,omitempty"`
}
