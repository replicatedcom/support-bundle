package producers

import (
	"context"

	dockertypes "github.com/docker/docker/api/types"
)

func (d *Docker) PSAll(ctx context.Context) (interface{}, error) {
	return d.client.ContainerList(ctx, dockertypes.ContainerListOptions{All: true})
}
