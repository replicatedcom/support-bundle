package planners

import (
	docker "github.com/docker/docker/client"
	dockerproducers "github.com/replicatedcom/support-bundle/pkg/plugins/docker/producers"
)

type Journald struct {
	dockerProducers *dockerproducers.Docker
	dockerClient    *docker.Client
}

func New(producers *dockerproducers.Docker, client *docker.Client) *Journald {
	return &Journald{
		dockerProducers: producers,
		dockerClient:    client,
	}
}
