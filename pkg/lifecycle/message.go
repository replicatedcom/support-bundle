package lifecycle

import (
	"fmt"

	"github.com/replicatedcom/support-bundle/pkg/types"
)

func MessageTask(task *types.LifecycleTask) Task {
	return func(l *Lifecycle) (bool, error) {
		fmt.Println(task.Message.Contents)
		return true, nil
	}
}
