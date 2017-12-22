package planners

import (
	"context"
	"fmt"
	"path/filepath"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
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

	var ts []types.Task
	containers, err := d.client.ContainerList(
		context.Background(),
		spec.DockerContainerLogs.ContainerListOptions.ToDockerContainerListOptions())
	if err != nil {
		err := fmt.Errorf("failed to list containers: %v", err)
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	for _, container := range containers {
		ts = append(ts, d.containerLogsTask(container.ID, container.Names, spec))
	}
	return ts
}

func (d *Docker) containerLogsTask(id string, names []string, spec types.Spec) types.Task {
	basename := id
	if len(names) > 0 {
		basename = names[0]
	}
	task := plans.StreamSource{
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

// get the ID of a container given the name
func (d *Docker) getContainerID(ctx context.Context, name string) (string, error) {
	filter := filters.NewArgs()
	filter.Add("name", name)
	containers, err := d.client.ContainerList(ctx, dockertypes.ContainerListOptions{
		All:     true,
		Filters: filter,
	})
	if err != nil {
		return "", err
	} else if len(containers) == 0 {
		return "", fmt.Errorf("unable to find container with name %s", name)
	}
	return containers[0].ID, nil
}
