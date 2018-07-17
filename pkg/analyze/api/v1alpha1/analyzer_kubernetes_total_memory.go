package v1alpha1

import (
	"errors"
	"text/template"

	"github.com/replicatedcom/support-bundle/pkg/analyze/api/common"
)

var (
	KubernetesTotalMemoryAnalyzerSpecTemplate *template.Template
)

func init() {
	KubernetesTotalMemoryAnalyzerSpecTemplate = mustAsset(
		"KubernetesTotalMemoryAnalyzerSpec", "specs/analyzer_kubernetes_total_memory.yml")
}

type KubernetesTotalMemoryAnalyzer struct {
	Severity common.Severity `json:"severity,omitempty" yaml:"severity,omitempty" hcl:"severity,omitempty"`
	Min      string          `json:"minimum" yaml:"minimum" hcl:"minimum"`
}

func (a *KubernetesTotalMemoryAnalyzer) Validate(spec AnalyzerSpec) error {
	if len(spec.CollectRefs) == 0 {
		return errors.New("collect_refs required")
	}
	if a.Min == "" {
		return errors.New("minimum required")
	}
	return nil
}

func (a *KubernetesTotalMemoryAnalyzer) GetSpec() (AnalyzerSpec, error) {
	return unmarshalSpec(KubernetesTotalMemoryAnalyzerSpecTemplate, a)
}
