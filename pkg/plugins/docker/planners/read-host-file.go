package planners

import (
	"time"

	"github.com/pkg/errors"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (d *Docker) ReadHostFile(spec types.Spec) []types.Task {
	if spec.Config.Image == "" {
		err := errors.New("spec for docker.read-host-file requires a image name within config")
		return []types.Task{plans.PreparedError(err, spec)}
	}

	if spec.Config.FilePath == "" {
		err := errors.New("spec for docker.read-host-file requires a file path within config")
		return []types.Task{plans.PreparedError(err, spec)}
	}

	scrubber, err := plans.RawScrubber(spec.Config.Scrub)
	if err != nil {
		err = errors.Wrap(err, "create scrubber for docker.read-host-file")
		return []types.Task{plans.PreparedError(err, spec)}
	}

	producer := d.producers.ReadHostFile(spec.Config.Image, spec.Config.FilePath)

	task := &plans.StreamSource{
		Producer:     producer,
		StreamFormat: "tar",
		RawScrubber:  scrubber,
		RawPath:      spec.Raw,
		JSONPath:     spec.JSON,
		HumanPath:    spec.Human,
	}

	if spec.TimeoutSeconds != 0 {
		task.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
	}

	return []types.Task{task}
}
