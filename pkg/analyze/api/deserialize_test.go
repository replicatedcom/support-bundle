package api_test

import (
	"testing"

	"github.com/ghodss/yaml"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSerialization(t *testing.T) {
	spec := `analyze:
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
      os: {}
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
      name: ps
`

	var obj api.Doc
	err := yaml.Unmarshal([]byte(spec), &obj)
	require.NoError(t, err)

	serialized, err := yaml.Marshal(obj)
	require.NoError(t, err)

	assert.Equal(t, spec, string(serialized))
}
