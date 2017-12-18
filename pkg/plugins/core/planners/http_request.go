package planners

import (
	"github.com/pkg/errors"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/plugins/core/producers"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func HTTPRequest(spec types.Spec) []types.Task {
	if spec.CoreHTTPRequest == nil {
		err := errors.New("spec for os.http-request options required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	if spec.CoreHTTPRequest.URL == "" {
		err := errors.New("spec for os.http-request url required")
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}

	task := plans.StreamsSource{
		Spec:     spec,
		Producer: producers.HTTPRequest(*spec.CoreHTTPRequest),
	}
	var err error
	task, err = setCommonFieldsStreamsSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
