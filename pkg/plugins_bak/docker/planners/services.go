package planners

import (
	"time"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

// Services returns the results of docker service ls
func (d *Docker) Services(spec types.Spec) []types.Task {
	scrubber, err := plans.RawScrubber(spec.Config.Scrub)
	if err != nil {
		err = errors.Wrap(err, "create scrubber for docker.services")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	services := &plans.StructuredSource{
		Producer:    d.producers.Services(),
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
