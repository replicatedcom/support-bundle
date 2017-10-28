package planners

import (
	"time"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/plugins/docker/producers"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) RunCommand(spec types.Spec) []types.Task {
	opts, err := producers.RunCommandOptionsFromSpec(spec.Config)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	task := &plans.StreamsSource{
		Producer:  d.producers.RunCommand(opts),
		RawPath:   spec.Raw,
		JSONPath:  spec.JSON,
		HumanPath: spec.Human,
	}

	if spec.TimeoutSeconds != 0 {
		task.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
	}

	return []types.Task{task}
}
