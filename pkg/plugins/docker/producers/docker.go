package producers

import (
	"context"

	docker "github.com/docker/docker/client"
)

type Docker struct {
	client *docker.Client
}

func New(client *docker.Client) *Docker {
	client.NegotiateAPIVersion(context.Background())
	return &Docker{client}
}
