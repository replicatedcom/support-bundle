package bundle

import (
	"github.com/replicatedcom/support-bundle/pkg/collect/types"

	jww "github.com/spf13/jwalterweatherman"
)

type Planner struct {
	Plugins []types.Plugin
}

func (p *Planner) AddPlugin(plugin types.Plugin) {
	p.Plugins = append(p.Plugins, plugin)
}

func (p *Planner) Plan(specs []types.Spec) []types.Task {
	var tasks []types.Task

Loop:
	for _, spec := range specs {
		for _, plugin := range p.Plugins {
			if planner := plugin.Plan(spec); planner != nil {
				tasks = append(tasks, planner(spec)...)
				continue Loop
			}
		}
		// don't show these to the customer, they're weird and scary
		jww.DEBUG.Printf("Producer not defined for spec %+v", spec)
	}

	return tasks
}
