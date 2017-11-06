package planners

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) ContainerInspect(spec types.Spec) []types.Task {
	if spec.DockerContainerInspect == nil {
		err := errors.New("spec for docker.container-inspect required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	containers, err := d.client.ContainerList(context.Background(), spec.DockerContainerInspect.ContainerListOptions)
	if err != nil {
		err := fmt.Errorf("container-ls-inspect failed to list containers: %v", err)
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	var tasks []types.Task

	for _, container := range containers {
		task := &plans.StructuredSource{
			Producer: d.producers.Inspect(container.ID),
		}

		basename := container.ID
		if len(container.Names) > 0 {
			basename = container.Names[0]
		}

		if spec.Raw != "" {
			task.RawPath = filepath.Join(spec.Raw, basename+".json")
		}
		if spec.JSON != "" {
			task.JSONPath = filepath.Join(spec.JSON, basename+".json")
		}
		if spec.Human != "" {
			task.HumanPath = filepath.Join(spec.Human, basename+".json")
		}

		tasks = append(tasks, task)
	}

	return tasks
}
