package planners

import (
	"time"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

// SwarmNodes returns the list of nodes in the swarm
func (d *Docker) SwarmNodes(spec types.Spec) []types.Task {
	scrubber, err := plans.RawScrubber(spec.Config.Scrub)
	if err != nil {
		err = errors.Wrap(err, "create scrubber for docker.swarm-nodes")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	swarmNodes := &plans.StructuredSource{
		Producer:    d.producers.SwarmNodes(),
		RawScrubber: scrubber,
		RawPath:     spec.Raw,
		JSONPath:    spec.JSON,
		HumanPath:   spec.Human,
	}

	if spec.TimeoutSeconds != 0 {
		swarmNodes.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
	}

	return []types.Task{
		swarmNodes,
	}
}
