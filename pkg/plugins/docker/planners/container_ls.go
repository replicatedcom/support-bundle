package planners

import (
	"errors"
	"path/filepath"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
	jww "github.com/spf13/jwalterweatherman"
)

func (d *Docker) ContainerLs(spec types.Spec) []types.Task {
	if spec.DockerContainerLs == nil {
		if spec.DockerPs != nil {
			jww.DEBUG.Println("spec for docker.container-ls aliased from docker.ps")
			spec.DockerContainerLs = spec.DockerPs
		} else {
			err := errors.New("spec for docker.container-ls required")
			task := plans.PreparedError(err, spec)
			return []types.Task{task}
		}
	}

	task := plans.StructuredSource{
		Producer: d.producers.ContainerLs(*spec.DockerContainerLs),
		JSONPath: filepath.Join(spec.Shared().OutputDir, "container_ls.json"),
	}
	var err error
	task, err = plans.SetCommonFieldsStructuredSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
