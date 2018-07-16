package producers

import (
	docker "github.com/docker/docker/client"
)

type Journald struct {
	dockerClient docker.CommonAPIClient
}

func New(dockerClient docker.CommonAPIClient) *Journald {
	return &Journald{dockerClient}
}
