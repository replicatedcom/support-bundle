package lifecycle

import "github.com/replicatedcom/support-bundle/pkg/types"

func PromptTask(task *types.LifecycleTask) Task {
	return func() error {
		return nil
	}
}
