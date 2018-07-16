package planners

import (
	"context"
	"fmt"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/collect/plans"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func (d *Docker) ServiceLogs(spec types.Spec) []types.Task {
	if spec.DockerServiceLogs == nil {
		err := errors.New("spec for docker.service-logs required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	if spec.DockerServiceLogs.Service != "" {
		return []types.Task{d.serviceLogsTask(spec.DockerServiceLogs.Service, "", spec, spec.DockerServiceLogs.ContainerLogsOptions)}
	}

	services, err := d.client.ServiceList(
		context.Background(),
		spec.DockerServiceLogs.ServiceListOptions.ToDockerServiceListOptions())
	if err != nil {
		err := fmt.Errorf("failed to list services: %v", err)
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	var ts []types.Task
	for _, service := range services {
		ts = append(ts, d.serviceLogsTask(service.ID, service.Spec.Name, spec, spec.DockerServiceLogs.ContainerLogsOptions))
	}
	return ts
}

func (d *Docker) serviceLogsTask(id string, name string, spec types.Spec, opts *dockertypes.ContainerLogsOptions) types.Task {
	basename := id
	if name != "" {
		basename = name
	}
	task := plans.StreamsSource{
		Producer: d.producers.ServiceLogs(id, basename, opts),
		RawPath:  spec.Shared().OutputDir,
	}
	var err error
	task, err = plans.SetCommonFieldsStreamsSource(task, spec)
	if err != nil {
		return plans.PreparedError(err, spec)
	}
	return &task
}
