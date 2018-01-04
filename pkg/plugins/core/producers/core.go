package producers

import (
	docker "github.com/docker/docker/client"
)

type Core struct {
	dockerClient docker.CommonAPIClient
}

func New(dockerClient docker.CommonAPIClient) *Core {
	return &Core{dockerClient}
}
