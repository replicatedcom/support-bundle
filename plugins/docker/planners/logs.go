package planners

import (
	"errors"

	"github.com/replicatedcom/support-bundle/plans"
	"github.com/replicatedcom/support-bundle/types"
)

func (d *Docker) Logs(spec types.Spec) []types.Task {
	if spec.Config.ContainerID == "" {
		err := errors.New("spec requires a container id")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	task := &plans.StreamSource{
		Producer:  d.producers.Logs(spec.Config.ContainerID),
		RawPath:   spec.Raw,
		JSONPath:  spec.JSON,
		HumanPath: spec.Human,
	}

	return []types.Task{task}
}
