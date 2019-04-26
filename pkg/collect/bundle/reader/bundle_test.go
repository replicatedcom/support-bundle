package reader

import (
	"testing"

	"github.com/replicatedcom/support-bundle/pkg/collect/bundle/reader/testfixtures"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBundle_GetIndex(t *testing.T) {
	tests := []struct {
		name       string
		bundlePath string
		wantLen    int
		wantErr    bool
	}{
		{
			name:       "basic",
			bundlePath: "bundle",
			wantLen:    2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			f, err := fs.Create("./bundle.tgz")
			require.NoError(t, err)
			err = testfixtures.WriteBundle(f, tt.bundlePath)
			require.NoError(t, err)
			bundle, err := NewBundle(fs, "./bundle.tgz", "")
			require.NoError(t, err)
			index := bundle.GetIndex()
			assert.Len(t, index, tt.wantLen)
		})
	}
}
