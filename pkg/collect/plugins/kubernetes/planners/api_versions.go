package planners

import (
	"path/filepath"

	"github.com/replicatedcom/support-bundle/pkg/collect/plans"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func (k *Kubernetes) APIVersions(spec types.Spec) []types.Task {
	task := plans.StructuredSource{
		Producer: k.producers.APIVersions,
		JSONPath: filepath.Join(spec.Shared().OutputDir, "api_versions.json"),
	}
	var err error
	task, err = plans.SetCommonFieldsStructuredSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
