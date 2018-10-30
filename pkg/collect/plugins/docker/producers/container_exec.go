package producers

import (
	"context"
	"io"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/docker/util"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func (d *Docker) ContainerExec(containerID string, opts dockertypes.ExecConfig) types.StreamsProducer {
	return func(ctx context.Context) (map[string]io.Reader, error) {
		stdoutR, stderrR, _, err := util.ContainerExec(ctx, d.client, containerID, opts)
		// FIXME: stdoutR, stderrR never closed
		if err != nil {
			return nil, err
		}
		return map[string]io.Reader{
			"stdout.raw": stdoutR,
			"stderr.raw": stderrR,
		}, nil
	}
}
