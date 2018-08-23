package reader

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/mholt/archiver"
	"github.com/replicatedcom/support-bundle/pkg/meta"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBundle_ReaderFromRef(t *testing.T) {
	tests := []struct {
		name       string
		bundlePath string
		ref        meta.Ref
		want       string
		wantErr    bool
	}{
		{
			name:       "basic",
			bundlePath: "testfixtures/bundle",
			ref: meta.Ref{
				Selector: meta.Selector{
					"analyze": "kubernetes-version",
				},
			},
			want: `{
  "major": "1",
  "minor": "8+",
  "gitVersion": "v1.8.10-gke.0",
  "gitCommit": "16ebd0de8e0ab2d1ef86d5b16ab1899b624a77cd",
  "gitTreeState": "clean",
  "buildDate": "2018-03-20T20:21:01Z",
  "goVersion": "go1.8.3b4",
  "compiler": "gc",
  "platform": "linux/amd64"
}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			f, err := fs.Create("./bundle.tgz")
			require.NoError(t, err)
			err = writeBundle(f, tt.bundlePath)
			require.NoError(t, err)
			bundle, err := NewBundle(fs, "./bundle.tgz")
			require.NoError(t, err)
			r, err := bundle.ReaderFromRef(tt.ref)
			require.NoError(t, err)
			got, err := ioutil.ReadAll(r)
			require.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func writeBundle(w io.Writer, bundlePath string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir(bundlePath); err != nil {
		return err
	}
	defer os.Chdir(cwd)

	var filePaths []string
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return err
	}
	for _, info := range files {
		filePaths = append(filePaths, info.Name())
	}

	return archiver.TarGz.Write(w, filePaths)
}
