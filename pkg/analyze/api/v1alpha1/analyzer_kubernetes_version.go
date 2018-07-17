package v1alpha1

import (
	"errors"
	"text/template"

	"github.com/replicatedcom/support-bundle/pkg/analyze/api/common"
)

var (
	KubernetesVersionAnalyzerSpecTemplate *template.Template
)

func init() {
	KubernetesVersionAnalyzerSpecTemplate = mustAsset(
		"KubernetesVersionAnalyzerSpec", "specs/analyzer_kubernetes_version.yml")
}

type KubernetesVersionAnalyzer struct {
	Severity  common.Severity `json:"severity,omitempty" yaml:"severity,omitempty" hcl:"severity,omitempty"`
	SemverMin string          `json:"semver_minimum" yaml:"semver_minimum" hcl:"semver_minimum"`
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

func (a *KubernetesVersionAnalyzer) GetSpec() (AnalyzerSpec, error) {
	return unmarshalSpec(KubernetesVersionAnalyzerSpecTemplate, a)
}
