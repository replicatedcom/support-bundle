package producers

import (
	"context"

	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) ServiceLs(opts types.DockerServiceLsOptions) types.StructuredProducer {
	return func(ctx context.Context) (interface{}, error) {
		return d.client.ServiceList(ctx, opts.ToDockerServiceListOptions())
	}
}
