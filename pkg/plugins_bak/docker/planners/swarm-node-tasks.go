package planners

import (
	"time"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

// SwarmNodeTasks returns the tasks running on each node in the swarm
func (d *Docker) SwarmNodeTasks(spec types.Spec) []types.Task {
	scrubber, err := plans.RawScrubber(spec.Config.Scrub)
	if err != nil {
		err = errors.Wrap(err, "create scrubber for docker.swarm-node-tasks")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	swarmNodeTasks := &plans.StreamsSource{
		Producer:    d.producers.SwarmNodeTasks(),
		RawScrubber: scrubber,
		RawPath:     spec.Raw,
		JSONPath:    spec.JSON,
		HumanPath:   spec.Human,
	}

	if spec.TimeoutSeconds != 0 {
		swarmNodeTasks.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
	}

	return []types.Task{
		swarmNodeTasks,
	}
}
