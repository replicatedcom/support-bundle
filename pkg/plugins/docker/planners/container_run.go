package planners

import (
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
	jww "github.com/spf13/jwalterweatherman"
)

func (d *Docker) ContainerRun(spec types.Spec) []types.Task {
	if spec.DockerContainerRun == nil {
		if spec.DockerRun != nil {
			jww.DEBUG.Println("spec for docker.container-run aliased from docker.run")
			spec.DockerContainerRun = spec.DockerRun
		} else {
			err := errors.New("spec for docker.container-run options required")
			task := plans.PreparedError(err, spec)
			return []types.Task{task}
		}
	}
	var err error
	if spec.DockerContainerRun.ContainerCreateConfig.Config == nil {
		err = errors.New("spec for docker.container-run ContainerCreateConfig.Config required")
	} else if spec.DockerContainerRun.ContainerCreateConfig.Config.Image == "" {
		err = errors.New("spec for docker.container-run ContainerCreateConfig.Config.Image required")
	}
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	task := plans.StreamsSource{
		Producer: d.producers.ContainerRun(*spec.DockerContainerRun),
	}
	task, err = plans.SetCommonFieldsStreamsSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
