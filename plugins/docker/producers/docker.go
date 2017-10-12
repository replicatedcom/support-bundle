package producers

import (
	docker "github.com/docker/docker/client"
)

type Docker struct {
	client *docker.Client
}

func New(client *docker.Client) *Docker {
	return &Docker{client}
}
