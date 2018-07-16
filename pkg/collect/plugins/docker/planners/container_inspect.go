package planners

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/collect/plans"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func (d *Docker) ContainerInspect(spec types.Spec) []types.Task {
	if spec.DockerContainerInspect == nil {
		err := errors.New("spec for docker.container-inspect required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	if spec.DockerContainerInspect.Container != "" {
		return []types.Task{d.containerInspectTask(spec.DockerContainerInspect.Container, nil, spec)}
	}

	var ts []types.Task
	containers, err := d.client.ContainerList(
		context.Background(),
		spec.DockerContainerInspect.ContainerListOptions.ToDockerContainerListOptions())
	if err != nil {
		err := fmt.Errorf("failed to list containers: %v", err)
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	for _, container := range containers {
		ts = append(ts, d.containerInspectTask(container.ID, container.Names, spec))
	}
	return ts
}

func (d *Docker) containerInspectTask(id string, names []string, spec types.Spec) types.Task {
	basename := id
	if len(names) > 0 {
		basename = names[0]
	}
	task := plans.StructuredSource{
		Producer: d.producers.ContainerInspect(id),
		JSONPath: filepath.Join(spec.Shared().OutputDir, basename+".json"),
	}
	var err error
	task, err = plans.SetCommonFieldsStructuredSource(task, spec)
	if err != nil {
		return plans.PreparedError(err, spec)
	}
	return &task
}
