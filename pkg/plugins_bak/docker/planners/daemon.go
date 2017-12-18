package planners

import (
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

// path returns "" if dir is empty, otherwise returns the joined pathnme
func maybePath(dir, filename string) string {
	if dir == "" {
		return ""
	}
	return filepath.Join(dir, filename)
}

// Daemon generates tasks to collect general information from Docker. The paths
// in the spec are interpreted as directories.
func (d *Docker) Daemon(spec types.Spec) []types.Task {
	scrubber, err := plans.RawScrubber(spec.Config.Scrub)
	if err != nil {
		err = errors.Wrap(err, "create scrubber for docker.daemon")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	info := &plans.StructuredSource{
		Producer:    d.producers.Info,
		RawScrubber: scrubber,
		RawPath:     maybePath(spec.Raw, "docker_info"),
		JSONPath:    maybePath(spec.JSON, "docker_info.json"),
		HumanPath:   maybePath(spec.Human, "docker_info"),
	}

	ps := &plans.StructuredSource{
		Producer:    d.producers.PSAll,
		RawScrubber: scrubber,
		RawPath:     maybePath(spec.Raw, "docker_ps_all"),
		JSONPath:    maybePath(spec.JSON, "docker_ps_all.json"),
		HumanPath:   maybePath(spec.Human, "docker_ps_all"),
	}

	images := &plans.StructuredSource{
		Producer:    d.producers.ImagesAll,
		RawScrubber: scrubber,
		RawPath:     maybePath(spec.Raw, "docker_images_all"),
		JSONPath:    maybePath(spec.JSON, "docker_images_all.json"),
		HumanPath:   maybePath(spec.Human, "docker_images_all"),
	}

	if spec.TimeoutSeconds != 0 {
		info.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
		ps.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
		images.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
	}

	return []types.Task{
		info,
		ps,
		images,
	}
}
