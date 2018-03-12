package planners

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) StackServicePs(spec types.Spec) []types.Task {
	var err error
	if spec.DockerStackServicePs == nil {
		err = errors.New("spec for docker.stack-service-ps required")
	} else if spec.DockerStackServicePs.Namespace == "" {
		err = errors.New("spec for docker.stack-service-ps namespace required")
	}
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	opts := spec.DockerStackServicePs.DockerServicePsOptions
	if opts == nil {
		opts = &types.DockerServicePsOptions{}
	}
	opts.Filters = dockerStackNamespaceFilter(opts.Filters, spec.DockerStackServicePs.Namespace)
	task := plans.StructuredSource{
		Producer: d.producers.TaskLs(opts.DockerTaskLsOptions),
		JSONPath: filepath.Join(spec.Shared().OutputDir, "service_ps.json"),
	}
	task, err = plans.SetCommonFieldsStructuredSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
