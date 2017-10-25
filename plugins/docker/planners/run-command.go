package planners

import (
	"errors"

	"github.com/replicatedcom/support-bundle/plans"
	"github.com/replicatedcom/support-bundle/types"
)

func parseCommandConfig(src interface{}) (string, []string) {
	id := ""
	var command []string

	m, ok := src.(map[interface{}]interface{})
	if !ok {
		return id, command
	}
	for k, v := range m {
		if key, ok := k.(string); ok {
			switch key {
			case "container_id":
				if val, ok := v.(string); ok {
					id = val
				}
			case "command":
				if val, ok := v.([]interface{}); ok {
					for _, argInterface := range val {
						if arg, ok := argInterface.(string); ok {
							command = append(command, arg)
						}
					}
				}
			}
		}
	}
	return id, command
}

func (d *Docker) RunCommand(spec types.Spec) []types.Task {
	id, command := parseCommandConfig(spec.Config)
	if id == "" || len(command) == 0 {
		err := errors.New("spec requires a container ID and command within config")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	// task := &plans.StreamsSource{
	// 	Producer:  d.producers.RunCommand(id, command),
	// 	RawPath:   spec.Raw,
	// 	JSONPath:  spec.JSON,
	// 	HumanPath: spec.Human,
	// }

	err := errors.New("This task type not yet implemented")
	task := plans.PreparedError(err, spec)

	return []types.Task{task}
}
