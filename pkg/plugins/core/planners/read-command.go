package planners

import (
	"time"

	"github.com/pkg/errors"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/plugins/core/producers"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func ReadCommand(spec types.Spec) []types.Task {
	if spec.Config.Command == "" {
		err := errors.New("spec for core.read-command requires a command within config")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	scrubber, err := plans.RawScrubber(spec.Config.Scrub)
	if err != nil {
		err = errors.Wrap(err, "create scrubber for docker.read-file")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	task := &plans.ByteSource{
		Producer:    producers.ReadCommand(spec.Config.Command, spec.Config.Args...),
		RawScrubber: scrubber,
		RawPath:     spec.Raw,
		JSONPath:    spec.JSON,
		HumanPath:   spec.Human,
	}

	if spec.TimeoutSeconds != 0 {
		task.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
	}

	return []types.Task{task}
}
