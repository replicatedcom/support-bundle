package planners

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/collect/plans"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
	v1 "k8s.io/api/core/v1"
)

func (k *Kubernetes) Logs(spec types.Spec) []types.Task {
	var err error
	podNameProvided := spec.KubernetesLogs.Pod != ""
	podListOptionsProvided := spec.KubernetesLogs.ListOptions != nil
	namespaceProvided := spec.KubernetesLogs.Namespace != ""

	if spec.KubernetesLogs == nil {
		err = errors.New("spec for kubernetes.logs required")
	}

	if !podNameProvided && !podListOptionsProvided {
		err = errors.New("spec for kubernetes.logs pod or list_options required")
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
			if spec.KubernetesLogs.PodLogOptions != nil && spec.KubernetesLogs.PodLogOptions.Container != "" {
				podLocations = append(podLocations, podLocation{
					PodName:       pod.Name,
					Namespace:     pod.Namespace,
					ContainerName: spec.KubernetesLogs.PodLogOptions.Container,
				})
			} else {
				// get logs for all containers in the pod
				for _, container := range pod.Spec.Containers {
					podLocations = append(podLocations, podLocation{
						PodName:       pod.Name,
						Namespace:     pod.Namespace,
						ContainerName: container.Name,
					})
				}
			}
		}
	}

	if len(podLocations) == 0 {
		err := errors.New("Unable to find any pods matching the provided pod/selector")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	var tasks []types.Task
	for _, podLocation := range podLocations {
		currentLogOptions := *spec.KubernetesLogs
		currentLogOptions.Pod = podLocation.PodName
		if currentLogOptions.PodLogOptions == nil {
			currentLogOptions.PodLogOptions = new(v1.PodLogOptions)
		}
		currentLogOptions.PodLogOptions.Container = podLocation.ContainerName
		currentLogOptions.Namespace = podLocation.Namespace

		rawPath := spec.Shared().OutputDir
		if !namespaceProvided {
			rawPath = filepath.Join(rawPath, podLocation.Namespace)
		}

		task := plans.StreamSource{
			Producer: k.producers.Logs(currentLogOptions),
			RawPath:  filepath.Join(rawPath, fmt.Sprintf("%s-%s.log", podLocation.PodName, podLocation.ContainerName)),
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
