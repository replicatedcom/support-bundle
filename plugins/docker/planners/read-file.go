package planners

import (
	"errors"

	"github.com/replicatedcom/support-bundle/plans"
	"github.com/replicatedcom/support-bundle/types"
)

func (d *Docker) ReadFile(spec types.Spec) []types.Task {
	if spec.Config.ContainerID == "" || spec.Config.FilePath == "" {
		err := errors.New("spec requires a container ID and filename within config")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	task := &plans.StreamSource{
		Producer:  d.producers.ReadFile(spec.Config.ContainerID, spec.Config.FilePath),
		RawPath:   spec.Raw,
		JSONPath:  spec.JSON,
		HumanPath: spec.Human,
	}

	return []types.Task{task}
}
