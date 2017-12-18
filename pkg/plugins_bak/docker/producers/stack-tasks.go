package producers

import (
	"context"

	"github.com/replicatedcom/support-bundle/pkg/types"

	dockertypes "github.com/docker/docker/api/types"
)

func (d *Docker) StackTasks(namespace string) types.StructuredProducer {
	return func(ctx context.Context) (interface{}, error) {
		return d.client.TaskList(ctx, dockertypes.TaskListOptions{Filters: dockerNamespaceFilter(namespace)})
	}
}
