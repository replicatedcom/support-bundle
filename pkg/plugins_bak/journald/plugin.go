package journald

import (
	"github.com/docker/docker/client"
	dockerproducers "github.com/replicatedcom/support-bundle/pkg/plugins/docker/producers"
	"github.com/replicatedcom/support-bundle/pkg/plugins/journald/planners"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func New() (types.Plugin, error) {
	c, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	producers := dockerproducers.New(c)

	journald := planners.New(producers, c)

	return map[string]types.Planner{
		"logs": journald.Logs,
	}, nil
}
