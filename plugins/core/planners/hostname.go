package planners

import (
	"github.com/replicatedcom/support-bundle/plans"
	"github.com/replicatedcom/support-bundle/plugins/core/producers"
	"github.com/replicatedcom/support-bundle/types"
)

func Hostname(spec types.Spec) []types.Task {
	task := &plans.ByteSource{
		Producer:  producers.ReadCommand("hostname"),
		RawPath:   spec.Raw,
		JSONPath:  spec.JSON,
		HumanPath: spec.Human,
	}

	return []types.Task{task}
}
