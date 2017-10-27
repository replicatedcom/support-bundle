package bundle

import (
	"strings"

	"github.com/replicatedcom/support-bundle/pkg/types"
)

type Planner struct {
	Plugins map[string]types.Plugin
}

func (p Planner) Plan(specs []types.Spec) []types.Task {
	var tasks []types.Task

	for _, spec := range specs {
		parts := strings.Split(spec.Builtin, ".")
		if len(parts) != 2 {
			continue
		}
		plugin, ok := p.Plugins[parts[0]]
		if !ok {
			continue
		}
		planner, ok := plugin[parts[1]]
		if !ok {
			continue
		}
		tasks = append(tasks, planner(spec)...)
	}

	return tasks
}
