package v1

import (
	"errors"
	"text/template"
)

var (
	KubernetesVersionRequirementSpecTemplate *template.Template
)

func init() {
	KubernetesVersionRequirementSpecTemplate = mustAsset(
		"KubernetesVersionRequirementSpec", "requirements/kubernetes_version.yml")
}

type KubernetesVersionRequirement struct {
	SemverMin string `json:"semver_minimum" yaml:"semver_minimum" hcl:"semver_minimum"`
}

func (a *KubernetesVersionRequirement) Validate(spec AnalyzerSpec) error {
	if len(spec.CollectRefs) == 0 {
		return errors.New("collect_refs required")
	}
	if a.SemverMin == "" {
		return errors.New("semver_minimum required")
	}
	return nil
}

func (a *KubernetesVersionRequirement) GetRawSpec() (RawSpec, error) {
	return unmarshalSpec(KubernetesVersionRequirementSpecTemplate, a)
}
