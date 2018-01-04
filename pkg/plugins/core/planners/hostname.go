package planners

import (
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (c *Core) Hostname(spec types.Spec) []types.Task {
	if spec.CoreHostname == nil {
		err := errors.New("spec for os.hostname required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	task := plans.StreamsSource{
		Producer: c.producers.RunCommand(types.CoreRunCommandOptions{Name: "hostname"}),
	}
	var err error
	task, err = plans.SetCommonFieldsStreamsSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
