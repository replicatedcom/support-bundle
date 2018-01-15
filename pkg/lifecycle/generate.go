package lifecycle

import "github.com/replicatedcom/support-bundle/pkg/types"

func GenerateTask(task *types.LifecycleTask) Task {
	return func() error {
		return nil
	}
}
