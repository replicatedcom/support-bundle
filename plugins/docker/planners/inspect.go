package planners

import (
	"errors"
	"time"

	"github.com/replicatedcom/support-bundle/plans"
	"github.com/replicatedcom/support-bundle/types"
)

func (d *Docker) Inspect(spec types.Spec) []types.Task {
	if spec.Config.ContainerID == "" && spec.Config.ContainerName == "" {
		err := errors.New("spec requires a container id or name")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	producer := d.producers.Inspect(spec.Config.ContainerID)

	if spec.Config.ContainerName != "" {
		producer = d.producers.InspectName(spec.Config.ContainerName)
	}

	task := &plans.StructuredSource{
		Producer:  producer,
		RawPath:   spec.Raw,
		JSONPath:  spec.JSON,
		HumanPath: spec.Human,
	}

	if spec.TimeoutSeconds != 0 {
		task.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
	}

	return []types.Task{task}
}
