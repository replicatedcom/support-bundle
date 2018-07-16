package planners

import (
	docker "github.com/docker/docker/client"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/journald/producers"
)

type Journald struct {
	producers    *producers.Journald
	dockerClient docker.CommonAPIClient
	inContainer  bool
}

func New(journald *producers.Journald, dockerClient docker.CommonAPIClient, inContainer bool) *Journald {
	return &Journald{
		producers:    journald,
		dockerClient: dockerClient,
		inContainer:  inContainer,
	}
}
