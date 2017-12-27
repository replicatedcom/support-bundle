package planners

import (
	"github.com/pkg/errors"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (c *Core) HTTPRequest(spec types.Spec) []types.Task {
	var err error
	if spec.CoreHTTPRequest == nil {
		err = errors.New("spec for os.http-request options required")
	} else if spec.CoreHTTPRequest.URL == "" {
		err = errors.New("spec for os.http-request url required")
	}
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	task := plans.StreamsSource{
		Producer: c.producers.HTTPRequest(*spec.CoreHTTPRequest),
	}
	task, err = plans.SetCommonFieldsStreamsSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
