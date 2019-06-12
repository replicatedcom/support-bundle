package producers

import (
	"context"
	"io"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/docker/util"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func (d *Docker) ContainerLogs(container string, basename string, opts *dockertypes.ContainerLogsOptions) types.StreamsProducer {
	return func(ctx context.Context) (map[string]io.Reader, error) {
		if opts == nil {
			opts = &dockertypes.ContainerLogsOptions{}
		}
		if !opts.ShowStdout && !opts.ShowStderr {
			opts.ShowStdout = true
			opts.ShowStderr = true
		}
		opts.Timestamps = true
		opts.Follow = false
		reader, err := d.client.ContainerLogs(ctx, container, *opts)
		if err != nil {
			return nil, err
		}
		return util.DemuxLogs(ctx, reader, basename)
	}
}
