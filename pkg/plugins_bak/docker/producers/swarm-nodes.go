package producers

import (
	"context"

	"github.com/replicatedcom/support-bundle/pkg/types"

	dockertypes "github.com/docker/docker/api/types"
)

func (d *Docker) SwarmNodes() types.StructuredProducer {
	return func(ctx context.Context) (interface{}, error) {
		return d.client.NodeList(ctx, dockertypes.NodeListOptions{})
	}
}
