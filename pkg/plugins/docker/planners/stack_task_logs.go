package planners

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) StackTaskLogs(spec types.Spec) []types.Task {
	var err error
	if spec.DockerStackTaskLogs == nil {
		err = errors.New("spec for docker.stack-task-logs options required")
	} else if spec.DockerStackTaskLogs.Namespace == "" {
		err = errors.New("spec for docker.stack-task-logs namespace required")
	}
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	opts := spec.DockerStackTaskLogs.DockerServicePsOptions
	if opts == nil {
		opts = &types.DockerServicePsOptions{}
	}
	opts.Filters = dockerStackNamespaceFilter(opts.Filters, spec.DockerStackTaskLogs.Namespace)
	tasks, err := d.client.TaskList(
		context.Background(),
		opts.ToDockerTaskListOptions())
	if err != nil {
		err := fmt.Errorf("failed to list tasks: %v", err)
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	var ts []types.Task
	for _, task := range tasks {
		ts = append(ts, d.taskLogsTask(task.ID, spec, spec.DockerStackTaskLogs.ContainerLogsOptions))
	}
	return ts
}
