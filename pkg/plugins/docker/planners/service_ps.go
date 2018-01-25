package planners

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) ServicePs(spec types.Spec) []types.Task {
	if spec.DockerServicePs == nil {
		err := errors.New("spec for docker.service-ps required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	task := plans.StructuredSource{
		Producer:  d.producers.TaskLs(spec.DockerServicePs.DockerTaskLsOptions),
		RawPath:   filepath.Join(spec.Shared().OutputDir, "service_ps.raw"),
		JSONPath:  filepath.Join(spec.Shared().OutputDir, "service_ps.json"),
		HumanPath: filepath.Join(spec.Shared().OutputDir, "service_ps.human"),
	}
	var err error
	task, err = plans.SetCommonFieldsStructuredSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
