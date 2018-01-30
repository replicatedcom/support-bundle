package bundle

import (
	"github.com/replicatedcom/support-bundle/pkg/types"
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
		// dont show these to the customer, they're weird and scary
		//jww.ERROR.Printf("Producer not defined for spec %+v", spec)
	}

	return tasks
}
