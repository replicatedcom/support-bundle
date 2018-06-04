package bundle

import (
	"context"
	"encoding/json"

	"github.com/replicatedcom/support-bundle/pkg/types"
	jww "github.com/spf13/jwalterweatherman"
)

func Exec(ctx context.Context, rootDir string, tasks []types.Task) []*types.Result {
	results := make(chan []*types.Result)

	for _, task := range tasks {
		go func(task types.Task) {
			select {
			case results <- task.Exec(ctx, rootDir):
			case <-ctx.Done():
				b, _ := json.Marshal(task.GetSpec())
				jww.WARN.Println("Task failed to complete before context was canceled:", string(b))
			}
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
