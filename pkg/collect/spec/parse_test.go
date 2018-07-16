package spec

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test that yml files are parsed into spec lists properly
func TestParse(t *testing.T) {
	yml := `
specs:
  - os.loadavg: {}
    output_dir: /metrics/loadavg/
  - docker.container-logs:
      container: testExample
      output_dir: /docker/container/logs/testExample/
  - docker.info:
      output_dir: /docker/info/
  - docker.container-inspect:
      output_dir: /docker/inspect/
`

	actual, err := Parse([]byte(yml))
	require.NoError(t, err)

	require.NotNil(t, actual[0].CoreLoadavg)
	require.Equal(t, "/metrics/loadavg/", actual[0].Shared().OutputDir)

	require.NotNil(t, actual[1].DockerContainerLogs)
	require.Equal(t, "/docker/container/logs/testExample/", actual[1].Shared().OutputDir)

	require.NotNil(t, actual[2].DockerInfo)
	require.Equal(t, "/docker/info/", actual[2].Shared().OutputDir)

	require.NotNil(t, actual[3].DockerContainerInspect)
	require.Equal(t, "/docker/inspect/", actual[3].Shared().OutputDir)

}
