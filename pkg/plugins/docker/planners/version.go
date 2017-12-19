package planners

import (
	"path/filepath"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) Version(spec types.Spec) []types.Task {
	task := plans.StructuredSource{
		Spec:      spec,
		Producer:  d.producers.Version,
		RawPath:   filepath.Join(spec.OutputDir, "docker_version.raw"),
		JSONPath:  filepath.Join(spec.OutputDir, "docker_version.json"),
		HumanPath: filepath.Join(spec.OutputDir, "docker_version.human"),
	}
	var err error
	task, err = plans.SetCommonFieldsStructuredSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
