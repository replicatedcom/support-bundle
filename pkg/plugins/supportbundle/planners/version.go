package planners

import (
	"path/filepath"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (s *SupportBundle) Version(spec types.Spec) []types.Task {
	task := plans.StructuredSource{
		Producer:  s.producers.Version,
		RawPath:   filepath.Join(spec.OutputDir, "VERSION.raw"),
		JSONPath:  filepath.Join(spec.OutputDir, "VERSION.json"),
		HumanPath: filepath.Join(spec.OutputDir, "VERSION.human"),
	}
	var err error
	task, err = plans.SetCommonFieldsStructuredSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
