package spec

import (
	"reflect"
	"testing"

	"github.com/replicatedcom/support-bundle/pkg/collect/types"
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

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name string
		yaml string
		want types.Doc
	}{
		{
			name: "specs",
			yaml: `
specs:
  - os.http-request:
      url: http://test.test
    output_dir: /os/http-request/`,
			want: types.Doc{
				Collect: types.Collect{
					V1: []types.Spec{
						{
							SpecShared:      types.SpecShared{OutputDir: "/os/http-request/"},
							CoreHTTPRequest: &types.CoreHTTPRequestOptions{URL: "http://test.test"},
						},
					},
				},
			},
		},
		{
			name: "collect.v1",
			yaml: `
collect:
  v1:
    - os.http-request:
        url: http://test.test
      output_dir: /os/http-request/`,
			want: types.Doc{
				Collect: types.Collect{
					V1: []types.Spec{
						{
							SpecShared:      types.SpecShared{OutputDir: "/os/http-request/"},
							CoreHTTPRequest: &types.CoreHTTPRequestOptions{URL: "http://test.test"},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Unmarshal([]byte(tt.yaml)); err != nil {
				t.Errorf("Unmarshal() error = %v", err)
			} else if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("Unmarshal() = %+v, want %+v", got, tt.want)
			}
		})
	}

}
