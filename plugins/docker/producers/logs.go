package producers

import (
	"context"
	"io"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/replicatedcom/support-bundle/types"
)

func (d *Docker) Logs(containerID string) types.StreamProducer {
	return func(ctx context.Context) (io.Reader, error) {
		options := dockertypes.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		}
		return d.client.ContainerLogs(ctx, containerID, options)
	}
}
