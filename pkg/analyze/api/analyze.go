package api

import v1 "github.com/replicatedcom/support-bundle/pkg/analyze/api/v1"

// Analyze is the top level analyze object
type Analyze struct {
	V1 []v1.Analyzer `json:"v1,omitempty" yaml:"v1,omitempty" hcl:"v1,omitempty"`
}
