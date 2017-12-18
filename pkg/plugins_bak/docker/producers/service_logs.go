package producers

import (
	"context"
	"io"

	"github.com/replicatedcom/support-bundle/pkg/types"

	dockertypes "github.com/docker/docker/api/types"
)

func (d *Docker) ServiceLogs(name string) types.StreamProducer {
	return func(ctx context.Context) (io.Reader, error) {
		return d.client.ServiceLogs(ctx, name, dockertypes.ContainerLogsOptions{ShowStderr: true, ShowStdout: true})
	}
}
