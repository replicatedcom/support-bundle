package planners

import (
	"errors"

	"github.com/replicatedcom/support-bundle/plans"
	"github.com/replicatedcom/support-bundle/types"
)

func parseFileConfig(src interface{}) (string, string) {
	id, filename := "", ""

	m, ok := src.(map[interface{}]interface{})
	if !ok {
		return id, filename
	}
	for k, v := range m {
		if key, ok := k.(string); ok {
			switch key {
			case "container_id":
				if val, ok := v.(string); ok {
					id = val
				}
			case "file_path":
				if val, ok := v.(string); ok {
					filename = val
				}
			}
		}
	}
	return id, filename
}

func (d *Docker) ReadFile(spec types.Spec) []types.Task {
	id, filename := parseFileConfig(spec.Config)
	if id == "" || filename == "" {
		err := errors.New("spec requires a container ID and filename within config")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	task := &plans.StreamSource{
		Producer:  d.producers.ReadFile(id, filename),
		RawPath:   spec.Raw,
		JSONPath:  spec.JSON,
		HumanPath: spec.Human,
	}

	return []types.Task{task}
}
