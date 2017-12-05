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
		"daemon":               docker.Daemon,
		"logs":                 docker.Logs,
		"inspect":              docker.Inspect,
		"read-file":            docker.ReadFile,
		"read-host-file":       docker.ReadHostFile,
		"exec-command":         docker.ExecCommand,
		"run-command":          docker.RunCommand,
		"container-ls-logs":    docker.ContainerLogs,
		"container-ls-inspect": docker.ContainerInspect,
		"services":             docker.Services,
		"service-logs":         docker.ServiceLogs,
		"swarm-nodes":          docker.SwarmNodes,
		"swarm-node-tasks":     docker.SwarmNodeTasks,
		"stack-services":       docker.StackServices,
		"stack-services-logs":  docker.StackServiceLogs,
		"stack-tasks":          docker.StackTasks,
		"stack-tasks-logs":     docker.StackTaskLogs,
	}, nil
}
