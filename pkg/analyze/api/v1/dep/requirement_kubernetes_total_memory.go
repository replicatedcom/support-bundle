package v1

import (
	"errors"
	"text/template"
)

var (
	KubernetesTotalMemoryRequirementSpecTemplate *template.Template
)

func init() {
	KubernetesTotalMemoryRequirementSpecTemplate = mustAsset(
		"KubernetesTotalMemoryRequirementSpec", "requirements/kubernetes_total_memory.yml")
}

type KubernetesTotalMemoryRequirement struct {
	Min string `json:"minimum" yaml:"minimum" hcl:"minimum"`
}

func (a *KubernetesTotalMemoryRequirement) Validate(spec AnalyzerSpec) error {
	if len(spec.CollectRefs) == 0 {
		return errors.New("collect_refs required")
	}
	if a.Min == "" {
		return errors.New("minimum required")
	}
	return nil
}

func (a *KubernetesTotalMemoryRequirement) GetRawSpec() (RawSpec, error) {
	return unmarshalSpec(KubernetesTotalMemoryRequirementSpecTemplate, a)
}
