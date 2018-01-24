package planners

import (
	"path/filepath"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (s *SupportBundle) CustomerMeta(spec types.Spec) []types.Task {
	task := plans.StructuredSource{
		Producer:  s.producers.CustomerMeta(spec.CustomerMeta),
		RawPath:   filepath.Join(spec.Shared().OutputDir, "meta.raw"),
		JSONPath:  filepath.Join(spec.Shared().OutputDir, "meta.json"),
		HumanPath: filepath.Join(spec.Shared().OutputDir, "meta.human"),
	}
	var err error
	task, err = plans.SetCommonFieldsStructuredSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
