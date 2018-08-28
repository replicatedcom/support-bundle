package plans

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func TestStreamsSourceScrubber(t *testing.T) {
	type fields struct {
		Filepath     string
		StreamFormat string
		RawScrubber  func([]byte) []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "basic",
			fields: fields{
				Filepath: "testfixtures/retraced-api.input.json",
			},
			want: "testfixtures/retraced-api.input.json",
		},
		{
			name: "scrub",
			fields: fields{
				Filepath: "testfixtures/retraced-api.input.json",
				RawScrubber: mustRawScrubber(t, &types.Scrub{
					Regex:   "(?m)(\"(?:POSTGRES_PASSWORD|ELASTICSEARCH_NODES|HMAC_SECRET_ADMIN|HMAC_SECRET_VIEWER|BOOTSTRAP_API_TOKEN)=)([^\"]*)(\",?)",
					Replace: "${1}***HIDDEN***${3}",
				}),
			},
			want: "testfixtures/retraced-api.scrub.json",
		},
		{
			name: "scrub tgz",
			fields: fields{
				Filepath:     "testfixtures/retraced-api.input.json",
				StreamFormat: StreamFormatTar,
				RawScrubber: mustRawScrubber(t, &types.Scrub{
					Regex:   `(?m)("(?:POSTGRES_PASSWORD|ELASTICSEARCH_NODES|HMAC_SECRET_ADMIN|HMAC_SECRET_VIEWER|BOOTSTRAP_API_TOKEN)=)([^"]*)(",?)`,
					Replace: "${1}***HIDDEN***${3}",
				}),
			},
			want: "testfixtures/retraced-api.scrub.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootDir := tempDir(t, "")
			defer os.RemoveAll(rootDir)

			task := &StreamsSource{
				Spec: types.Spec{
					CoreReadFile: &types.CoreReadFileOptions{
						Filepath: tt.fields.Filepath,
					},
				},
				Producer: func(context.Context) (map[string]io.Reader, error) {
					if tt.fields.StreamFormat == StreamFormatTar {
						return map[string]io.Reader{
							"": tgzReader(t, tt.fields.Filepath),
						}, nil
					}
					return map[string]io.Reader{
						"": reader(t, tt.fields.Filepath),
					}, nil
				},
				StreamFormat: tt.fields.StreamFormat,
				RawScrubber:  tt.fields.RawScrubber,
				RawPath:      "actual",
			}
			got := task.Exec(context.Background(), rootDir)
			require.Len(t, got, 1)
			require.NoError(t, got[0].Error)
			expected := mustAsset(t, tt.want)
			actual := mustAsset(t, filepath.Join(rootDir, got[0].Path))
			assert.Equal(t, expected, actual)
		})
	}
}
