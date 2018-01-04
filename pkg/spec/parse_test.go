package spec

import (
	"testing"

	"github.com/replicatedcom/support-bundle/pkg/types"
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
  - docker.info: {}
    output_dir: /docker/info/
`

	actual, err := Parse([]byte(yml))
	require.NoError(t, err)

	expected := []types.Spec{
		types.Spec{
			CoreLoadavg: &types.CoreLoadavgOptions{},
			OutputDir:   "/metrics/loadavg/",
		},
		types.Spec{
			DockerContainerLogs: &types.DockerContainerLogsOptions{
				Container: "testExample",
			},
			OutputDir: "/docker/container/logs/testExample/",
		},
		types.Spec{
			DockerInfo: &types.DockerInfoOptions{},
			OutputDir:  "/docker/info/",
		},
	}

	require.EqualValues(t, expected, actual)
}
