package planners

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/collect/plans"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func (d *Docker) ServiceLs(spec types.Spec) []types.Task {
	if spec.DockerServiceLs == nil {
		err := errors.New("spec for docker.service-ls required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	task := plans.StructuredSource{
		Producer: d.producers.ServiceLs(*spec.DockerServiceLs),
		JSONPath: filepath.Join(spec.Shared().OutputDir, "service_ls.json"),
	}
	var err error
	task, err = plans.SetCommonFieldsStructuredSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
