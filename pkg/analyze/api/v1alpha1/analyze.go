package v1alpha1

import (
	"github.com/replicatedcom/support-bundle/pkg/meta"
)

// Analyzer is a spec to generate one or more deployment assets
type AnalyzerSpec struct {
	AnalyzerShared `json:",inline" yaml:",inline" hcl:",inline"`

	KubernetesVersion     *KubernetesVersionAnalyzer     `json:"kubernetes.version,omitempty" yaml:"kubernetes.version,omitempty" hcl:"kubernetes.version,omitempty"`
	KubernetesTotalMemory *KubernetesTotalMemoryAnalyzer `json:"kubernetes.total-memory,omitempty" yaml:"kubernetes.total-memory,omitempty" hcl:"kubernetes.total-memory,omitempty"`

	Raw *RawAnalyzer `json:"raw,omitempty" yaml:"raw,omitempty" hcl:"raw,omitempty"`
}

type AnalyzerShared struct {
	Message     string     `json:"message,omitempty" yaml:"message,omitempty" hcl:"message,omitempty"`
	Meta        *meta.Meta `json:"meta,omitempty" yaml:"meta,omitempty" hcl:"meta,omitempty"`
	CollectRefs []meta.Ref `json:"collect_refs,omitempty" yaml:"collect_refs,omitempty" hcl:"collect_refs,omitempty"`
}

func (a AnalyzerSpec) GetAnalyzer() Analyzer {
	switch {
	case a.KubernetesVersion != nil:
		return a.KubernetesVersion
	case a.KubernetesTotalMemory != nil:
		return a.KubernetesTotalMemory
	case a.Raw != nil:
		return a.Raw
	default:
		return nil
	}
}
