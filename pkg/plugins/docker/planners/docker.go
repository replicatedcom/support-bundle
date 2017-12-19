package planners

import (
	"context"
	"fmt"

	dockertypes "github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	"github.com/replicatedcom/support-bundle/pkg/plugins/docker/producers"
)

type Docker struct {
	producers *producers.Docker
	client    *dockerclient.Client
}

func New(docker *producers.Docker, client *dockerclient.Client) *Docker {
	return &Docker{
		producers: docker,
		client:    client,
	}
}

// get the ID of a container given the name
func (d *Docker) getContainerID(ctx context.Context, containerName string) (string, error) {
	containers, err := d.client.ContainerList(ctx, dockertypes.ContainerListOptions{All: true})
	if err != nil {
		return "", err
	}

	for _, container := range containers {
		for _, name := range container.Names {
			if name == containerName || name == "/"+containerName {
				return container.ID, nil
			}
		}
	}

	return "", fmt.Errorf("unable to find container with name %s", containerName)
}
