package planners

import (
	"path/filepath"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (k *Kubernetes) Version(spec types.Spec) []types.Task {
	task := plans.StructuredSource{
		Producer:  k.producers.Version,
		RawPath:   filepath.Join(spec.Shared().OutputDir, "server_version.raw"),
		JSONPath:  filepath.Join(spec.Shared().OutputDir, "server_version.json"),
		HumanPath: filepath.Join(spec.Shared().OutputDir, "server_version.human"),
	}
	var err error
	task, err = plans.SetCommonFieldsStructuredSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
