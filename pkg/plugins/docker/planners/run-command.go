package planners

import (
	"errors"
	"time"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) RunCommand(spec types.Spec) []types.Task {
	fullCommand := append([]string{spec.Config.Command}, spec.Config.Args...)

	if spec.Config.Image == "" {
		err := errors.New("spec for docker.run-command requires an image within config")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	task := &plans.StreamsSource{
		Producer:  d.producers.RunCommand(spec.Config.Image, fullCommand, spec.Config.Env, spec.Config.Binds),
		RawPath:   spec.Raw,
		JSONPath:  spec.JSON,
		HumanPath: spec.Human,
	}

	if spec.TimeoutSeconds != 0 {
		task.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
	}

	return []types.Task{task}
}
