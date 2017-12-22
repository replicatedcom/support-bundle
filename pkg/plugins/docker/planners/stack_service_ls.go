package planners

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) StackServiceLs(spec types.Spec) []types.Task {
	if spec.DockerStackServiceLs == nil {
		err := errors.New("spec for docker.stack-service-ls options required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	} else if spec.DockerStackServiceLs.Namespace == "" {
		err := errors.New("spec for docker.stack-service-ls namespace required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	opts := spec.DockerStackServiceLs.DockerServiceLsOptions
	if opts == nil {
		opts = &types.DockerServiceLsOptions{}
	}
	opts.Filters = dockerStackNamespaceFilter(opts.Filters, spec.DockerStackServiceLs.Namespace)
	task := plans.StructuredSource{
		Spec:      spec,
		Producer:  d.producers.ServiceLs(*opts),
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

func dockerStackNamespaceFilter(args map[string][]string, namespace string) map[string][]string {
	if args == nil {
		args = map[string][]string{}
	}
	args["label"] = append(args["label"], "com.docker.stack.namespace="+namespace)
	return args
}
