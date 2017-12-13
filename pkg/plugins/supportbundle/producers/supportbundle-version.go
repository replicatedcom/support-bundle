package producers

import (
	"context"

	"github.com/replicatedcom/support-bundle/pkg/types"
	"github.com/replicatedcom/support-bundle/pkg/version"
)

func Version() types.StructuredProducer {
	return func(ctx context.Context) (interface{}, error) {
		version.Init()
		version := version.GetBuild()

		return version, nil
	}
}
