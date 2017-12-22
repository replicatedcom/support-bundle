package producers

import (
	"context"
	"io"

	"github.com/replicatedcom/support-bundle/pkg/plugins/docker/util"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) ContainerRun(opts types.DockerContainerRunOptions) types.StreamsProducer {
	return func(ctx context.Context) (map[string]io.Reader, error) {
		stdout, stderr, err := util.ContainerRun(ctx, d.client, opts.ContainerCreateConfig, opts.EnablePull)
		if err != nil {
			return nil, err
		}
		return map[string]io.Reader{
			"stdout.raw": stdout,
			"stderr.raw": stderr,
		}, nil
	}
}
