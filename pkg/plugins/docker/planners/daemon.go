package planners

import (
	"path/filepath"
	"time"

	"github.com/replicatedcom/support-bundle/pkg/types"
	"github.com/replicatedcom/support-bundle/pkg/plans"
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
	info := &plans.StructuredSource{
		Producer:  d.producers.Info,
		RawPath:   maybePath(spec.Raw, "docker_info"),
		JSONPath:  maybePath(spec.JSON, "docker_info.json"),
		HumanPath: maybePath(spec.Human, "docker_info"),
	}

	ps := &plans.StructuredSource{
		Producer:  d.producers.PSAll,
		RawPath:   maybePath(spec.Raw, "docker_ps_all"),
		JSONPath:  maybePath(spec.JSON, "docker_ps_all.json"),
		HumanPath: maybePath(spec.Human, "docker_ps_all"),
	}

	if spec.TimeoutSeconds != 0 {
		info.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
		ps.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
	}

	return []types.Task{
		info,
		ps,
	}
}
