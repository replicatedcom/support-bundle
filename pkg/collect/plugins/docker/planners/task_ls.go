package planners

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/collect/plans"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func (d *Docker) TaskLs(spec types.Spec) []types.Task {
	if spec.DockerTaskLs == nil {
		err := errors.New("spec for docker.task-ls required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	task := plans.StructuredSource{
		Producer: d.producers.TaskLs(*spec.DockerTaskLs),
		JSONPath: filepath.Join(spec.Shared().OutputDir, "task_ls.json"),
	}
	var err error
	task, err = plans.SetCommonFieldsStructuredSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
