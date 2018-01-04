package planners

import (
	"path/filepath"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (k *Kubernetes) ClusterInfo(spec types.Spec) []types.Task {
	task := plans.StructuredSource{
		Producer:  k.producers.ClusterInfo,
		RawPath:   filepath.Join(spec.OutputDir, "cluster_info.raw"),
		JSONPath:  filepath.Join(spec.OutputDir, "cluster_info.json"),
		HumanPath: filepath.Join(spec.OutputDir, "cluster_info.human"),
	}
	var err error
	task, err = plans.SetCommonFieldsStructuredSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}