package planners

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) ContainerExec(spec types.Spec) []types.Task {
	if spec.DockerContainerExec == nil {
		if spec.DockerExec != nil {
			jww.DEBUG.Println("spec for docker.container-exec aliased from docker.exec")
			spec.DockerContainerExec = spec.DockerExec
		} else {
			err := errors.New("spec for docker.container-exec options required")
			task := plans.PreparedError(err, spec)
			return []types.Task{task}
		}
	}
	var err error
	if spec.DockerContainerExec.ExecConfig.Cmd == nil {
		err = errors.New("spec for docker.container-exec exec_config.Cmd required")
	}
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	task := plans.StreamsSource{
		Producer: d.producers.ContainerExec(*spec.DockerContainerExec),
	}
	task, err = plans.SetCommonFieldsStreamsSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
