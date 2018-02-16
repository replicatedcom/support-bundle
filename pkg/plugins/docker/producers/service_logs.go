package producers

import (
	"context"
	"io"

	"github.com/replicatedcom/support-bundle/pkg/types"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/replicatedcom/support-bundle/pkg/plugins/docker/util"
)

func (d *Docker) ServiceLogs(serviceID string, opts *dockertypes.ContainerLogsOptions) types.StreamsProducer {
	return func(ctx context.Context) (map[string]io.Reader, error) {
		if opts == nil {
			opts = &dockertypes.ContainerLogsOptions{}
		}
		if !opts.ShowStdout && !opts.ShowStderr {
			opts.ShowStdout = true
			opts.ShowStderr = true
		}
		opts.Timestamps = true
		reader, err := d.client.ServiceLogs(ctx, serviceID, *opts)
		if err != nil {
			return nil, err
		}
		return util.DemuxLogs(ctx, reader, "")
	}
}
