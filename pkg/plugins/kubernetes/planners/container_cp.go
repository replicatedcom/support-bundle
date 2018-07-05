package planners

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
	"k8s.io/api/core/v1"
)

func (k *Kubernetes) ContainerCp(spec types.Spec) []types.Task {
	var err error
	podNameProvided := spec.KubernetesContainerCp.Pod != ""
	labelSelectorProvided :=
		spec.KubernetesContainerCp.PodListOptions != nil &&
			spec.KubernetesContainerCp.PodListOptions.LabelSelector != ""

	if spec.KubernetesContainerCp == nil {
		err = errors.New("spec for kubernetes.ContainerCp required")
	}

	if !podNameProvided && !labelSelectorProvided {
		err = errors.New("spec for kubernetes.ContainerCp pod or list_options required")
	}

	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	var podNames []string

	if labelSelectorProvided {
		resourceListOpts := types.KubernetesResourceListOptions{
			Kind:        "pods",
			Namespace:   spec.KubernetesContainerCp.Namespace,
			ListOptions: spec.KubernetesContainerCp.PodListOptions,
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
			if !podNameProvided || spec.KubernetesContainerCp.Pod == pod.Name {
				podNames = append(podNames, pod.Name)
			}
		}

		if len(pods) == 0 {
			err := errors.New("unable to find any pods matching the provided pod/selector")
			task := plans.PreparedError(err, spec)
			return []types.Task{task}
		}
	} else {
		podNames = []string{spec.KubernetesContainerCp.Pod}
	}

	var tasks []types.Task
	for _, podName := range podNames {
		rawPath := spec.Shared().OutputDir
		if !podNameProvided {
			rawPath = filepath.Join(rawPath, podName)
		}

		task := plans.StreamsSource{
			RawPath:      rawPath,
			StreamFormat: plans.StreamFormatTar,
			Producer: k.producers.ContainerCp(
				podName,
				spec.KubernetesContainerCp.Container,
				spec.KubernetesContainerCp.Namespace,
				filepath.Clean(spec.KubernetesContainerCp.SrcPath)),
		}

		task, err = plans.SetCommonFieldsStreamsSource(task, spec)
		if err != nil {
			tasks = append(tasks, plans.PreparedError(err, spec))
		} else {
			tasks = append(tasks, &task)
		}
	}

	return tasks
}
