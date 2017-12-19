package producers

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
)

type Docker struct {
	client *docker.Client
}

func New(client *docker.Client) *Docker {
	client.NegotiateAPIVersion(context.Background())
	return &Docker{client}
}

// get the ID of a container given the name
func (d *Docker) getContainerID(ctx context.Context, containerName string) (string, error) {
	matchingID := ""

	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return matchingID, err
	}

	for _, container := range containers {
		for _, name := range container.Names {
			if name == containerName || name == "/"+containerName {
				matchingID = container.ID
				return matchingID, nil
			}
		}
	}

	err = errors.New("unable to find container with name " + containerName)
	return matchingID, err
}
