package planners

import (
	"time"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

// ServiceLogs gets the logs of a service, as with `docker service logs SERVICENAME`
func (d *Docker) ServiceLogs(spec types.Spec) []types.Task {
	scrubber, err := plans.RawScrubber(spec.Config.Scrub)
	if err != nil {
		err = errors.Wrap(err, "create scrubber for docker.swarm")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	if spec.DockerSwarmOptions.ServiceName == "" {
		err := errors.New("spec for docker.service-logs requires a service name")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	logs := &plans.StreamsSource{
		Producer:    d.producers.ServiceLogs(spec.DockerSwarmOptions.ServiceName),
		RawScrubber: scrubber,
		RawPath:     spec.Raw,
		JSONPath:    spec.JSON,
		HumanPath:   spec.Human,
	}

	if spec.TimeoutSeconds != 0 {
		logs.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
	}

	return []types.Task{
		logs,
	}
}
