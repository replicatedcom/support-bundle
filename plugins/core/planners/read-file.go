package planners

import (
	"errors"

	"github.com/replicatedcom/support-bundle/plans"
	"github.com/replicatedcom/support-bundle/plugins/core/producers"
	"github.com/replicatedcom/support-bundle/types"
)

func parseFileConfig(src interface{}) string {
	filename := ""

	m, ok := src.(map[interface{}]interface{})
	if !ok {
		return filename
	}
	for k, v := range m {
		if key, ok := k.(string); ok {
			switch key {
			case "file_path":
				if val, ok := v.(string); ok {
					filename = val
				}
			}
		}
	}
	return filename
}

func ReadFile(spec types.Spec) []types.Task {
	filename := parseFileConfig(spec.Config)
	if filename == "" {
		err := errors.New("spec requires a filename within config")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	task := &plans.ByteSource{
		Producer:  producers.ReadFile(filename),
		RawPath:   spec.Raw,
		JSONPath:  spec.JSON,
		HumanPath: spec.Human,
	}

	return []types.Task{task}
}
