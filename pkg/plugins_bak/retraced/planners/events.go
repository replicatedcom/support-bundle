package planners

import (
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/plugins/retraced/producers"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func Events(spec types.Spec) []types.Task {
	var err error
	if spec.RetracedEventsCommand == nil {
		err = errors.New("spec for retraced.events required")
	} else if spec.RetracedEventsCommand.APIEndpoint == "" { // require an endpoint, I don't think defaulting to SaaS Retraced makes sense
		err = errors.New("retraced.events spec missing api_endpoint")
	} else if spec.RetracedEventsCommand.APIToken == "" {
		err = errors.New("retraced.events spec missing api_token")
	} else if spec.RetracedEventsCommand.ProjectID == "" {
		err = errors.New("retraced.events spec missing project_id")
	}
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{
		&plans.StreamSource{
			Producer: producers.Events(spec),
			RawPath:  spec.Raw,
		},
	}
}
