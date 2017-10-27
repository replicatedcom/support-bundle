package producers

import (
	"context"
	"io"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/replicatedcom/support-bundle/pkg/types"
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

func (d *Docker) LogsName(containerName string) types.StreamProducer {
	return func(ctx context.Context) (io.Reader, error) {
		containerID, err := d.getContainerID(ctx, containerName)
		if err != nil {
			return nil, err
		}
		return d.Logs(containerID)(ctx)
	}
}
