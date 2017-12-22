package planners

import (
	"context"
	"fmt"
	"path/filepath"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) ServiceLogs(spec types.Spec) []types.Task {
	if spec.DockerServiceLogs == nil {
		err := errors.New("spec for docker.service-logs options required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	if spec.DockerServiceLogs.ID != "" {
		return []types.Task{d.serviceLogsTask(spec.DockerServiceLogs.ID, "", spec, spec.DockerServiceLogs.ContainerLogsOptions)}
	}

	if spec.DockerServiceLogs.Name != "" {
		serviceID, err := d.getServiceID(context.Background(), spec.DockerServiceLogs.Name)
		if err != nil {
			task := plans.PreparedError(err, spec)
			return []types.Task{task}
		}
		return []types.Task{d.serviceLogsTask(serviceID, spec.DockerServiceLogs.Name, spec, spec.DockerServiceLogs.ContainerLogsOptions)}
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
	task := plans.StreamSource{
		Spec:     spec,
		Producer: d.producers.ServiceLogs(id, opts),
		RawPath:  filepath.Join(spec.OutputDir, basename+".raw"),
	}
	var err error
	task, err = plans.SetCommonFieldsStreamSource(task, spec)
	if err != nil {
		return plans.PreparedError(err, spec)
	}
	return &task
}

// get the ID of a service given the name
func (d *Docker) getServiceID(ctx context.Context, name string) (string, error) {
	filter := filters.NewArgs()
	filter.Add("name", name)
	services, err := d.client.ServiceList(ctx, dockertypes.ServiceListOptions{
		Filters: filter,
	})
	if err != nil {
		return "", err
	} else if len(services) == 0 {
		return "", fmt.Errorf("unable to find service with name %s", name)
	}
	return services[0].ID, nil
}
