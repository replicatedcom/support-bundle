package lifecycle

import "github.com/replicatedcom/support-bundle/pkg/types"

func UploadTask(task *types.LifecycleTask) Task {
	return func() error {
		return nil
	}
}
