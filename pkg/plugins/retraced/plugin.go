package retraced

import (
	"github.com/replicatedcom/support-bundle/pkg/plugins/retraced/planners"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func New() types.Plugin {
	return map[string]types.Planner{
		"events": planners.Events,
	}
}
