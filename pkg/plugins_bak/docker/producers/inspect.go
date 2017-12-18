package producers

import (
	"context"

	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) Inspect(containerID string) types.StructuredProducer {
	return func(ctx context.Context) (interface{}, error) {
		return d.client.ContainerInspect(ctx, containerID)
	}
}

func (d *Docker) InspectName(containerName string) types.StructuredProducer {
	return func(ctx context.Context) (interface{}, error) {
		containerID, err := d.getContainerID(ctx, containerName)
		if err != nil {
			return nil, err
		}
		return d.Inspect(containerID)(ctx)
	}
}
