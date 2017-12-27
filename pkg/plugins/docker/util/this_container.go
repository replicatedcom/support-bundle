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
	thisContainer   *dockertypes.Container
	thisContainerMu sync.Mutex
)

// ThisContainer will find a running container with image containing "support-bundle"
func ThisContainer(ctx context.Context, client docker.CommonAPIClient) (*dockertypes.Container, error) {
	thisContainerMu.Lock()
	defer thisContainerMu.Unlock()

	if thisContainer != nil {
		return thisContainer, nil
	}

	containers, err := client.ContainerList(ctx, dockertypes.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	for _, container := range containers {
		b, _ := json.Marshal(container)
		fmt.Println(string(b))
		// TODO: LABEL "com.replicated.support-bundle"="true"
		if strings.Contains(container.Image, "support-bundle") {
			thisContainer = &container
			return thisContainer, nil
		}
	}

	return nil, errors.New("not found")
}
