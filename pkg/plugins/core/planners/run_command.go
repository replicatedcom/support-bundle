package planners

import (
	"github.com/pkg/errors"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (c *Core) RunCommand(spec types.Spec) []types.Task {
	var err error
	if spec.CoreRunCommand == nil {
		err = errors.New("spec for os.run-command required")
	} else if spec.CoreRunCommand.Name == "" {
		err = errors.New("spec for os.run-command name required")
	}
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	task := plans.StreamsSource{
		Producer: c.producers.RunCommand(*spec.CoreRunCommand),
	}
	task, err = plans.SetCommonFieldsStreamsSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
