package planners

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/collect/plans"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
	v1 "k8s.io/api/core/v1"
)

func (k *Kubernetes) ContainerCp(spec types.Spec) []types.Task {
	var err error
	podNameProvided := spec.KubernetesContainerCp.Pod != ""
	podListOptionsProvided := spec.KubernetesContainerCp.PodListOptions != nil
	namespaceProvided := spec.KubernetesContainerCp.Namespace != ""

	if spec.KubernetesContainerCp == nil {
		err = errors.New("spec for kubernetes.ContainerCp required")
	}

	if !podNameProvided && !podListOptionsProvided {
		err = errors.New("spec for kubernetes.ContainerCp pod or list_options required")
	}

	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	type podLocation struct {
		PodName       string
		ContainerName string
		Namespace     string
	}

	var podLocations []podLocation

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
			l := podLocation{PodName: pod.Name, Namespace: pod.Namespace}
			if spec.KubernetesContainerCp.Container != "" {
				l.ContainerName = spec.KubernetesContainerCp.Container
			} else if len(pod.Spec.Containers) > 1 {
				// choose first container in pod
				l.ContainerName = pod.Spec.Containers[0].Name
			}
			podLocations = append(podLocations, l)
		}
	}

	if len(podLocations) == 0 {
		err := errors.New("unable to find any pods matching the provided pod/selector")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	var tasks []types.Task
	for _, podLocation := range podLocations {
		rawPath := spec.Shared().OutputDir
		if !namespaceProvided {
			rawPath = filepath.Join(rawPath, podLocation.Namespace)
		}
		if !podNameProvided {
			rawPath = filepath.Join(rawPath, podLocation.PodName)
		}

		task := plans.StreamsSource{
			RawPath:      rawPath,
			StreamFormat: plans.StreamFormatTar,
			Producer: k.producers.ContainerCp(
				podLocation.PodName,
				podLocation.ContainerName,
				podLocation.Namespace,
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
