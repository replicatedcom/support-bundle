package docker

import (
	"github.com/docker/docker/client"
	"github.com/replicatedcom/support-bundle/pkg/plugins/docker/planners"
	"github.com/replicatedcom/support-bundle/pkg/plugins/docker/producers"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

type Docker struct {
	planner *planners.Docker
}

func New() (*Docker, error) {
	c, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	producers := producers.New(c)
	return &Docker{
		planner: planners.New(producers, c),
	}, nil
}

func (p *Docker) Plan(spec types.Spec) types.Planner {
	switch {
	case spec.DockerInfo != nil:
		return p.planner.Info
	case spec.DockerVersion != nil:
		return p.planner.Version
	default:
		return nil
	}
}
