package planners

import (
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/plugins/supportbundle/producers"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func Version(spec types.Spec) []types.Task {
	return []types.Task{
		&plans.StructuredSource{
			Producer:  producers.Version(),
			RawPath:   spec.Raw,
			JSONPath:  spec.JSON,
			HumanPath: spec.Human,
		},
	}
}
