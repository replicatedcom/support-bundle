package api

import "github.com/replicatedcom/support-bundle/pkg/analyze/api/v1alpha1"

// Analyze is the top level analyze object
type Analyze struct {
	V1Alpha1 []v1alpha1.AnalyzerSpec `json:"v1alpha1,omitempty" yaml:"v1alpha1,omitempty" hcl:"v1alpha1,omitempty"`
}
