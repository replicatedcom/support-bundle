package planners

import (
	"errors"

	"github.com/replicatedcom/support-bundle/plans"
	"github.com/replicatedcom/support-bundle/types"
)

func (d *Docker) Inspect(spec types.Spec) []types.Task {
	config := parseContainerConfig(spec.Config)
	if config.ContainerID == "" {
		err := errors.New("spec requires container config")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	task := &plans.StructuredSource{
		Producer:  d.producers.Inspect(config.ContainerID),
		RawPath:   spec.Raw,
		JSONPath:  spec.JSON,
		HumanPath: spec.Human,
	}

	return []types.Task{task}
}
