package docker

import (
	docker "github.com/docker/docker/client"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/docker/planners"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/docker/producers"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

type Docker struct {
	planner *planners.Docker
}

func New(client docker.CommonAPIClient) *Docker {
	return &Docker{
		planner: planners.New(
			producers.New(client),
			client,
		),
	}
}

func (p *Docker) Plan(spec types.Spec) types.Planner {
	switch {
	case spec.DockerContainerCp != nil:
		return p.planner.ContainerCp
	case spec.DockerContainerExec != nil:
		return p.planner.ContainerExec
	case spec.DockerContainerInspect != nil:
		return p.planner.ContainerInspect
	case spec.DockerContainerLs != nil:
		return p.planner.ContainerLs
	case spec.DockerContainerLogs != nil:
		return p.planner.ContainerLogs
	case spec.DockerContainerRun != nil:
		return p.planner.ContainerRun
	case spec.DockerExec != nil:
		return p.planner.ContainerExec
	case spec.DockerImageLs != nil:
		return p.planner.ImageLs
	case spec.DockerImages != nil:
		return p.planner.ImageLs
	case spec.DockerInfo != nil:
		return p.planner.Info
	case spec.DockerLogs != nil:
		return p.planner.ContainerLogs
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
