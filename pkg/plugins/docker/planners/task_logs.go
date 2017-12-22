package planners

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) TaskLogs(spec types.Spec) []types.Task {
	if spec.DockerTaskLogs == nil {
		err := errors.New("spec for docker.task-logs options required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	if spec.DockerTaskLogs.ID != "" {
		return []types.Task{d.taskLogsTask(spec.DockerTaskLogs.ID, spec)}
	}

	var ts []types.Task
	tasks, err := d.client.TaskList(
		context.Background(),
		spec.DockerTaskLogs.TaskListOptions.ToDockerTaskListOptions())
	if err != nil {
		err := fmt.Errorf("failed to list tasks: %v", err)
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	for _, task := range tasks {
		ts = append(ts, d.taskLogsTask(task.ID, spec))
	}
	return ts
}

func (d *Docker) taskLogsTask(id string, spec types.Spec) types.Task {
	basename := id
	task := plans.StreamSource{
		Spec:     spec,
		Producer: d.producers.TaskLogs(basename, spec.DockerTaskLogs.ContainerLogsOptions),
		RawPath:  filepath.Join(spec.OutputDir, basename+".raw"),
	}
	var err error
	task, err = plans.SetCommonFieldsStreamSource(task, spec)
	if err != nil {
		return plans.PreparedError(err, spec)
	}
	return &task
}
