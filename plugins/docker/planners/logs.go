package planners

import (
	"errors"

	"github.com/replicatedcom/support-bundle/plans"
	"github.com/replicatedcom/support-bundle/types"
)

func parseContainerConfig(src interface{}) types.ContainerConfig {
	config := types.ContainerConfig{}

	m, ok := src.(map[interface{}]interface{})
	if !ok {
		return config
	}
	for k, v := range m {
		if key, ok := k.(string); ok {
			switch key {
			case "container_id":
				if val, ok := v.(string); ok {
					config.ContainerID = val
				}
			}
		}
	}
	return config
}

func (d *Docker) Logs(spec types.Spec) []types.Task {
	config := parseContainerConfig(spec.Config)
	if config.ContainerID == "" {
		err := errors.New("spec requires container config")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	task := &plans.StreamSource{
		Producer:  d.producers.Logs(config.ContainerID),
		RawPath:   spec.Raw,
		JSONPath:  spec.JSON,
		HumanPath: spec.Human,
	}

	return []types.Task{task}
}
