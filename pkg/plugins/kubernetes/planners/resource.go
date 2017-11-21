package planners

import (
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (k *Kubernetes) Resource(spec types.Spec) []types.Task {
	resourceType := spec.Config.ResourceName
	return []types.Task{
		&plans.StructuredSource{
			Producer:  k.p.Resource(resourceType),
			RawPath:   spec.Raw,
			JSONPath:  spec.JSON,
			HumanPath: spec.Human,
		},
	}
}
