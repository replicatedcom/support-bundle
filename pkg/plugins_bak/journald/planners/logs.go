package planners

import (
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/plugins/journald/producers"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (j *Journald) Logs(spec types.Spec) []types.Task {
	var err error
	if spec.JournaldLogs == nil {
		err = errors.New("spec for journald.logs required")
	}
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	return []types.Task{
		&plans.ByteSource{
			Producer: producers.Logs(j.dockerClient, *spec.JournaldLogs),
			RawPath:  spec.Raw,
		},
	}
}
