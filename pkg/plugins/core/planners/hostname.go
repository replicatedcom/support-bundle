package planners

import (
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/plugins/core/producers"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func Hostname(spec types.Spec) []types.Task {
	if spec.CoreHostname == nil {
		err := errors.New("spec for os.hostname options required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	task := plans.StreamsSource{
		Producer: producers.RunCommand(types.CoreRunCommandOptions{Name: "hostname"}),
	}
	var err error
	task, err = plans.SetCommonFieldsStreamsSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
