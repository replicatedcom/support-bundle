package util

import (
	"context"
	"os"
	"strings"

	dockertypes "github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/pkg/errors"
)

func ThisContainer(ctx context.Context, client *docker.Client) (*dockertypes.Container, error) {
	thisCmd := strings.Join(os.Args, " ")

	containers, err := client.ContainerList(ctx, dockertypes.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	for _, c := range containers {
		if c.Command == thisCmd {
			return &c, nil
		}
	}
	return nil, errors.New("not found")
}
