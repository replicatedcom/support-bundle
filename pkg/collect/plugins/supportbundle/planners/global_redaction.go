package planners

import (
	"context"
	"path/filepath"

	"github.com/replicatedcom/support-bundle/pkg/collect/plans"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func (s *SupportBundle) GlobalRedaction(spec types.Spec) []types.Task {
	scrubs := spec.GlobalRedaction.Scrubs

	for _, scrub := range scrubs {
		scrubber, err := plans.RawScrubber(&types.Scrub{Regex: scrub.Regex, Replace: scrub.Replace})
		if err != nil {
			task := plans.PreparedError(err, spec)
			return []types.Task{task}
		}

		// GlobalScrubbers is used within StreamsSource to scrub all files
		// any task-specific scrubbers will be called first, followed by the global scrubbers in order of insertion
		plans.GlobalScrubbersLock.Lock()
		plans.GlobalScrubbers = append(plans.GlobalScrubbers, scrubber)
		plans.GlobalScrubbersLock.Unlock()
	}

	scubbersTask := plans.StructuredSource{
		JSONPath: filepath.Join(spec.Shared().OutputDir, "scrubs.json"),
		Producer: func(ctx context.Context) (interface{}, error) { return scrubs, nil },
	}

	return []types.Task{&scubbersTask}
}
