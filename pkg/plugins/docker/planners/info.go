package planners

import (
	"path/filepath"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) Info(spec types.Spec) []types.Task {
	task := plans.StructuredSource{
		Spec:      spec,
		Producer:  d.producers.Info,
		RawPath:   filepath.Join(spec.OutputDir, "docker_info.raw"),
		JSONPath:  filepath.Join(spec.OutputDir, "docker_info.json"),
		HumanPath: filepath.Join(spec.OutputDir, "docker_info.human"),
	}
	var err error
	task, err = plans.SetCommonFieldsStructuredSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
