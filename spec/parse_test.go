package spec

import (
	"testing"

	"github.com/replicatedcom/support-bundle/types"
	"github.com/stretchr/testify/require"
)

// Test that yml files are parsed into spec lists properly
func TestParse(t *testing.T) {
	yml := `
specs:
  - builtin: core.loadavg
    raw: /raw/metrics/loadavg
    json: /json/metrics/loadavg.json
    human: /human/metrics/loadavg
  - builtin: docker.logs
    raw: /raw/containers/testExample/logs.txt
    config:
      container_id: testExample
  - builtin: docker.daemon
    raw: /raw/docker/
    json: /json/docker/
`

	specs, err := Parse([]byte(yml))
	require.NoError(t, err)

	require.Contains(t, specs, types.Spec{
		Builtin: "core.loadavg",
		Raw:     "/raw/metrics/loadavg",
		JSON:    "/json/metrics/loadavg.json",
		Human:   "/human/metrics/loadavg",
	})

	require.Contains(t, specs, types.Spec{
		Builtin: "docker.logs",
		Raw:     "/raw/containers/testExample/logs.txt",
		Config:  map[interface{}]interface{}{"container_id": "testExample"},
	})

	require.Contains(t, specs, types.Spec{
		Builtin: "docker.daemon",
		Raw:     "/raw/docker/",
		JSON:    "/json/docker/",
	})
}
