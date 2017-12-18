package docker

import (
	"github.com/docker/docker/client"
	"github.com/replicatedcom/support-bundle/pkg/plugins/docker/planners"
	"github.com/replicatedcom/support-bundle/pkg/plugins/docker/producers"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func New() (types.Plugin, error) {
	c, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	producers := producers.New(c)
	docker := planners.New(producers, c)
	return map[string]types.Planner{
		"info":              docker.Info,
		"version":           docker.Version,
		"ps":                docker.Ps,
		"images":            docker.Images,
		"container-logs":    docker.ContainerLogs,
		"container-inspect": docker.ContainerInspect,
		"container-cp":      docker.ContainerCp,
		// "read-host-file":       docker.ReadHostFile,
		"container-exec": docker.ContainerExec,
		"container-run":  docker.ContainerRun,
		// "container-ls-logs":    docker.ContainerLogs,
		// "container-ls-inspect": docker.ContainerInspect,
		"services":           docker.Services,
		"service-logs":       docker.ServiceLogs,
		"swarm-nodes":        docker.SwarmNodes,
		"swarm-node-tasks":   docker.SwarmNodeTasks,
		"stack-services":     docker.StackServices,
		"stack-service-logs": docker.StackServiceLogs,
		"stack-tasks":        docker.StackTasks,
		"stack-task-logs":    docker.StackTaskLogs,
	}, nil
}
