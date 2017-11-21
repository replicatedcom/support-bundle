package planners

import (
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (k *Kubernetes) Logs(spec types.Spec) []types.Task {
	opts := spec.KubernetesContainerLogsOptions
	if spec.Config.PodName == "" {
		err := errors.New("spec for kubernettes.logs requires a pod name (and a container name for multi-container pods")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	scrubber, err := plans.RawScrubber(opts.Scrub)
	if err != nil {
		err = errors.Wrap(err, "create scrubber for kubernetes.logs")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	return []types.Task{
		&plans.StreamSource{
			Producer:    k.p.Logs(opts.PodName, opts.ContainerName, opts.Namespace),
			RawScrubber: scrubber,
			RawPath:     spec.Raw,
			JSONPath:    spec.JSON,
			HumanPath:   spec.Human,
		},
	}
}
