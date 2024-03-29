package planners

import (
	"context"
	"fmt"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/collect/plans"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func (d *Docker) TaskLogs(spec types.Spec) []types.Task {
	if spec.DockerTaskLogs == nil {
		err := errors.New("spec for docker.task-logs required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	if spec.DockerTaskLogs.ID != "" {
		return []types.Task{d.taskLogsTask(spec.DockerTaskLogs.ID, spec, spec.DockerTaskLogs.ContainerLogsOptions)}
	}

	tasks, err := d.client.TaskList(
		context.Background(),
		spec.DockerTaskLogs.TaskListOptions.ToDockerTaskListOptions())
	if err != nil {
		err := fmt.Errorf("failed to list tasks: %v", err)
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	var ts []types.Task
	for _, task := range tasks {
		ts = append(ts, d.taskLogsTask(task.ID, spec, spec.DockerTaskLogs.ContainerLogsOptions))
	}
	return ts
}

func (d *Docker) taskLogsTask(id string, spec types.Spec, opts *dockertypes.ContainerLogsOptions) types.Task {
	basename := id
	task := plans.StreamsSource{
		Producer: d.producers.TaskLogs(basename, opts),
		RawPath:  spec.Shared().OutputDir,
	}
	var err error
	task, err = plans.SetCommonFieldsStreamsSource(task, spec)
	if err != nil {
		return plans.PreparedError(err, spec)
	}
	return &task
}
