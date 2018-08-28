package plans

import (
	"archive/tar"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/replicatedcom/support-bundle/pkg/collect/types"
	"github.com/stretchr/testify/require"
)

func tempDir(t *testing.T, prefix string) string {
	dir, err := ioutil.TempDir("", prefix)
	require.NoError(t, err)
	return dir
}

func mustAsset(t *testing.T, filename string) string {
	b, err := ioutil.ReadFile(filename)
	require.NoError(t, err)
	return string(b)
}

func mustRawScrubber(t *testing.T, scrubSpec *types.Scrub) types.BytesScrubber {
	scrubber, err := RawScrubber(scrubSpec)
	require.NoError(t, err)
	return scrubber
}

func reader(t *testing.T, filename string) io.Reader {
	return strings.NewReader(mustAsset(t, filename))
}

func tgzReader(t *testing.T, filename string) io.Reader {
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		asset := mustAsset(t, filename)
		tw := tar.NewWriter(pw)
		defer tw.Close()
		hdr := &tar.Header{
			Name: filename,
			Mode: 0600,
			Size: int64(len(asset)),
		}
		err := tw.WriteHeader(hdr)
		require.NoError(t, err)
		_, err = tw.Write([]byte(asset))
		require.NoError(t, err)
	}()
	return pr
}
