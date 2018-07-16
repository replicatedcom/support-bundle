package producers

import (
	"context"

	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func (d *Docker) ContainerInspect(container string) types.StructuredProducer {
	return func(ctx context.Context) (interface{}, error) {
		return d.client.ContainerInspect(ctx, container)
	}
}
