package planners

import (
	"path/filepath"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) Version(spec types.Spec) []types.Task {
	task := plans.StructuredSource{
		Producer: d.producers.Version,
		JSONPath: filepath.Join(spec.Shared().OutputDir, "docker_version.json"),
	}
	var err error
	task, err = plans.SetCommonFieldsStructuredSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
