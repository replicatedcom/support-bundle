package producers

import (
	"context"

	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) ContainerInspect(containerID string) types.StructuredProducer {
	return func(ctx context.Context) (interface{}, error) {
		return d.client.ContainerInspect(ctx, containerID)
	}
}
