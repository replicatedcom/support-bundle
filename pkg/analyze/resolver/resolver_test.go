package resolver_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/replicatedcom/support-bundle/pkg/analyze/condition"
	"github.com/replicatedcom/support-bundle/pkg/analyze/insight"
	"github.com/replicatedcom/support-bundle/pkg/analyze/variable"
	"github.com/replicatedcom/support-bundle/pkg/meta"

	"github.com/go-kit/kit/log"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api/common"
	v1 "github.com/replicatedcom/support-bundle/pkg/analyze/api/v1"
	. "github.com/replicatedcom/support-bundle/pkg/analyze/resolver"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestResolverResolveSpec(t *testing.T) {
	eqOsCentosEval := condition.EvalCondition(`{{repl eq .os "centos"}}`)
	eqRefCentosEval := condition.EvalCondition(`{{repl eq .Ref "centos"}}`)
	osVersionGte1604Eval := condition.EvalCondition(`{{repl lt .osVersion "16.04" | not}}`)

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
  - evaluateConditions:
    - condition:
        or:
        - eval: '{{repl eq .os "centos"}}'
        - eval: '{{repl eq .Ref "centos"}}'
          variableRef: os
        - stringCompare:
            eq: centos
          variableRef: os
      insightOnError:
        detail: CentOS version must be at least 7.5
        primary: Failed to detect OS
        severity: debug
      insightOnFalse:
        detail: CentOS version must be at least 7.5
        primary: OS is not CentOS
        severity: debug
    - condition:
        eval: '{{repl lt .osVersion "16.04" | not}}'
      insightOnError:
        detail: CentOS version must be at least 7.5
        primary: Failed to detect CentOS version
        severity: debug
      insightOnFalse:
        detail: CentOS version must be at least 7.5
        primary: CentOS version {{repl .osVersion}}
        severity: warn
    insight:
      detail: CentOS version must be at least 7.5
      primary: CentOS version {{repl .osVersion}}
      severity: info
    name: centos-min-version
    registerVariables:
    - name: os
      os: {}`},
			inline: []string{`
analyze:
  v1:
  - evaluateConditions:
    - condition:
        not:
          empty: {}
          variableRef: ps
    - condition:
        regexpMatch:
          regexp: /chef-client
        variableRef: ps
    insight:
      detail: The server must not be running the Chef Client
      primary: Chef Client detected
      severity: warn
    name: chef-client
    registerVariables:
    - collectRef:
        selector:
          analyze: ps-aux
      name: ps`},

			expect: api.Doc{
				Analyze: api.Analyze{
					V1: []v1.Analyzer{
						{
							Meta: meta.Meta{
								Name: "centos-min-version",
							},
							RegisterVariables: []v1.Variable{
								{
									Meta: meta.Meta{
										Name: "os",
									},
									Os: &variable.Os{},
								},
							},
							EvaluateConditions: []v1.EvaluateCondition{
								{
									Condition: v1.Condition{
										Or: &v1.OrPredicate{
											{
												EvalCondition: &eqOsCentosEval,
											},
											{
												EvalCondition: &eqRefCentosEval,
												VariableRef:   "os",
											},
											{
												StringCompare: &condition.StringCompare{
													Compare: condition.Compare{
														Eq: "centos",
													},
												},
												VariableRef: "os",
											},
										},
									},
									InsightOnError: &insight.Insight{
										Primary:  "Failed to detect OS",
										Detail:   "CentOS version must be at least 7.5",
										Severity: common.SeverityDebug,
									},
									InsightOnFalse: &insight.Insight{
										Primary:  "OS is not CentOS",
										Detail:   "CentOS version must be at least 7.5",
										Severity: common.SeverityDebug,
									},
								},
								{
									Condition: v1.Condition{
										EvalCondition: &osVersionGte1604Eval,
									},
									InsightOnError: &insight.Insight{
										Primary:  "Failed to detect CentOS version",
										Detail:   "CentOS version must be at least 7.5",
										Severity: common.SeverityDebug,
									},
									InsightOnFalse: &insight.Insight{
										Primary:  "CentOS version {{repl .osVersion}}",
										Detail:   "CentOS version must be at least 7.5",
										Severity: common.SeverityWarn,
									},
								},
							},
							Insight: &insight.Insight{
								Primary:  "CentOS version {{repl .osVersion}}",
								Detail:   "CentOS version must be at least 7.5",
								Severity: common.SeverityInfo,
							},
						},
						{
							Meta: meta.Meta{
								Name: "chef-client",
							},
							RegisterVariables: []v1.Variable{
								{
									Meta: meta.Meta{
										Name: "ps",
									},
									CollectRef: &variable.CollectRef{
										Ref: meta.Ref{
											Selector: meta.Selector{
												"analyze": "ps-aux",
											},
										},
									},
								},
							},
							EvaluateConditions: []v1.EvaluateCondition{
								{
									Condition: v1.Condition{
										Not: &v1.NotPredicate{
											Condition: v1.Condition{
												Empty:       &condition.Empty{},
												VariableRef: "ps",
											},
										},
									},
								},
								{
									Condition: v1.Condition{
										RegexpMatch: &condition.RegexpMatch{
											Regexp: "/chef-client",
										},
										VariableRef: "ps",
									},
								},
							},
							Insight: &insight.Insight{
								Primary:  "Chef Client detected",
								Detail:   "The server must not be running the Chef Client",
								Severity: common.SeverityWarn,
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
			input := Input{
				Files:  files,
				Inline: test.inline,
			}
			spec, err := resolver.ResolveSpec(context.Background(), input, true)
			req.NoError(err)
			req.Equal(test.expect, spec)
		})
	}
}
