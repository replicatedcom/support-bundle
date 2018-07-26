package v1

import (
	"errors"
	"text/template"
)

var (
	KubernetesVersionAnalyzerSpecTemplate *template.Template
)

func init() {
	KubernetesVersionAnalyzerSpecTemplate = mustAsset(
		"KubernetesVersionAnalyzerSpec", "specs/analyzer_kubernetes_version.yml")
}

type KubernetesVersionAnalyzer struct {
	SemverMin string `json:"semver_minimum" yaml:"semver_minimum" hcl:"semver_minimum"`
}

func (a *KubernetesVersionAnalyzer) Validate(spec AnalyzerSpec) error {
	if len(spec.CollectRefs) == 0 {
		return errors.New("collect_refs required")
	}
	if a.SemverMin == "" {
		return errors.New("semver_minimum required")
	}
	return nil
}

func (a *KubernetesVersionAnalyzer) GetRawSpec() (RawSpec, error) {
	return unmarshalSpec(KubernetesVersionAnalyzerSpecTemplate, a)
}
