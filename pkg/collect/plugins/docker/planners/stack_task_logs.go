package planners

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/swarm"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/collect/plans"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func (d *Docker) StackTaskLogs(spec types.Spec) []types.Task {
	var err error
	if spec.DockerStackTaskLogs == nil {
		err = errors.New("spec for docker.stack-task-logs required")
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

	filteredTasks := []swarm.Task{}
	for _, task := range tasks {
		matched := true
		for labelKey, labelVal := range spec.DockerStackTaskLogs.Labels {
			svcLabelVal, ok := task.Spec.ContainerSpec.Labels[labelKey]
			matched = matched && ok && svcLabelVal == labelVal
		}
		if matched {
			filteredTasks = append(filteredTasks, task)
		}
	}

	if len(filteredTasks) == 0 {
		err := fmt.Errorf("no tasks found in namespace %s with labels %+v",
			spec.DockerStackTaskLogs.Namespace,
			spec.DockerStackTaskLogs.Labels)
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	var ts []types.Task
	for _, task := range filteredTasks {
		ts = append(ts, d.taskLogsTask(task.ID, spec, spec.DockerStackTaskLogs.ContainerLogsOptions))
	}
	return ts
}
