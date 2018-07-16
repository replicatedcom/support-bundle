package planners

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/collect/plans"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func (k *Kubernetes) ResourceList(spec types.Spec) []types.Task {
	var err error
	if spec.KubernetesResourceList == nil {
		err = errors.New("spec for kubernetes.resource-list required")
	} else if spec.KubernetesResourceList.Kind == "" {
		err = errors.New("spec for kubernetes.resource-list kind required")
	}
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	task := plans.StructuredSource{
		Producer: k.producers.ResourceList(*spec.KubernetesResourceList),
		JSONPath: filepath.Join(spec.Shared().OutputDir, "resource.json"),
	}
	task, err = plans.SetCommonFieldsStructuredSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
