package producers

import (
	"context"
	"io"

	"github.com/replicatedcom/support-bundle/pkg/plugins/docker/util"
	"github.com/replicatedcom/support-bundle/pkg/types"

	dockertypes "github.com/docker/docker/api/types"
)

func (d *Docker) ServiceLogs(name string) types.StreamsProducer {
	return func(ctx context.Context) (map[string]io.Reader, error) {
		reader, err := d.client.ServiceLogs(ctx, name, dockertypes.ContainerLogsOptions{ShowStderr: true, ShowStdout: true})

		if err != nil {
			return nil, err
		}

		return util.DemuxLogs(ctx, reader, "")
	}
}
