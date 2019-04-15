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
      - eval: '{{repl eq .Ref "centos"}}'
        ref: os
      - ref: os
        stringCompare:
          eq: centos
    registerVariables:
    - name: os
      os: {}
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
      name: ps
`

	var obj api.Doc
	err := yaml.Unmarshal([]byte(spec), &obj)
	require.NoError(t, err)

	serialized, err := yaml.Marshal(obj)
	require.NoError(t, err)

	assert.Equal(t, spec, string(serialized))
}
