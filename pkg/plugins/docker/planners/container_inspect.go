package planners

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) ContainerInspect(spec types.Spec) []types.Task {
	if spec.DockerContainerInspect == nil {
		err := errors.New("spec for docker.container-inspect options required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	if spec.DockerContainerInspect.ID != "" {
		return []types.Task{d.containerInspectTask(spec.DockerContainerInspect.ID, nil, spec)}
	}

	if spec.DockerContainerInspect.Name != "" {
		containerID, err := d.getContainerID(context.Background(), spec.DockerContainerInspect.Name)
		if err != nil {
			task := plans.PreparedError(err, spec)
			return []types.Task{task}
		}
		return []types.Task{d.containerInspectTask(containerID, []string{spec.DockerContainerInspect.Name}, spec)}
	}

	var tasks []types.Task
	containers, err := d.client.ContainerList(
		context.Background(),
		spec.DockerContainerInspect.ContainerListOptions.ToDockerContainerListOptions())
	if err != nil {
		err := fmt.Errorf("failed to list containers: %v", err)
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	for _, container := range containers {
		tasks = append(tasks, d.containerInspectTask(container.ID, container.Names, spec))
	}
	return tasks
}

func (d *Docker) containerInspectTask(id string, names []string, spec types.Spec) types.Task {
	basename := id
	if len(names) > 0 {
		basename = names[0]
	}
	task := plans.StructuredSource{
		Spec:      spec,
		Producer:  d.producers.ContainerInspect(id),
		RawPath:   filepath.Join(spec.OutputDir, basename+".raw"),
		JSONPath:  filepath.Join(spec.OutputDir, basename+".json"),
		HumanPath: filepath.Join(spec.OutputDir, basename+".human"),
	}
	var err error
	task, err = plans.SetCommonFieldsStructuredSource(task, spec)
	if err != nil {
		return plans.PreparedError(err, spec)
	}
	return &task
}
