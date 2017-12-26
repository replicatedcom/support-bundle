package planners

import (
	"github.com/pkg/errors"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) ContainerCp(spec types.Spec) []types.Task {
	var err error
	if spec.DockerContainerCp == nil {
		err = errors.New("spec for docker.container-cp options required")
	} else if spec.DockerContainerCp.Container == "" {
		err = errors.New("spec for docker.container-cp container required")
	} else if spec.DockerContainerCp.SrcPath == "" {
		err = errors.New("spec for docker.container-cp src_path required")
	}
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	task := plans.StreamSource{
		Producer:     d.producers.ContainerCp(spec.DockerContainerCp.Container, spec.DockerContainerCp.SrcPath),
		StreamFormat: plans.StreamFormatTar,
	}
	task, err = plans.SetCommonFieldsStreamSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
