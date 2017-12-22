package producers

import (
	"context"

	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) ServicePs(opts types.DockerServicePsOptions) types.StructuredProducer {
	return func(ctx context.Context) (interface{}, error) {
		return d.client.TaskList(ctx, opts.ToDockerTaskListOptions())
	}
}
