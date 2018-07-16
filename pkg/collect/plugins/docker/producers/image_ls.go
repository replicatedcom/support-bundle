package producers

import (
	"context"

	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func (d *Docker) ImageLs(opts types.DockerImageLsOptions) types.StructuredProducer {
	return func(ctx context.Context) (interface{}, error) {
		return d.client.ImageList(ctx, opts.ToDockerImageListOptions())
	}
}
