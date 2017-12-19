package bundle

import (
	"github.com/replicatedcom/support-bundle/pkg/plugins/core"
	"github.com/replicatedcom/support-bundle/pkg/plugins/docker"
	"github.com/replicatedcom/support-bundle/pkg/types"
	jww "github.com/spf13/jwalterweatherman"
)

type Planner struct {
	Core   *core.Core
	Docker *docker.Docker
}

func (p *Planner) Plan(specs []types.Spec) []types.Task {
	var tasks []types.Task

	for _, spec := range specs {
		if p.Core != nil {
			if planner := p.Core.Plan(spec); planner != nil {
				tasks = append(tasks, planner(spec)...)
				continue
			}
		}
		if p.Docker != nil {
			if planner := p.Docker.Plan(spec); planner != nil {
				tasks = append(tasks, planner(spec)...)
				continue
			}
		}
		jww.ERROR.Printf("Producer not defined for spec %v\n", spec)
	}

	return tasks
}
