package util

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	dockertypes "github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/pkg/errors"
)

var (
	thisContainer   *dockertypes.ContainerJSON
	thisContainerMu sync.Mutex
)

// ThisContainer will find a running container with image containing "support-bundle"
func ThisContainer(ctx context.Context, client docker.CommonAPIClient) (*dockertypes.ContainerJSON, error) {
	thisContainerMu.Lock()
	defer thisContainerMu.Unlock()

	if thisContainer != nil {
		return thisContainer, nil
	}

	containers, err := client.ContainerList(ctx, dockertypes.ContainerListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "container list")
	}
	for _, container := range containers {
		b, _ := json.Marshal(container)
		fmt.Println(string(b))
		// TODO: LABEL "com.replicated.support-bundle"="true"
		if strings.Contains(container.Image, "support-bundle") {
			c, err := client.ContainerInspect(ctx, container.ID)
			if err != nil {
				return nil, errors.Wrapf(err, "container inspect %s", container.ID)
			}
			thisContainer = &c
			return thisContainer, nil
		}
	}

	return nil, errors.New("not found")
}
