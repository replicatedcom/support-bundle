package planners

import (
	"errors"

	"github.com/replicatedcom/support-bundle/plans"
	"github.com/replicatedcom/support-bundle/types"
)

func (d *Docker) ExecCommand(spec types.Spec) []types.Task {
	// fullCommand := append([]string{spec.Config.Command}, spec.Config.Args...)

	if spec.Config.ContainerID == "" && spec.Config.ContainerName == "" {
		err := errors.New("spec for docker.run-command requires a container ID or Name within config")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	} else if spec.Config.Command == "" {
		err := errors.New("spec for docker.run-command requires a command within config")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	// task := &plans.StreamsSource{
	// 	Producer:  d.producers.ExecCommand(spec.Config.ContainerID, fullCommand),
	// 	RawPath:   spec.Raw,
	// 	JSONPath:  spec.JSON,
	// 	HumanPath: spec.Human,
	// }

	// if spec.Config.ContainerName != "" {
	// 	task = &plans.StreamsSource{
	// 		Producer:  d.producers.ExecCommandByName(spec.Config.ContainerName),
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
