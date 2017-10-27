package planners

import (
	"regexp"
	"time"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/plugins/core/producers"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func ReadFile(spec types.Spec) []types.Task {
	if spec.Config.FilePath == "" {
		err := errors.New("spec for core.read-file requires a filename within config")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	scrubber, err := rawScrubber(spec.Config.Scrub)
	if err != nil {
		err := errors.New("spec for core.read-file has invalid scrubber spec")
		task := plans.PreparedError(err, spec)

		return []types.Task{task}
	}

	task := &plans.ByteSource{
		Producer:    producers.ReadFile(spec.Config.FilePath),
		RawScrubber: scrubber,
		RawPath:     spec.Raw,
		JSONPath:    spec.JSON,
		HumanPath:   spec.Human,
	}

	if spec.TimeoutSeconds != 0 {
		task.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
	}

	return []types.Task{task}
}

func rawScrubber(scrubSpec types.Scrub) (types.BytesScrubber, error) {
	if scrubSpec.Regex == "" {
		return nil, nil
	}

	regex, err := regexp.Compile(scrubSpec.Regex)
	if err != nil {
		return nil, errors.Wrapf(err, "parse regex %s", scrubSpec.Regex)
	}

	return func(in []byte) []byte {
		return regex.ReplaceAll(in, []byte(scrubSpec.Replace))
	}, nil

}
