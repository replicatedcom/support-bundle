package planners

import (
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (k *Kubernetes) Resource(spec types.Spec) []types.Task {
	resource := spec.KubernetesResourceCommand
	return []types.Task{
		&plans.StructuredSource{
			Producer:  k.p.Resource(resource.Type, resource.Namespace),
			RawPath:   spec.Raw,
			JSONPath:  spec.JSON,
			HumanPath: spec.Human,
		},
	}
}
