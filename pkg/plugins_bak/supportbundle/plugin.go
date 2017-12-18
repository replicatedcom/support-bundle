package supportbundle

import (
	"github.com/replicatedcom/support-bundle/pkg/plugins/supportbundle/planners"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func New() types.Plugin {
	return map[string]types.Planner{
		"version": planners.Version,
	}
}
