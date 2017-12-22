package docker

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/replicatedcom/support-bundle/pkg/plugins/docker/planners"
	"github.com/replicatedcom/support-bundle/pkg/plugins/docker/producers"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

type Docker struct {
	planner *planners.Docker
}

func New() (*Docker, error) {
	c, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	c.NegotiateAPIVersion(context.Background())
	producers := producers.New(c)
	return &Docker{
		planner: planners.New(producers, c),
	}, nil
}

func (p *Docker) Plan(spec types.Spec) types.Planner {
	switch {
	case spec.DockerContainerInspect != nil:
		return p.planner.ContainerInspect
	case spec.DockerContainerLs != nil:
		return p.planner.ContainerLs
	case spec.DockerContainerLogs != nil:
		return p.planner.ContainerLogs
	case spec.DockerContainerRun != nil:
		return p.planner.ContainerRun
	case spec.DockerImageLs != nil:
		return p.planner.ImageLs
	case spec.DockerImages != nil:
		return p.planner.ImageLs
	case spec.DockerInfo != nil:
		return p.planner.Info
	case spec.DockerNodeLs != nil:
		return p.planner.NodeLs
	case spec.DockerPs != nil:
		return p.planner.ContainerLs
	case spec.DockerRun != nil:
		return p.planner.ContainerRun
	case spec.DockerServiceLogs != nil:
		return p.planner.ServiceLogs
	case spec.DockerServiceLs != nil:
		return p.planner.ServiceLs
	case spec.DockerServicePs != nil:
		return p.planner.ServicePs
	case spec.DockerStackServiceLogs != nil:
		return p.planner.StackServiceLogs
	case spec.DockerStackServiceLs != nil:
		return p.planner.StackServiceLs
	case spec.DockerStackServicePs != nil:
		return p.planner.StackServicePs
	case spec.DockerStackTaskLogs != nil:
		return p.planner.StackTaskLogs
	case spec.DockerTaskLogs != nil:
		return p.planner.TaskLogs
	case spec.DockerTaskLs != nil:
		return p.planner.TaskLs
	case spec.DockerVersion != nil:
		return p.planner.Version
	default:
		return nil
	}
}
