package planners

import (
	"time"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) Inspect(spec types.Spec) []types.Task {
	if spec.Config.ContainerID == "" && spec.Config.ContainerName == "" {
		err := errors.New("spec for docker.inspect requires a container id or name")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	scrubber, err := plans.RawScrubber(spec.Config.Scrub)
	if err != nil {
		err = errors.Wrap(err, "create scrubber for docker.inspect")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	producer := d.producers.Inspect(spec.Config.ContainerID)

	if spec.Config.ContainerName != "" {
		producer = d.producers.InspectName(spec.Config.ContainerName)
	}

	task := &plans.StructuredSource{
		Producer:    producer,
		RawScrubber: scrubber,
		RawPath:     spec.Raw,
		JSONPath:    spec.JSON,
		HumanPath:   spec.Human,
	}

	if spec.TimeoutSeconds != 0 {
		task.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
	}

	return []types.Task{task}
}
