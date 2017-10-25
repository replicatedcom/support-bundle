package core

import (
	"github.com/replicatedcom/support-bundle/plugins/core/planners"
	"github.com/replicatedcom/support-bundle/types"
)

func New() types.Plugin {
	return map[string]types.Planner{
		"loadavg":      planners.PlanLoadAverage,
		"hostname":     planners.Hostname,
		"uptime":       planners.Uptime,
		"read-file":    planners.ReadFile,
		"read-command": planners.ReadCommand,
	}
}
