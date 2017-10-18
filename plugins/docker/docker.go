package docker

import (
	"github.com/docker/docker/client"
	"github.com/replicatedcom/support-bundle/plugins/docker/planners"
	"github.com/replicatedcom/support-bundle/plugins/docker/producers"
	"github.com/replicatedcom/support-bundle/types"
)

func New() (types.Plugin, error) {
	c, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	producers := producers.New(c)
	docker := planners.New(producers)
	return map[string]types.Planner{
		"daemon": docker.Daemon,
		"logs":   docker.Logs,
	}, nil
}
