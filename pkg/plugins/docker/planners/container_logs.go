package planners

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
	jww "github.com/spf13/jwalterweatherman"
)

func (d *Docker) ContainerLogs(spec types.Spec) []types.Task {
	if spec.DockerContainerLogs == nil {
		if spec.DockerLogs != nil {
			jww.DEBUG.Println("spec for docker.container-logs aliased from docker.logs")
			spec.DockerContainerLogs = spec.DockerLogs
		} else {
			err := errors.New("spec for docker.container-logs required")
			task := plans.PreparedError(err, spec)
			return []types.Task{task}
		}
	}

	if spec.DockerContainerLogs.Container != "" {
		return []types.Task{d.containerLogsTask(spec.DockerContainerLogs.Container, nil, spec)}
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
	task := plans.StreamsSource{
		Producer: d.producers.ContainerLogs(id, basename, spec.DockerContainerLogs.ContainerLogsOptions),
		RawPath:  spec.Shared().OutputDir,
	}
	var err error
	task, err = plans.SetCommonFieldsStreamsSource(task, spec)
	if err != nil {
		return plans.PreparedError(err, spec)
	}
	return &task
}
