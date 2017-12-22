package planners

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) ServiceLs(spec types.Spec) []types.Task {
	if spec.DockerServiceLs == nil {
		err := errors.New("spec for docker.service-ls options required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	task := plans.StructuredSource{
		Spec:      spec,
		Producer:  d.producers.ServiceLs(*spec.DockerServiceLs),
		RawPath:   filepath.Join(spec.OutputDir, "service_ls.raw"),
		JSONPath:  filepath.Join(spec.OutputDir, "service_ls.json"),
		HumanPath: filepath.Join(spec.OutputDir, "service_ls.human"),
	}
	var err error
	task, err = plans.SetCommonFieldsStructuredSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
