package planners

import (
	"github.com/pkg/errors"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/plugins/core/producers"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func RunCommand(spec types.Spec) []types.Task {
	if spec.CoreRunCommand == nil {
		err := errors.New("spec for os.run-command options required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	if spec.CoreRunCommand.Name == "" {
		err := errors.New("spec for os.run-command name required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	task := plans.StreamsSource{
		Spec:     spec,
		Producer: producers.RunCommand(*spec.CoreRunCommand),
	}
	var err error
	task, err = plans.SetCommonFieldsStreamsSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
