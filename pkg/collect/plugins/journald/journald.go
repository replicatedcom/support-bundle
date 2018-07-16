package journald

import (
	"context"
	"os"

	docker "github.com/docker/docker/client"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/journald/planners"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/journald/producers"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

type Journald struct {
	planner *planners.Journald
}

func New() (*Journald, error) {
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

	dockerClient, err := docker.NewEnvClient()
	if err != nil {
		return nil, err
	}
	dockerClient.NegotiateAPIVersion(context.Background())

	producers := producers.New(dockerClient)
	return &Journald{
		planner: planners.New(producers, dockerClient, inContainer),
	}, nil
}

func (p *Journald) Plan(spec types.Spec) types.Planner {
	switch {
	case spec.JournaldLogs != nil:
		return p.planner.Logs
	default:
		return nil
	}
}
