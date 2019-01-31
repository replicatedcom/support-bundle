package planners

import (
	"github.com/replicatedcom/support-bundle/pkg/collect/plans"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func (s *SupportBundle) GlobalRedaction(spec types.Spec) []types.Task {
	scrubs := spec.GlobalRedaction.Scrubs

	for _, scrub := range scrubs {
		scrubber, err := plans.RawScrubber(&scrub)
		if err != nil {
			task := plans.PreparedError(err, spec)
			return []types.Task{task}
		}

		plans.GlobalScrubbersLock.Lock()
		plans.GlobalScrubbers = append(plans.GlobalScrubbers, scrubber)
		plans.GlobalScrubbersLock.Unlock()
	}

	return []types.Task{}
}
