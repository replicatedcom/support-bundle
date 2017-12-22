package planners

import (
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/plugins/core/producers"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func ReadFile(spec types.Spec) []types.Task {
	var err error
	if spec.CoreReadFile == nil {
		err = errors.New("spec for os.read-file options required")
	} else if spec.CoreReadFile.Filepath == "" {
		err = errors.New("spec for os.read-file filepath required")
	}
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	task := plans.StreamsSource{
		Producer: producers.ReadFile(*spec.CoreReadFile),
	}
	task, err = plans.SetCommonFieldsStreamsSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
