package producers

import (
	"context"

	"github.com/replicatedcom/support-bundle/pkg/types"
)

func ProducerOf(thing interface{}) types.StructuredProducer {
	return func(ctx context.Context) (interface{}, error) {
		return thing, nil
	}
}

func (s *SupportBundle) CustomerMeta(options *types.CustomerMetaOptions) types.StructuredProducer {
	return ProducerOf(options)
}
