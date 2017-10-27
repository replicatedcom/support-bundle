package bundle

import (
	"context"

	"github.com/replicatedcom/support-bundle/pkg/types"
)

func Exec(ctx context.Context, rootDir string, tasks []types.Task) []*types.Result {
	results := make(chan []*types.Result)

	for _, task := range tasks {
		go func(task types.Task) {
			results <- task.Exec(ctx, rootDir)
		}(task)
	}

	pending := len(tasks)
	var accm []*types.Result

	for {
		select {
		case r := <-results:
			accm = append(accm, r...)
			pending--
			if pending == 0 {
				return accm
			}
		case <-ctx.Done():
			return accm
		}
	}
}
