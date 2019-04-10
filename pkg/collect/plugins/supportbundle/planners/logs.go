package planners

import (
	"path/filepath"

	"github.com/replicatedcom/support-bundle/pkg/collect/plans"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func (s *SupportBundle) Logs(spec types.Spec) []types.Task {
	task := plans.StreamSource{
		Producer: s.producers.Logs,
		RawPath:  filepath.Join(spec.Shared().OutputDir, "logs"),
	}
	var err error
	task, err = plans.SetCommonFieldsStreamSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
