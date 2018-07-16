package producers

import (
	"context"

	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func (d *Docker) ContainerLs(opts types.DockerContainerLsOptions) types.StructuredProducer {
	return func(ctx context.Context) (interface{}, error) {
		return d.client.ContainerList(ctx, opts.ToDockerContainerListOptions())
	}
}
