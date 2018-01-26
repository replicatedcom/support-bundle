package producers

import (
	"context"
	"io"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/replicatedcom/support-bundle/pkg/plugins/docker/util"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) Logs(containerID string) types.StreamsProducer {
	return func(ctx context.Context) (map[string]io.Reader, error) {
		options := dockertypes.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		}

		reader, err := d.client.ContainerLogs(ctx, containerID, options)
		if err != nil {
			return nil, err
		}

		return util.DemuxLogs(ctx, reader, "")
	}
}

func (d *Docker) LogsName(containerName string) types.StreamsProducer {
	return func(ctx context.Context) (map[string]io.Reader, error) {
		containerID, err := d.getContainerID(ctx, containerName)
		if err != nil {
			return nil, err
		}
		return d.Logs(containerID)(ctx)
	}
}
