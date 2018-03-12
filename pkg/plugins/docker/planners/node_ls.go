package planners

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) NodeLs(spec types.Spec) []types.Task {
	if spec.DockerNodeLs == nil {
		err := errors.New("spec for docker.node-ls required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	task := plans.StructuredSource{
		Producer: d.producers.NodeLs(*spec.DockerNodeLs),
		JSONPath: filepath.Join(spec.Shared().OutputDir, "node_ls.json"),
	}
	var err error
	task, err = plans.SetCommonFieldsStructuredSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
