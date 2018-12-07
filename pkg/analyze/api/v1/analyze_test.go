package v1_test

import (
	"testing"

	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	. "github.com/replicatedcom/support-bundle/pkg/analyze/api/v1"
	"github.com/replicatedcom/support-bundle/pkg/meta"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

func TestAnalyzeDeserialize(t *testing.T) {
	tests := []struct {
		name   string
		yaml   string
		expect AnalyzerSpec
	}{
		{
			name: "kubernetes.version",
			yaml: `
---
analyze:
  v1:
    - kubernetes.version:
        semver_minimum: 1.10.0
      collect_refs:
        - selector:
            analyze: kubernetes-version`,

			expect: AnalyzerSpec{
				KubernetesVersion: &KubernetesVersionRequirement{
					SemverMin: "1.10.0",
				},
				AnalyzerShared: AnalyzerShared{
					CollectRefs: []meta.Ref{{
						Selector: map[string]string{
							"analyze": "kubernetes-version",
						},
					}},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := require.New(t)

			var doc api.Doc

			err := yaml.Unmarshal([]byte(test.yaml), &doc)
			req.NoError(err)
			req.Len(doc.Analyze.V1, 1)
			req.Equal(test.expect, doc.Analyze.V1[0])
		})
	}
}
