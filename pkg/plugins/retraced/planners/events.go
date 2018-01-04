package planners

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (r *Retraced) Events(spec types.Spec) []types.Task {
	var err error
	if spec.RetracedEvents == nil {
		err = errors.New("spec for retraced.events required")
	} else if spec.RetracedEvents.APIEndpoint == "" { // require an endpoint, I don't think defaulting to SaaS Retraced makes sense
		err = errors.New("spec for retraced.events api_endpoint required")
	} else if spec.RetracedEvents.APIToken == "" {
		err = errors.New("spec for retraced.events api_token required")
	} else if spec.RetracedEvents.ProjectID == "" {
		err = errors.New("spec for retraced.events project_id required")
	}
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	client, err := getClient((*spec.RetracedEvents).RetracedAPIClientOptions)
	if err != nil {
		err = errors.Wrap(err, "create retraced client")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	task := plans.StreamSource{
		Producer: r.producers.Events(client, *spec.RetracedEvents),
		RawPath:  filepath.Join(spec.OutputDir, "audit_events.csv"),
	}
	task, err = plans.SetCommonFieldsStreamSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
