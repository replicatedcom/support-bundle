package resolver_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api/v1alpha1"
	. "github.com/replicatedcom/support-bundle/pkg/analyze/resolver"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
	"github.com/replicatedcom/support-bundle/pkg/meta"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestResolverResolveSpec(t *testing.T) {
	tests := []struct {
		name   string
		files  map[string]string
		inline []string
		expect api.Spec
	}{
		{
			name: "resolve",
			files: map[string]string{
				"/spec/1.yml": `
collect:
  v1:
    - kubernetes.resource-list:
        kind: nodes
      output_dir: /kubernetes/resource-list/nodes
      meta:
        labels:
          analyze: resource-list-nodes

analyze:
  v1alpha1:
    - kubernetes.total-memory:
        minimum: 10Gi
      collect_refs:
        - selector:
            analyze: resource-list-nodes`},
			inline: []string{`
collect:
  v1:
    - kubernetes.version: {}
      output_dir: /kubernetes/version
      meta:
        labels:
          analyze: kubernetes-version
analyze:
  v1alpha1:
    - kubernetes.version:
        semver_minimum: 1.10.0
      collect_refs:
        - selector:
            analyze: kubernetes-version`},

			expect: api.Spec{
				Collect: api.Collect{
					V1: []collecttypes.Spec{
						{
							KubernetesResourceList: &collecttypes.KubernetesResourceListOptions{
								Kind: "nodes",
							},
							SpecShared: collecttypes.SpecShared{
								OutputDir: "/kubernetes/resource-list/nodes",
								Meta: &meta.Meta{
									Labels: map[string]string{
										"analyze": "resource-list-nodes",
									},
								},
							},
						},
						{
							KubernetesVersion: &collecttypes.KubernetesVersionOptions{},
							SpecShared: collecttypes.SpecShared{
								OutputDir: "/kubernetes/version",
								Meta: &meta.Meta{
									Labels: map[string]string{
										"analyze": "kubernetes-version",
									},
								},
							},
						},
					},
				},
				Analyze: api.Analyze{
					V1Alpha1: []v1alpha1.AnalyzerSpec{
						{
							KubernetesTotalMemory: &v1alpha1.KubernetesTotalMemoryAnalyzer{
								Min: "10Gi",
							},
							AnalyzerShared: v1alpha1.AnalyzerShared{
								CollectRefs: []meta.Ref{{
									Selector: map[string]string{
										"analyze": "resource-list-nodes",
									},
								}},
							},
						},
						{
							KubernetesVersion: &v1alpha1.KubernetesVersionAnalyzer{
								SemverMin: "1.10.0",
							},
							AnalyzerShared: v1alpha1.AnalyzerShared{
								CollectRefs: []meta.Ref{{
									Selector: map[string]string{
										"analyze": "kubernetes-version",
									},
								}},
							},
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := require.New(t)

			memFS := afero.NewMemMapFs()

			var files []string
			for filename, data := range test.files {
				memFS.MkdirAll(filepath.Dir(filename), 0755)
				afero.WriteFile(memFS, filename, []byte(data), 0644)
				files = append(files, filename)
			}

			resolver := &Resolver{
				Logger: log.NewNopLogger(),
				Fs:     memFS,
			}
			spec, err := resolver.ResolveSpec(context.Background(), files, test.inline)
			req.NoError(err)
			req.Equal(test.expect, spec)
		})
	}
}
