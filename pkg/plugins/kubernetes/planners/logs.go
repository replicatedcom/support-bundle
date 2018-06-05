package planners

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
	"k8s.io/api/core/v1"
)

func (k *Kubernetes) Logs(spec types.Spec) []types.Task {
	var err error
	podNameProvided := spec.KubernetesLogs.Pod != ""
	labelSelectorProvided :=
		spec.KubernetesLogs.ListOptions != nil &&
			spec.KubernetesLogs.ListOptions.LabelSelector != ""

	if spec.KubernetesLogs == nil {
		err = errors.New("spec for kubernetes.logs required")
	}

	if !podNameProvided && !labelSelectorProvided {
		err = errors.New("spec for kubernetes.logs pod or list_options required")
	}

	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	var podNames []string
	if labelSelectorProvided {
		resourceListOpts := types.KubernetesResourceListOptions{
			Kind:        "pods",
			Namespace:   spec.KubernetesLogs.Namespace,
			ListOptions: spec.KubernetesLogs.ListOptions,
		}

		resources, err := k.producers.ResourceList(resourceListOpts)(context.Background())
		if err != nil {
			err := errors.Wrap(err, "Failed to list pods")
			task := plans.PreparedError(err, spec)
			return []types.Task{task}
		}

		podList := resources.(*v1.PodList)
		pods := podList.Items
		for _, pod := range pods {
			if !podNameProvided || spec.KubernetesLogs.Pod == pod.Name {
				podNames = append(podNames, pod.Name)
			}
		}

		if len(podNames) == 0 {
			err := errors.New("Unable to find any pods matching the provided pod/selector")
			task := plans.PreparedError(err, spec)
			return []types.Task{task}
		}
	} else {
		podNames = []string{spec.KubernetesLogs.Pod}
	}

	var tasks []types.Task
	for _, podName := range podNames {
		currentLogOptions := spec.KubernetesLogs
		currentLogOptions.Pod = podName

		// To support backwards compatibility
		// for non-label selectored log queries
		logFileName := fmt.Sprintf("logs-%s.raw", podName)
		if podNameProvided {
			logFileName = "logs.raw"
		}

		task := plans.StreamSource{
			Producer: k.producers.Logs(*currentLogOptions),
			RawPath:  filepath.Join(spec.Shared().OutputDir, logFileName),
		}

		task, err = plans.SetCommonFieldsStreamSource(task, spec)
		if err != nil {
			tasks = append(tasks, plans.PreparedError(err, spec))
		} else {
			tasks = append(tasks, &task)
		}
	}

	return tasks
}
