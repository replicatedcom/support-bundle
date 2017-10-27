package planners

import (
	"errors"
	"time"

	"github.com/replicatedcom/support-bundle/plans"
	"github.com/replicatedcom/support-bundle/types"
)

func (d *Docker) ReadFile(spec types.Spec) []types.Task {
	if spec.Config.ContainerID == "" && spec.Config.ContainerName == "" {
		err := errors.New("spec for docker.read-file requires a container ID or Name within config")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	} else if spec.Config.FilePath == "" {
		err := errors.New("spec for docker.read-file requires a file path within config")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	producer := d.producers.ReadFile(spec.Config.ContainerID, spec.Config.FilePath)

	if spec.Config.ContainerName != "" {
		producer = d.producers.ReadFileByName(spec.Config.ContainerName, spec.Config.FilePath)
	}

	task := &plans.StreamSource{
		Producer:  producer,
		RawPath:   spec.Raw,
		JSONPath:  spec.JSON,
		HumanPath: spec.Human,
	}

	if spec.TimeoutSeconds != 0 {
		task.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
	}

	return []types.Task{task}
}
