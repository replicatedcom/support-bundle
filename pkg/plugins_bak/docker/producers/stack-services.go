package producers

import (
	"context"

	"github.com/replicatedcom/support-bundle/pkg/types"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

func dockerNamespaceFilter(namespace string) filters.Args {
	filter := filters.NewArgs()
	filter.Add("label", "com.docker.stack.namespace="+namespace)
	return filter
}

func (d *Docker) StackServices(namespace string) types.StructuredProducer {
	return func(ctx context.Context) (interface{}, error) {
		return d.client.ServiceList(ctx, dockertypes.ServiceListOptions{Filters: dockerNamespaceFilter(namespace)})
	}
}
