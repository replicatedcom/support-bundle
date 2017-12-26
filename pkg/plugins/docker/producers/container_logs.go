package producers

import (
	"context"
	"io"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) ContainerLogs(container string, opts *dockertypes.ContainerLogsOptions) types.StreamProducer {
	return func(ctx context.Context) (io.Reader, error) {
		if opts == nil {
			opts = &dockertypes.ContainerLogsOptions{}
		}
		if !opts.ShowStdout && !opts.ShowStderr {
			opts.ShowStdout = true
			opts.ShowStderr = true
		}
		opts.Timestamps = true
		return d.client.ContainerLogs(ctx, container, *opts)
	}
}
