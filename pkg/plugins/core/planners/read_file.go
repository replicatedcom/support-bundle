package planners

import (
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
	jww "github.com/spf13/jwalterweatherman"
)

func (c *Core) ReadFile(spec types.Spec) []types.Task {
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

	var task plans.StreamsSource
	if c.inContainer {
		jww.DEBUG.Println("os.read-file in container")
		task.Producer = c.producers.DockerReadFile(*spec.CoreReadFile)
		task.StreamFormat = plans.StreamFormatTar
	} else {
		task.Producer = c.producers.ReadFile(*spec.CoreReadFile)
	}
	task, err = plans.SetCommonFieldsStreamsSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
