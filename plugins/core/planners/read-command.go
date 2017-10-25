package planners

import (
	"errors"

	"github.com/replicatedcom/support-bundle/plans"
	"github.com/replicatedcom/support-bundle/plugins/core/producers"
	"github.com/replicatedcom/support-bundle/types"
)

func parseCommandConfig(src interface{}) (string, []string) {
	command := ""
	var args []string

	m, ok := src.(map[interface{}]interface{})
	if !ok {
		return command, args
	}
	for k, v := range m {
		if key, ok := k.(string); ok {
			switch key {
			case "command":
				if val, ok := v.(string); ok {
					command = val
				}
			case "args":
				if val, ok := v.([]interface{}); ok {
					for _, argInterface := range val {
						if arg, ok := argInterface.(string); ok {
							args = append(args, arg)
						}
					}
				}
			}
		}
	}
	return command, args
}

func ReadCommand(spec types.Spec) []types.Task {
	command, args := parseCommandConfig(spec.Config)
	if command == "" {
		err := errors.New("spec requires a command within config")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	task := &plans.ByteSource{
		Producer:  producers.ReadCommand(command, args...),
		RawPath:   spec.Raw,
		JSONPath:  spec.JSON,
		HumanPath: spec.Human,
	}

	return []types.Task{task}
}
