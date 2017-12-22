package planners

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) StackServicePs(spec types.Spec) []types.Task {
	if spec.DockerStackServicePs == nil {
		err := errors.New("spec for docker.stack-service-ps options required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	} else if spec.DockerStackServicePs.Namespace == "" {
		err := errors.New("spec for docker.stack-service-ps namespace required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	opts := spec.DockerStackServicePs.DockerServicePsOptions
	if opts == nil {
		opts = &types.DockerServicePsOptions{}
	}
	opts.Filters = dockerStackNamespaceFilter(opts.Filters, spec.DockerStackServicePs.Namespace)
	task := plans.StructuredSource{
		Spec:      spec,
		Producer:  d.producers.ServicePs(*opts),
		RawPath:   filepath.Join(spec.OutputDir, "service_ps.raw"),
		JSONPath:  filepath.Join(spec.OutputDir, "service_ps.json"),
		HumanPath: filepath.Join(spec.OutputDir, "service_ps.human"),
	}
	var err error
	task, err = plans.SetCommonFieldsStructuredSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
