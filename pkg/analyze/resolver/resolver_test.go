package resolver_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/replicatedcom/support-bundle/pkg/analyze/condition"
	"github.com/replicatedcom/support-bundle/pkg/analyze/message"
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
	eqRefCentosEval := condition.EvalCondition(`{{repl eq .ref "centos"}}`)
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
  - condition:
      eval: '{{repl lt .osVersion "16.04" | not}}'
    messages:
      conditionError:
        detail: CentOS version must be at least 7.5
        primary: Failed to detect CentOS version
        severity: debug
      conditionFalse:
        detail: CentOS version must be at least 7.5
        primary: CentOS version {{repl .osVersion}}
        severity: warn
      conditionTrue:
        detail: CentOS version must be at least 7.5
        primary: CentOS version {{repl .osVersion}}
        severity: info
      preconditionError:
        detail: CentOS version must be at least 7.5
        primary: Failed to detect OS
        severity: debug
      preconditionFalse:
        detail: CentOS version must be at least 7.5
        primary: OS is not CentOS
        severity: debug
    name: centos-min-version
    precondition:
      or:
      - eval: '{{repl eq .os "centos"}}'
      - eval: '{{repl eq .ref "centos"}}'
        ref: os
      - ref: os
        stringCompare:
          eq: centos
    registerVariables:
    - name: os
      os: {}`},
			inline: []string{`
analyze:
  v1:
  - condition:
      ref: ps
      regexpMatch:
        regexp: /chef-client
    messages:
      conditionTrue:
        detail: The server must not be running the Chef Client
        primary: Chef Client detected
        severity: warn
    name: chef-client
    precondition:
      not:
        empty: {}
        ref: ps
    registerVariables:
    - collectRef:
        selector:
          analyze: ps-aux
      name: ps`},

			expect: api.Doc{
				Analyze: api.Analyze{
					V1: []v1.Analyzer{
						{
							Name: "centos-min-version",
							RegisterVariables: []v1.Variable{
								{
									Name: "os",
									Os:   &variable.Os{},
								},
							},
							Precondition: &v1.Condition{
								Or: &v1.OrPredicate{
									{
										EvalCondition: &eqOsCentosEval,
									},
									{
										EvalCondition: &eqRefCentosEval,
										Ref:           "os",
									},
									{
										StringCompare: &condition.StringCompare{
											Compare: condition.Compare{
												Eq: "centos",
											},
										},
										Ref: "os",
									},
								},
							},
							Condition: &v1.Condition{
								EvalCondition: &osVersionGte1604Eval,
							},
							Messages: v1.Messages{
								PreconditionError: &message.Message{
									Primary:  "Failed to detect OS",
									Detail:   "CentOS version must be at least 7.5",
									Severity: common.SeverityDebug,
								},
								PreconditionFalse: &message.Message{
									Primary:  "OS is not CentOS",
									Detail:   "CentOS version must be at least 7.5",
									Severity: common.SeverityDebug,
								},
								ConditionError: &message.Message{
									Primary:  "Failed to detect CentOS version",
									Detail:   "CentOS version must be at least 7.5",
									Severity: common.SeverityDebug,
								},
								ConditionFalse: &message.Message{
									Primary:  "CentOS version {{repl .osVersion}}",
									Detail:   "CentOS version must be at least 7.5",
									Severity: common.SeverityWarn,
								},
								ConditionTrue: &message.Message{
									Primary:  "CentOS version {{repl .osVersion}}",
									Detail:   "CentOS version must be at least 7.5",
									Severity: common.SeverityInfo,
								},
							},
						},
						{
							Name: "chef-client",
							RegisterVariables: []v1.Variable{
								{
									Name: "ps",
									CollectRef: &variable.CollectRef{
										Ref: meta.Ref{
											Selector: meta.Selector{
												"analyze": "ps-aux",
											},
										},
									},
								},
							},
							Precondition: &v1.Condition{
								Not: &v1.NotPredicate{
									Condition: v1.Condition{
										Empty: &condition.Empty{},
										Ref:   "ps",
									},
								},
							},
							Condition: &v1.Condition{
								RegexpMatch: &condition.RegexpMatch{
									Regexp: "/chef-client",
								},
								Ref: "ps",
							},
							Messages: v1.Messages{
								ConditionTrue: &message.Message{
									Primary:  "Chef Client detected",
									Detail:   "The server must not be running the Chef Client",
									Severity: common.SeverityWarn,
								},
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
			spec, err := resolver.ResolveSpec(context.Background(), input)
			req.NoError(err)
			req.Equal(test.expect, spec)
		})
	}
}
