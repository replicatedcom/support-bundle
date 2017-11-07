package bundle

import (
	"strings"

	"github.com/replicatedcom/support-bundle/pkg/types"
	jww "github.com/spf13/jwalterweatherman"
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
			jww.ERROR.Printf("Plugin %s not defined\n", parts[0])
			continue
		}
		planner, ok := plugin[parts[1]]
		if !ok {
			jww.ERROR.Printf("Planner %s not defined\n", parts[0])
			continue
		}
		tasks = append(tasks, planner(spec)...)
	}

	return tasks
}
