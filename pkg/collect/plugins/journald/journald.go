package journald

import (
	docker "github.com/docker/docker/client"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/journald/planners"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/journald/producers"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

type Journald struct {
	planner *planners.Journald
}

func New(inContainer bool, dockerClient docker.CommonAPIClient) *Journald {
	return &Journald{
		planner: planners.New(
			producers.New(dockerClient),
			dockerClient,
			inContainer,
		),
	}
}

func (p *Journald) Plan(spec types.Spec) types.Planner {
	switch {
	case spec.JournaldLogs != nil:
		return p.planner.Logs
	default:
		return nil
	}
}
