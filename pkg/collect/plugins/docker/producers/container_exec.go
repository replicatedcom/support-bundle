package producers

import (
	"context"
	"io"

	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/docker/util"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func (d *Docker) ContainerExec(opts types.DockerContainerExecOptions) types.StreamsProducer {
	return func(ctx context.Context) (map[string]io.Reader, error) {
		stdoutR, stderrR, _, err := util.ContainerExec(ctx, d.client, opts.Container, opts.ExecConfig)
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
