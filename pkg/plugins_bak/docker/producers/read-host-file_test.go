package producers

import (
	"archive/tar"
	"context"
	"io"
	"io/ioutil"
	"testing"

	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadHostFile(t *testing.T) {
	t.SkipNow()
	c, err := client.NewEnvClient()
	require.NoError(t, err)

	d := New(c)
	producer := d.ReadHostFile("ubuntu:latest", "/go/src/github.com/replicatedcom/support-bundle/README.md")

	tarStream, err := producer(context.Background())
	require.NoError(t, err)

	tarReader := tar.NewReader(tarStream)
	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)

		n, err := io.Copy(ioutil.Discard, tarReader)
		require.NoError(t, err)
		require.Equal(t, hdr.Size, n)
		assert.Contains(t, hdr.Name, "README.md")
	}
}
