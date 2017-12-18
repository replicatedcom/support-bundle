package bundle

import (
	"github.com/replicatedcom/support-bundle/pkg/plugins/core"
	"github.com/replicatedcom/support-bundle/pkg/types"
	jww "github.com/spf13/jwalterweatherman"
)

type Planner struct {
	Plugins map[string]types.Plugin
}

func (p Planner) Plan(specs []types.Spec) []types.Task {
	var tasks []types.Task

	for _, spec := range specs {
		if planner := core.Plan(spec); planner != nil {
			tasks = append(tasks, planner(spec)...)
		} else {
			jww.ERROR.Println("Producer not defined")
		}
	}

	return tasks
}
