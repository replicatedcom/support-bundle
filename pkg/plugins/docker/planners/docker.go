package planners

import (
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
