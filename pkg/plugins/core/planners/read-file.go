package planners

import (
	"errors"
	"time"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/plugins/core/producers"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func ReadFile(spec types.Spec) []types.Task {
	if spec.Config.FilePath == "" {
		err := errors.New("spec for core.read-file requires a filename within config")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	task := &plans.ByteSource{
		Producer:  producers.ReadFile(spec.Config.FilePath),
		RawPath:   spec.Raw,
		JSONPath:  spec.JSON,
		HumanPath: spec.Human,
	}

	if spec.TimeoutSeconds != 0 {
		task.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
	}

	return []types.Task{task}
}
