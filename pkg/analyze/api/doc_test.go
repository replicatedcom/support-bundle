package api

import (
	"reflect"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api/v1alpha1"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func TestDoc_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		yaml string
		want Doc
	}{
		{
			name: "specs",
			yaml: `
specs:
  - kubernetes.resource-list:
      kind: nodes
    output_dir: /kubernetes/resource-list/nodes

analyze:
  v1alpha1:
    - kubernetes.total-memory:
        minimum: 1Gi`,
			want: Doc{
				Collect: Collect{
					V1: []collecttypes.Spec{
						{
							SpecShared:             collecttypes.SpecShared{OutputDir: "/kubernetes/resource-list/nodes"},
							KubernetesResourceList: &collecttypes.KubernetesResourceListOptions{Kind: "nodes"},
						},
					},
				},
				Analyze: Analyze{
					V1Alpha1: []v1alpha1.AnalyzerSpec{
						{
							KubernetesTotalMemory: &v1alpha1.KubernetesTotalMemoryAnalyzer{
								Min: "1Gi",
							},
						},
					},
				},
			},
		},
		{
			name: "collect.v1",
			yaml: `
collect:
  v1:
    - kubernetes.resource-list:
        kind: nodes
      output_dir: /kubernetes/resource-list/nodes

analyze:
  v1alpha1:
    - kubernetes.total-memory:
        minimum: 1Gi`,
			want: Doc{
				Collect: Collect{
					V1: []collecttypes.Spec{
						{
							SpecShared:             collecttypes.SpecShared{OutputDir: "/kubernetes/resource-list/nodes"},
							KubernetesResourceList: &collecttypes.KubernetesResourceListOptions{Kind: "nodes"},
						},
					},
				},
				Analyze: Analyze{
					V1Alpha1: []v1alpha1.AnalyzerSpec{
						{
							KubernetesTotalMemory: &v1alpha1.KubernetesTotalMemoryAnalyzer{
								Min: "1Gi",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Doc
			if err := yaml.Unmarshal([]byte(tt.yaml), &got); err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("json.Unmarshal() = %+v, want %+v", got, tt.want)
			}
		})
	}

}
