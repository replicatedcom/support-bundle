package core

import (
	"github.com/replicatedcom/support-bundle/pkg/plugins/core/planners"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

type Core struct {
}

func New() (*Core, error) {
	return &Core{}, nil
}

func (p *Core) Plan(spec types.Spec) types.Planner {
	switch {
	case spec.CoreHostname != nil:
		return planners.Hostname
	case spec.CoreHTTPRequest != nil:
		return planners.HTTPRequest
	case spec.CoreLoadavg != nil:
		return planners.Loadavg
	case spec.CoreReadFile != nil:
		return planners.ReadFile
	case spec.CoreRunCommand != nil:
		return planners.RunCommand
	case spec.CoreUptime != nil:
		return planners.Uptime
	default:
		return nil
	}
}
