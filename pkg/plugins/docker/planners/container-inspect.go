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
		err := errors.New("spec for docker.container-inspect required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	scrubber, err := plans.RawScrubber(spec.Config.Scrub)
	if err != nil {
		err = errors.Wrap(err, "create scrubber for docker.container-inspect")
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
			Producer:    d.producers.Inspect(container.ID),
			RawScrubber: scrubber,
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
