package planners

import (
	"errors"

	"github.com/replicatedcom/support-bundle/plans"
	"github.com/replicatedcom/support-bundle/types"
)

func (d *Docker) RunCommand(spec types.Spec) []types.Task {
	fullCommand := append([]string{spec.Config.Command}, spec.Config.Args...)

	if (spec.Config.ContainerID == "" && spec.Config.ContainerName == "") || len(fullCommand) == 0 || spec.Config.Command == "" {
		err := errors.New("spec requires a container ID or Name and command within config")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	// task := &plans.StreamsSource{
	// 	Producer:  d.producers.RunCommand(spec.Config.ContainerID, fullCommand),
	// 	RawPath:   spec.Raw,
	// 	JSONPath:  spec.JSON,
	// 	HumanPath: spec.Human,
	// }

	// if spec.Config.ContainerName != "" {
	// 	task = &plans.StreamsSource{
	// 		Producer:  d.producers.RunCommandByName(spec.Config.ContainerName),
	// 		RawPath:   spec.Raw,
	// 		JSONPath:  spec.JSON,
	// 		HumanPath: spec.Human,
	// 	}
	// }

	// if spec.TimeoutSeconds != 0 {
	// 	task.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
	// }

	err := errors.New("This task type not yet implemented")
	task := plans.PreparedError(err, spec)

	return []types.Task{task}
}
