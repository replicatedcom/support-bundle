package planners

import (
	docker "github.com/docker/docker/client"
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
