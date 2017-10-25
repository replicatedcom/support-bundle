package planners

import (
	"errors"

	"github.com/replicatedcom/support-bundle/plans"
	"github.com/replicatedcom/support-bundle/types"
)

func (d *Docker) Logs(spec types.Spec) []types.Task {
	if spec.Config.ContainerID == "" && spec.Config.ContainerName == "" {
		err := errors.New("spec requires a container id or name")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	task := &plans.StreamSource{
		Producer:  d.producers.Logs(spec.Config.ContainerID),
		RawPath:   spec.Raw,
		JSONPath:  spec.JSON,
		HumanPath: spec.Human,
	}

	if spec.Config.ContainerName != "" {
		task = &plans.StreamSource{
			Producer:  d.producers.LogsName(spec.Config.ContainerName),
			RawPath:   spec.Raw,
			JSONPath:  spec.JSON,
			HumanPath: spec.Human,
		}
	}

	return []types.Task{task}
}
