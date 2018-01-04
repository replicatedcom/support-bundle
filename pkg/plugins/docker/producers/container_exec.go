package producers

import (
	"context"
	"io"

	"github.com/replicatedcom/support-bundle/pkg/plugins/docker/util"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) ContainerExec(opts types.DockerContainerExecOptions) types.StreamsProducer {
	return func(ctx context.Context) (map[string]io.Reader, error) {
		stdout, stderr, _, err := util.ContainerExec(ctx, d.client, opts.Container, opts.ExecConfig)
		if err != nil {
			return nil, err
		}
		return map[string]io.Reader{
			"stdout.raw": stdout,
			"stderr.raw": stderr,
		}, nil
	}
}
