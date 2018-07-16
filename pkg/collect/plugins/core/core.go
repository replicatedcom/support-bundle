package core

import (
	"context"
	"os"

	docker "github.com/docker/docker/client"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/core/planners"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/core/producers"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

type Core struct {
	planner *planners.Core
}

func New() (*Core, error) {
	inContainer := os.Getenv("IN_CONTAINER") != ""
	var dockerClient docker.CommonAPIClient
	if inContainer {
		c, err := docker.NewEnvClient()
		if err != nil {
			return nil, err
		}
		c.NegotiateAPIVersion(context.Background())
		dockerClient = c
	}

	producers := producers.New(dockerClient)
	return &Core{
		planner: planners.New(producers, dockerClient, inContainer),
	}, nil
}

func (p *Core) Plan(spec types.Spec) types.Planner {
	switch {
	case spec.CoreHostname != nil:
		return p.planner.Hostname
	case spec.CoreHTTPRequest != nil:
		return p.planner.HTTPRequest
	case spec.CoreLoadavg != nil:
		return p.planner.Loadavg
	case spec.CoreReadFile != nil:
		return p.planner.ReadFile
	case spec.CoreRunCommand != nil:
		return p.planner.RunCommand
	case spec.CoreUptime != nil:
		return p.planner.Uptime
	default:
		return nil
	}
}
