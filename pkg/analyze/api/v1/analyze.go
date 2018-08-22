package v1

import (
	"github.com/replicatedcom/support-bundle/pkg/analyze/api/common"
	"github.com/replicatedcom/support-bundle/pkg/meta"
)

// AnalyzerSpec is a spec to analyze one or more collect assets
type AnalyzerSpec struct {
	AnalyzerShared `json:",inline" yaml:",inline" hcl:",inline"`

	Raw *RawRequirement `json:"raw,omitempty" yaml:"raw,omitempty" hcl:"raw,omitempty"`

	DockerVersion *DockerVersionRequirement `json:"docker.version,omitempty" yaml:"docker.version,omitempty" hcl:"docker.version,omitempty"`

	KubernetesVersion     *KubernetesVersionRequirement     `json:"kubernetes.version,omitempty" yaml:"kubernetes.version,omitempty" hcl:"kubernetes.version,omitempty"`
	KubernetesTotalMemory *KubernetesTotalMemoryRequirement `json:"kubernetes.total-memory,omitempty" yaml:"kubernetes.total-memory,omitempty" hcl:"kubernetes.total-memory,omitempty"`
}

// RawSpec is a translated spec from an AnalyzerSpec
type RawSpec struct {
	AnalyzerShared `json:",inline" yaml:",inline" hcl:",inline"`

	Raw *RawRequirement `json:"raw,omitempty" yaml:"raw,omitempty" hcl:"raw,omitempty"`
}

type AnalyzerShared struct {
	Message     string          `json:"message,omitempty" yaml:"message,omitempty" hcl:"message,omitempty"`
	Severity    common.Severity `json:"severity,omitempty" yaml:"severity,omitempty" hcl:"severity,omitempty"`
	Meta        *meta.Meta      `json:"meta,omitempty" yaml:"meta,omitempty" hcl:"meta,omitempty"`
	CollectRefs []meta.Ref      `json:"collect_refs,omitempty" yaml:"collect_refs,omitempty" hcl:"collect_refs,omitempty"`
}

func (a AnalyzerSpec) GetRequirement() Requirement {
	switch {

	case a.Raw != nil:
		return a.Raw

	case a.DockerVersion != nil:
		return a.DockerVersion

	case a.KubernetesVersion != nil:
		return a.KubernetesVersion
	case a.KubernetesTotalMemory != nil:
		return a.KubernetesTotalMemory

	default:
		return nil
	}
}
