package resolver_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	v1 "github.com/replicatedcom/support-bundle/pkg/analyze/api/v1"
	. "github.com/replicatedcom/support-bundle/pkg/analyze/resolver"
	"github.com/replicatedcom/support-bundle/pkg/meta"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestResolverResolveSpec(t *testing.T) {
	tests := []struct {
		name   string
		files  map[string]string
		inline []string
		expect api.Doc
	}{
		{
			name: "resolve",
			files: map[string]string{
				"/spec/1.yml": `
analyze:
  v1:
    - kubernetes.total-memory:
        minimum: 10Gi
      collect_refs:
        - selector:
            analyze: resource-list-nodes`},
			inline: []string{`
analyze:
  v1:
    - kubernetes.version:
        semver_minimum: 1.10.0
      collect_refs:
        - selector:
            analyze: kubernetes-version`},

			expect: api.Doc{
				Analyze: api.Analyze{
					V1: []v1.AnalyzerSpec{
						{
							KubernetesTotalMemory: &v1.KubernetesTotalMemoryRequirement{
								Min: "10Gi",
							},
							AnalyzerShared: v1.AnalyzerShared{
								CollectRefs: []meta.Ref{{
									Selector: map[string]string{
										"analyze": "resource-list-nodes",
									},
								}},
							},
						},
						{
							KubernetesVersion: &v1.KubernetesVersionRequirement{
								SemverMin: "1.10.0",
							},
							AnalyzerShared: v1.AnalyzerShared{
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
			// TODO: test customer spec
			input := ResolverInput{
				Files:  files,
				Inline: test.inline,
			}
			spec, err := resolver.ResolveSpec(context.Background(), input)
			req.NoError(err)
			req.Equal(test.expect, spec)
		})
	}
}
