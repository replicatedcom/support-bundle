package producers

import (
	"context"

	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func (d *Docker) NodeLs(opts types.DockerNodeLsOptions) types.StructuredProducer {
	return func(ctx context.Context) (interface{}, error) {
		return d.client.NodeList(ctx, opts.ToDockerNodeListOptions())
	}
}
