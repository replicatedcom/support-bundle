package core

import (
	"github.com/replicatedcom/support-bundle/pkg/plugins/core/planners"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func Plan(spec types.Spec) types.Planner {
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
