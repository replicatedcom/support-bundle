package plans

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStructuredSourceScrubber(t *testing.T) {
	type fields struct {
		Source      interface{}
		RawScrubber func([]byte) []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{
			name: "basic",
			fields: fields{
				Source: containerInspectSource(t, "testfixtures/retraced-api.input.json"),
			},
			want: containerInspectSource(t, "testfixtures/retraced-api.input.json"),
		},
		{
			name: "scrub",
			fields: fields{
				Source: containerInspectSource(t, "testfixtures/retraced-api.input.json"),
				RawScrubber: mustRawScrubber(t, &types.Scrub{
					Regex:   "(?m)(\"(?:POSTGRES_PASSWORD|ELASTICSEARCH_NODES|HMAC_SECRET_ADMIN|HMAC_SECRET_VIEWER|BOOTSTRAP_API_TOKEN)=)([^\"]*)(\",?)",
					Replace: "${1}***HIDDEN***${3}",
				}),
			},
			want: containerInspectSource(t, "testfixtures/retraced-api.scrub.json"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootDir := tempDir(t, "")
			defer os.RemoveAll(rootDir)

			task := &StructuredSource{
				Spec: types.Spec{
					DockerContainerInspect: &types.DockerContainerInspectOptions{
						Container: "my-container",
					},
				},
				Producer: func(context.Context) (interface{}, error) {
					return tt.fields.Source, nil
				},
				RawScrubber: tt.fields.RawScrubber,
				RawPath:     "actual",
			}
			got := task.Exec(context.Background(), rootDir)
			require.Len(t, got, 1)
			require.NoError(t, got[0].Error)
			expected, err := json.Marshal(tt.want)
			require.NoError(t, err)
			actual := mustAsset(t, filepath.Join(rootDir, got[0].Path))
			assert.JSONEq(t, string(expected), actual)
		})
	}
}

func containerInspectSource(t *testing.T, filename string) interface{} {
	asset := mustAsset(t, filename)
	var container dockertypes.ContainerJSON
	err := json.Unmarshal([]byte(asset), &container)
	require.NoError(t, err)
	return container
}
