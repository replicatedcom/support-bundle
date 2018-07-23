package producers

import (
	docker "github.com/docker/docker/client"
)

type Docker struct {
	client docker.CommonAPIClient
}

func New(dockerClient docker.CommonAPIClient) *Docker {
	return &Docker{dockerClient}
}
