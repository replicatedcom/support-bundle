package planners

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
	jww "github.com/spf13/jwalterweatherman"
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

	task := plans.StreamSource{
		RawPath: filepath.Join(spec.OutputDir, "logs.raw"),
	}
	if j.inContainer {
		jww.DEBUG.Println("journald.logs in container")
		task.Producer = j.producers.DockerLogs(*spec.JournaldLogs)
	} else {
		task.Producer = j.producers.Logs(*spec.JournaldLogs)
	}
	task, err = plans.SetCommonFieldsStreamSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
