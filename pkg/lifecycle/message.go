package lifecycle

import (
	"fmt"

	"github.com/replicatedcom/support-bundle/pkg/types"
)

func MessageTask(task *types.LifecycleTask) Task {
	return func() error {
		fmt.Println(task.Message.Contents)
		return nil
	}
}
