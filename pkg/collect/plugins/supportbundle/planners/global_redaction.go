package planners

import (
	"context"
	"path/filepath"

	"github.com/replicatedcom/support-bundle/pkg/collect/plans"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func (s *SupportBundle) GlobalRedaction(spec types.Spec) (tasks []types.Task) {
	files := spec.GlobalRedaction.Files

	if len(files) > 0 {
		fileRedactor, err := plans.FileRedactor(files)
		if err != nil {
			tasks = append(tasks, plans.PreparedError(err, spec))
			return
		}
		plans.GlobalLock.Lock()
		plans.GlobalFileRedactors = append(plans.GlobalFileRedactors, fileRedactor)
		plans.GlobalLock.Unlock()
	}

	tasks = append(tasks, &plans.StructuredSource{
		JSONPath: filepath.Join(spec.Shared().OutputDir, "file_redactions.json"),
		Producer: func(ctx context.Context) (interface{}, error) { return files, nil },
	})

	scrubs := spec.GlobalRedaction.Scrubs

	for _, scrub := range scrubs {
		scrubber, err := plans.RawScrubber(&types.Scrub{Regex: scrub.Regex, Replace: scrub.Replace})
		if err != nil {
			tasks = append(tasks, plans.PreparedError(err, spec))
			return
		}

		// GlobalScrubbers is used within StreamsSource to scrub all files
		// any task-specific scrubbers will be called first, followed by the global scrubbers in order of insertion
		plans.GlobalLock.Lock()
		plans.GlobalScrubbers = append(plans.GlobalScrubbers, scrubber)
		plans.GlobalLock.Unlock()
	}

	tasks = append(tasks, &plans.StructuredSource{
		JSONPath: filepath.Join(spec.Shared().OutputDir, "scrubs.json"),
		Producer: func(ctx context.Context) (interface{}, error) { return scrubs, nil },
	})

	return
}
