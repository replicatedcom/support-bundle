package api_test

import (
	"testing"

	yaml "gopkg.in/yaml.v2"

	"github.com/stretchr/testify/assert"

	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	"github.com/stretchr/testify/require"
)

func TestSerialization(t *testing.T) {
	spec := `analyze:
  v1:
  - name: centos-min-version
    registerVariables:
    - os: {}
      name: os
    precondition:
      or:
      - eval: '{{repl eq .os "centos"}}'
      - eval: '{{repl eq .ref "centos"}}'
        ref: os
      - stringCompare:
          eq: centos
        ref: os
    condition:
      eval: '{{repl lt .osVersion "16.04" | not}}'
    messages:
      preconditionError:
        primary: Failed to detect OS
        detail: CentOS version must be at least 7.5
        severity: debug
      preconditionFalse:
        primary: OS is not CentOS
        detail: CentOS version must be at least 7.5
        severity: debug
      conditionError:
        primary: Failed to detect CentOS version
        detail: CentOS version must be at least 7.5
        severity: debug
      conditionFalse:
        primary: CentOS version {{repl .osVersion}}
        detail: CentOS version must be at least 7.5
        severity: warn
      conditionTrue:
        primary: CentOS version {{repl .osVersion}}
        detail: CentOS version must be at least 7.5
        severity: info
  - name: chef-client
    registerVariables:
    - collectRef:
        selector:
          analyze: ps-aux
      name: ps
    precondition:
      not:
        empty: {}
        ref: ps
    condition:
      regexpMatch:
        regexp: /chef-client
      ref: ps
    messages:
      conditionTrue:
        primary: Chef Client detected
        detail: The server must not be running the Chef Client
        severity: warn
`

	var obj api.Doc
	err := yaml.Unmarshal([]byte(spec), &obj)
	require.NoError(t, err)

	serialized, err := yaml.Marshal(obj)
	require.NoError(t, err)

	assert.Equal(t, spec, string(serialized))
}
