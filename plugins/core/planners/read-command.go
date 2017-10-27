package planners

import (
	"errors"
	"time"

	"github.com/replicatedcom/support-bundle/plans"
	"github.com/replicatedcom/support-bundle/plugins/core/producers"
	"github.com/replicatedcom/support-bundle/types"
)

func ReadCommand(spec types.Spec) []types.Task {
	if spec.Config.Command == "" {
		err := errors.New("spec for core.read-command requires a command within config")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	task := &plans.ByteSource{
		Producer:  producers.ReadCommand(spec.Config.Command, spec.Config.Args...),
		RawPath:   spec.Raw,
		JSONPath:  spec.JSON,
		HumanPath: spec.Human,
	}

	if spec.TimeoutSeconds != 0 {
		task.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
	}

	return []types.Task{task}
}
