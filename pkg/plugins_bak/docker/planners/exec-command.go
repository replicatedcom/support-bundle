package planners

import (
	"time"

	"github.com/pkg/errors"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) ExecCommand(spec types.Spec) []types.Task {
	if spec.Config.ContainerID == "" && spec.Config.ContainerName == "" {
		err := errors.New("spec for docker.exec-command requires a container ID or Name within config")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	} else if spec.Config.Command == "" {
		err := errors.New("spec for docker.exec-command requires a command within config")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	scrubber, err := plans.RawScrubber(spec.Config.Scrub)
	if err != nil {
		err = errors.Wrap(err, "create scrubber for docker.exec-command")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	fullCommand := append([]string{spec.Config.Command}, spec.Config.Args...)

	task := &plans.StreamsSource{
		Producer:    d.producers.ExecCommand(spec.Config.ContainerID, fullCommand),
		RawScrubber: scrubber,
		RawPath:     spec.Raw,
		JSONPath:    spec.JSON,
		HumanPath:   spec.Human,
	}

	if spec.Config.ContainerName != "" {
		task.Producer = d.producers.ExecCommandByName(spec.Config.ContainerName, fullCommand)
	}

	if spec.TimeoutSeconds != 0 {
		task.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
	}

	return []types.Task{task}
}