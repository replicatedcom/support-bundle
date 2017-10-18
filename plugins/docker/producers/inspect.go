package producers

import (
	"context"

	"github.com/replicatedcom/support-bundle/types"
)

func (d *Docker) Inspect(containerID string) types.StructuredProducer {
	return func(ctx context.Context) (interface{}, error) {
		return d.client.ContainerInspect(ctx, containerID)
	}
}
