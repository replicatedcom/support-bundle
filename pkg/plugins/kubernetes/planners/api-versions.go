package planners

import (
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (k *Kubernetes) APIVersions(spec types.Spec) []types.Task {
	return []types.Task{
		&plans.StructuredSource{
			Producer:  k.p.APIVersions(),
			RawPath:   spec.Raw,
			JSONPath:  spec.JSON,
			HumanPath: spec.Human,
		},
	}
}
