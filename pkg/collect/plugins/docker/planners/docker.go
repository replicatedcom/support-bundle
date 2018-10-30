package planners

import (
	"context"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	docker "github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/docker/producers"
)

type Docker struct {
	producers *producers.Docker
	client    docker.CommonAPIClient
}

func New(docker *producers.Docker, client docker.CommonAPIClient) *Docker {
	return &Docker{
		producers: docker,
		client:    client,
	}
}

func (d *Docker) SelectContainerIDs(name string, labels []string) ([]string, error) {
	if name != "" {
		return []string{name}, nil
	}

	labelFilters := dockertypes.ContainerListOptions{}
	labelFilters.Filters = filters.NewArgs()
	for _, val := range labels {
		labelFilters.Filters.Add("label", val)
	}

	labelledContainers, err := d.client.ContainerList(context.Background(), labelFilters)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to find containers to read file from")
	}

	ids := make([]string, 0)
	for _, container := range labelledContainers {
		ids = append(ids, container.ID)
	}

	return ids, nil
}
