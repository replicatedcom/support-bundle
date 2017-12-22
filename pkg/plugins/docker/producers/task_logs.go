package producers

import (
	"context"
	"io"

	"github.com/replicatedcom/support-bundle/pkg/types"

	dockertypes "github.com/docker/docker/api/types"
)

func (d *Docker) TaskLogs(taskID string, opts *dockertypes.ContainerLogsOptions) types.StreamProducer {
	return func(ctx context.Context) (io.Reader, error) {
		if opts == nil {
			opts = &dockertypes.ContainerLogsOptions{}
		}
		if !opts.ShowStdout && !opts.ShowStderr {
			opts.ShowStdout = true
			opts.ShowStderr = true
		}
		opts.Timestamps = true
		return d.client.TaskLogs(ctx, taskID, *opts)
	}
}
