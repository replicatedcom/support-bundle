package producers

import (
	docker "github.com/docker/docker/client"
)

type Docker struct {
	client docker.CommonAPIClient
}

func New(client docker.CommonAPIClient) *Docker {
	return &Docker{client}
}
