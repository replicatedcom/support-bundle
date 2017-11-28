package planners

import (
	"time"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

// StackServiceLogs returns the logs of services within the given namespace
func (d *Docker) StackServiceLogs(spec types.Spec) []types.Task {
	if spec.DockerSwarmOptions == nil || spec.DockerSwarmOptions.Namespace == "" {
		err := errors.New("spec for docker.stack-service-logs requires a namespace")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	scrubber, err := plans.RawScrubber(spec.Config.Scrub)
	if err != nil {
		err = errors.Wrap(err, "create scrubber for docker.stack-service-logs")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	services := &plans.StreamsSource{
		Producer:    d.producers.StackServiceLogs(spec.DockerSwarmOptions.Namespace),
		RawScrubber: scrubber,
		RawPath:     spec.Raw,
		JSONPath:    spec.JSON,
		HumanPath:   spec.Human,
	}

	if spec.TimeoutSeconds != 0 {
		services.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
	}

	return []types.Task{
		services,
	}
}
