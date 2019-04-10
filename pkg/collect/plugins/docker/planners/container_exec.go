package planners

import (
	"path"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/collect/plans"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
	jww "github.com/spf13/jwalterweatherman"
)

func (d *Docker) ContainerExec(spec types.Spec) []types.Task {
	if spec.DockerContainerExec == nil {
		if spec.DockerExec != nil {
			jww.DEBUG.Println("spec for docker.container-exec aliased from docker.exec")
			spec.DockerContainerExec = spec.DockerExec
		} else {
			err := errors.New("spec for docker.container-exec required")
			task := plans.PreparedError(err, spec)
			return []types.Task{task}
		}
	}

	if spec.DockerContainerExec.Container == "" && len(spec.DockerContainerExec.Labels) == 0 {
		err := errors.New("spec for docker.container-exec container name or label required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	if spec.DockerContainerExec.ExecConfig.Cmd == nil {
		err := errors.New("spec for docker.container-exec exec_config.Cmd required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	containers, err := d.SelectContainerIDs(spec.DockerContainerExec.Container, spec.DockerContainerExec.Labels)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	if len(containers) == 0 {
		task := plans.PreparedError(errors.New("Unable to find any containers matching the provided labels"), spec)
		return []types.Task{task}
	}

	var tasks []types.Task
	for _, container := range containers {
		task := plans.StreamsSource{
			Producer: d.producers.ContainerExec(container, spec.DockerContainerExec.ExecConfig),
		}
		task, err := plans.SetCommonFieldsStreamsSource(task, spec)
		if len(containers) > 1 {
			//if multiple containers match, add the container id to the path to avoid collisions
			task.RawPath = path.Join(task.RawPath, container)
		}

		if err != nil {
			tasks = append(tasks, plans.PreparedError(err, spec))
		} else {
			tasks = append(tasks, &task)
		}
	}
	return tasks
}
