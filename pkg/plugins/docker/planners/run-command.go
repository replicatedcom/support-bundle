package planners

import (
	"errors"
	"time"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) RunCommand(spec types.Spec) []types.Task {
	var err error
	if spec.DockerRunCommand == nil {
		err = errors.New("spec for docker.run-command required")
	} else if spec.DockerRunCommand.ContainerCreateConfig.Config == nil {
		err = errors.New("spec for docker.run-command.ContainerCreateConfig.Config required")
	} else if spec.DockerRunCommand.ContainerCreateConfig.Config.Image == "" {
		err = errors.New("spec for docker.run-command.ContainerCreateConfig.Config.Image required")
	}
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	task := &plans.StreamsSource{
		Producer:  d.producers.RunCommand(*spec.DockerRunCommand),
		RawPath:   spec.Raw,
		JSONPath:  spec.JSON,
		HumanPath: spec.Human,
	}

	if spec.TimeoutSeconds != 0 {
		task.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
	}

	return []types.Task{task}
}
