package planners

import (
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (k *Kubernetes) Logs(spec types.Spec) []types.Task {
	return []types.Task{
		&plans.StreamSource{
			Producer:  k.p.Logs(),
			RawPath:   spec.Raw,
			JSONPath:  spec.JSON,
			HumanPath: spec.Human,
		},
	}
}
