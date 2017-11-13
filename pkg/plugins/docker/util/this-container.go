package util

import (
	"context"
	"strings"
	"sync"

	dockertypes "github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/pkg/errors"
)

var (
	thisContainer    *dockertypes.Container
	thisContainerMtx sync.Mutex
)

// find a running container with image containing "support-bundle"
func ThisContainer(ctx context.Context, client *docker.Client) (*dockertypes.Container, error) {
	thisContainerMtx.Lock()
	defer thisContainerMtx.Unlock()

	if thisContainer != nil {
		return thisContainer, nil
	}

	containers, err := client.ContainerList(ctx, dockertypes.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	for _, c := range containers {
		if strings.Contains(c.Image, "support-bundle") {
			thisContainer = &c
			return thisContainer, nil
		}
	}

	return nil, errors.New("not found")
}
