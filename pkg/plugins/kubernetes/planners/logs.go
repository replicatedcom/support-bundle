package planners

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (k *Kubernetes) Logs(spec types.Spec) []types.Task {
	var err error
	if spec.KubernetesLogs == nil {
		err = errors.New("spec for kubernetes.logs required")
	} else if spec.KubernetesLogs.Pod == "" {
		err = errors.New("spec for kubernetes.logs pod required")
	}
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	task := plans.StreamSource{
		Producer: k.producers.Logs(*spec.KubernetesLogs),
		RawPath:  filepath.Join(spec.OutputDir, "logs.raw"),
	}
	task, err = plans.SetCommonFieldsStreamSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
