package producers

import (
	"context"

	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func (d *Docker) TaskLs(opts types.DockerTaskLsOptions) types.StructuredProducer {
	return func(ctx context.Context) (interface{}, error) {
		return d.client.TaskList(ctx, opts.ToDockerTaskListOptions())
	}
}
