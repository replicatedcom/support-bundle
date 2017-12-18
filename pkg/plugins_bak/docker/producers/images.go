package producers

import (
	"context"

	dockertypes "github.com/docker/docker/api/types"
)

func (d *Docker) ImagesAll(ctx context.Context) (interface{}, error) {
	return d.client.ImageList(ctx, dockertypes.ImageListOptions{All: true})
}
