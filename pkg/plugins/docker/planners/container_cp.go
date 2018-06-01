package planners

import (
	"context"
	"path"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/pkg/errors"
)

func (d *Docker) ContainerCp(spec types.Spec) []types.Task {
	var err error
	if spec.DockerContainerCp == nil {
		err = errors.New("spec for docker.container-cp required")
	} else if spec.DockerContainerCp.Container == "" && len(spec.DockerContainerCp.Labels) == 0 {
		err = errors.New("spec for docker.container-cp container or label required")
	} else if spec.DockerContainerCp.SrcPath == "" {
		err = errors.New("spec for docker.container-cp src_path required")
	}
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	var containers []string
	if spec.DockerContainerCp.Container != "" {
		containers = append(containers, spec.DockerContainerCp.Container)
	} else {
		//get all containers with matching labels
		labelFilters := dockertypes.ContainerListOptions{}
		labelFilters.Filters = filters.NewArgs()
		for _, val := range spec.DockerContainerCp.Labels {
			labelFilters.Filters.Add("label", val)
		}

		labelledContainers, err := d.client.ContainerList(context.Background(), labelFilters)
		if err != nil {
			err = errors.Wrap(err, "Unable to find containers to read file from")
			task := plans.PreparedError(err, spec)
			return []types.Task{task}
		}

		for _, container := range labelledContainers {
			containers = append(containers, container.ID)
		}

		if len(containers) == 0 {
			task := plans.PreparedError(errors.New("Unable to find any containers matching the provided labels"), spec)
			return []types.Task{task}
		}
	}

	var tasks []types.Task

	for _, container := range containers {
		task := plans.StreamSource{
			Producer:     d.producers.ContainerCp(container, spec.DockerContainerCp.SrcPath),
			StreamFormat: plans.StreamFormatTar,
		}
		task, err = plans.SetCommonFieldsStreamSource(task, spec)

		if len(containers) > 1 {
			//if multiple containers match, add the container id to the path to avoid collisions
			task.RawPath = path.Join(task.RawPath, container)
		}

		if err != nil {
			task := plans.PreparedError(err, spec)
			tasks = append(tasks, task)
		} else {
			tasks = append(tasks, &task)
		}
	}

	return tasks
}
