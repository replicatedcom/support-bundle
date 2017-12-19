package planners

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) ContainerLogs(spec types.Spec) []types.Task {
	if spec.DockerContainerLogs == nil {
		err := errors.New("spec for docker.container-logs options required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	if spec.DockerContainerLogs.ID != "" {
		return []types.Task{d.containerLogsTask(spec.DockerContainerLogs.ID, nil, spec)}
	}

	if spec.DockerContainerLogs.Name != "" {
		containerID, err := d.getContainerID(context.Background(), spec.DockerContainerLogs.Name)
		if err != nil {
			task := plans.PreparedError(err, spec)
			return []types.Task{task}
		}
		return []types.Task{d.containerLogsTask(containerID, []string{spec.DockerContainerLogs.Name}, spec)}
	}

	var tasks []types.Task
	containers, err := d.client.ContainerList(
		context.Background(),
		spec.DockerContainerLogs.ContainerListOptions.ToDockerContainerListOptions())
	if err != nil {
		err := fmt.Errorf("failed to list containers: %v", err)
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	for _, container := range containers {
		tasks = append(tasks, d.containerLogsTask(container.ID, container.Names, spec))
	}
	return tasks
}

func (d *Docker) containerLogsTask(id string, names []string, spec types.Spec) types.Task {
	basename := id
	if len(names) > 0 {
		basename = names[0]
	}
	task := plans.StreamSource{
		Spec:     spec,
		Producer: d.producers.ContainerLogs(id, spec.DockerContainerLogs.ContainerLogsOptions),
		RawPath:  filepath.Join(spec.OutputDir, basename+".raw"),
	}
	var err error
	task, err = plans.SetCommonFieldsStreamSource(task, spec)
	if err != nil {
		return plans.PreparedError(err, spec)
	}
	return &task
}
