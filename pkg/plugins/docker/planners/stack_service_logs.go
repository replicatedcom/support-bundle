package planners

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/swarm"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) StackServiceLogs(spec types.Spec) []types.Task {
	var err error
	if spec.DockerStackServiceLogs == nil {
		err = errors.New("spec for docker.stack-service-logs required")
	} else if spec.DockerStackServiceLogs.Namespace == "" {
		err = errors.New("spec for docker.stack-service-logs namespace required")
	}
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	opts := spec.DockerStackServiceLogs.DockerServiceLsOptions
	if opts == nil {
		opts = &types.DockerServiceLsOptions{}
	}
	opts.Filters = dockerStackNamespaceFilter(opts.Filters, spec.DockerStackServiceLogs.Namespace)
	services, err := d.client.ServiceList(
		context.Background(),
		opts.ToDockerServiceListOptions())
	if err != nil {
		err := fmt.Errorf("failed to list services: %v", err)
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	filteredServices := []swarm.Service{}
	for _, service := range services {
		matched := true
		for labelKey, labelVal := range spec.DockerStackServiceLogs.Labels {
			svcLabelVal, ok := service.Spec.Labels[labelKey]
			matched = matched && ok && svcLabelVal == labelVal
		}
		if matched {
			filteredServices = append(filteredServices, service)
		}
	}

	var ts []types.Task
	for _, service := range filteredServices {
		ts = append(ts, d.serviceLogsTask(service.ID, service.Spec.Name, spec, spec.DockerStackServiceLogs.ContainerLogsOptions))
	}
	return ts
}
