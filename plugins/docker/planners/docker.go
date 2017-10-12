package planners

import (
	"github.com/replicatedcom/support-bundle/plugins/docker/producers"
)

type Docker struct {
	producers *producers.Docker
}

func New(docker *producers.Docker) *Docker {
	return &Docker{
		producers: docker,
	}
}
