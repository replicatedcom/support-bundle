package producers

import (
	"context"

	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (k *Kubernetes) Version() types.StructuredProducer {
	return func(ctx context.Context) (interface{}, error) {
		return k.client.ServerVersion()
	}
}
