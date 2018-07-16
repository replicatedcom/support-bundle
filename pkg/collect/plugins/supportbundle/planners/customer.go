package planners

import (
	"path/filepath"

	"github.com/replicatedcom/support-bundle/pkg/collect/plans"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func (s *SupportBundle) CustomerMeta(spec types.Spec) []types.Task {
	task := plans.StructuredSource{
		Producer: s.producers.CustomerMeta(spec.CustomerMeta),
		JSONPath: filepath.Join(spec.Shared().OutputDir, "meta.json"),
	}
	var err error
	task, err = plans.SetCommonFieldsStructuredSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
