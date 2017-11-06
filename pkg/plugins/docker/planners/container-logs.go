package planners

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) ContainerLogs(spec types.Spec) []types.Task {
	if spec.DockerContainerLogs == nil {
		err := errors.New("spec for docker.container-logs required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	containers, err := d.client.ContainerList(context.Background(), spec.DockerContainerLogs.ContainerListOptions)
	if err != nil {
		err := fmt.Errorf("container-ls-logs failed to list containers: %v", err)
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	var tasks []types.Task

	for _, container := range containers {
		task := &plans.StreamSource{
			Producer: d.producers.Logs(container.ID),
		}

		basename := container.ID
		if len(container.Names) > 0 {
			basename = container.Names[0]
		}

		if spec.Raw != "" {
			task.RawPath = filepath.Join(spec.Raw, basename+".log")
		}
		if spec.JSON != "" {
			task.JSONPath = filepath.Join(spec.JSON, basename+".log")
		}
		if spec.Human != "" {
			task.HumanPath = filepath.Join(spec.Human, basename+".log")
		}

		tasks = append(tasks, task)
	}

	return tasks
}
