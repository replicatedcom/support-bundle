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
	docker := planners.New(producers)
	return map[string]types.Planner{
		"daemon":       docker.Daemon,
		"logs":         docker.Logs,
		"inspect":      docker.Inspect,
		"read-file":    docker.ReadFile,
		"exec-command": docker.ExecCommand,
	}, nil
}
