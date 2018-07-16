package planners

import (
	docker "github.com/docker/docker/client"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/core/producers"
)

type Core struct {
	producers    *producers.Core
	dockerClient docker.CommonAPIClient
	inContainer  bool
}

func New(core *producers.Core, client docker.CommonAPIClient, inContainer bool) *Core {
	return &Core{
		producers:    core,
		dockerClient: client,
		inContainer:  inContainer,
	}
}
