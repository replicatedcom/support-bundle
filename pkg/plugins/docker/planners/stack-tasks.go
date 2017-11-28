package planners

import (
	"time"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

// StackTasks returns the list of tasks, filtered for a given stack
func (d *Docker) StackTasks(spec types.Spec) []types.Task {
	if spec.DockerSwarmOptions == nil || spec.DockerSwarmOptions.Namespace == "" {
		err := errors.New("spec for docker.stack-tasks requires a namespace")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	scrubber, err := plans.RawScrubber(spec.Config.Scrub)
	if err != nil {
		err = errors.Wrap(err, "create scrubber for docker.stack-tasks")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	services := &plans.StructuredSource{
		Producer:    d.producers.StackTasks(spec.DockerSwarmOptions.Namespace),
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
