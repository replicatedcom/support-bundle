package planners

import (
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (k *Kubernetes) Logs(spec types.Spec) []types.Task {
	opts := spec.KubernetesContainerLogsOptions
	return []types.Task{
		&plans.StreamSource{
			Producer:  k.p.Logs(opts.PodName, opts.ContainerName, opts.Namespace),
			RawPath:   spec.Raw,
			JSONPath:  spec.JSON,
			HumanPath: spec.Human,
		},
	}
}
