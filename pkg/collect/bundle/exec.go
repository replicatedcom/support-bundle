package bundle

import (
	"context"
	"encoding/json"
	"math"
	"runtime"
	"time"

	"github.com/replicatedcom/support-bundle/pkg/collect/plans"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
	jww "github.com/spf13/jwalterweatherman"
)

const (
	TaskConcurrencyMax     = 4
	DeferredTimeoutDefault = 5 * time.Second
)

type ExecOptions struct {
	DeferredTimeout time.Duration
}

type Option func(*ExecOptions)

func Exec(ctx context.Context, rootDir string, tasks []types.Task, opts ...Option) []*types.Result {
	var startTasks, deferredTasks []types.Task

	for _, task := range tasks {
		if task.GetSpec().Shared().Defer {
			deferredTasks = append(deferredTasks, task)
		} else {
			startTasks = append(startTasks, task)
		}
	}

	concurrency := int(math.Min(float64(runtime.NumCPU()), float64(TaskConcurrencyMax)))

	jww.DEBUG.Printf("Running tasks with a concurrency of %d", concurrency)

	var results []*types.Result

	jww.DEBUG.Println("Running tasks...")

	results = append(results, exec(ctx, rootDir, startTasks, concurrency)...)

	jww.DEBUG.Println("Running tasks complete")

	options := ExecOptions{
		DeferredTimeout: DeferredTimeoutDefault,
	}
	for _, opt := range opts {
		opt(&options)
	}

	// hardcode short deferred task timeout of 5 seconds
	deferredCtx, cancel := context.WithTimeout(context.Background(), options.DeferredTimeout)
	defer cancel()

	jww.DEBUG.Println("Running deferred tasks...")

	results = append(results, exec(deferredCtx, rootDir, deferredTasks, concurrency)...)

	jww.DEBUG.Println("Running deferred tasks complete")

	return results
}

func exec(ctx context.Context, rootDir string, tasks []types.Task, concurrency int) []*types.Result {
	results := make(chan []*types.Result)

	// allow for tasks to time themselves out before we cancel
	taskCtx, cancel := contextTimeoutSub(ctx, time.Second)
	defer cancel()

	lock := make(chan bool, concurrency)
	for i := 0; i < concurrency; i++ {
		lock <- true
	}

	for _, task := range tasks {
		result := make(chan []*types.Result, 1)

		go func(task types.Task) {
			<-lock
			result <- task.Exec(taskCtx, rootDir)
			lock <- true
		}(task)

		go func(task types.Task) {
			select {
			case r := <-result:
				results <- r
			case <-ctx.Done():
				b, _ := json.Marshal(task.GetSpec())
				jww.WARN.Println("Task failed to complete before context was canceled:", string(b))

				// TODO: the task can probably send a better error here
				results <- plans.PreparedError(ctx.Err(), task.GetSpec()).Exec(ctx, rootDir)
			}
		}(task)
	}

	var accm []*types.Result
	for i := 0; i < len(tasks); i++ {
		accm = append(accm, <-results...)
	}
	return accm
}

func contextTimeoutSub(ctx context.Context, less time.Duration) (context.Context, context.CancelFunc) {
	if deadline, ok := ctx.Deadline(); ok {
		dur := time.Until(deadline)
		// the minimum timeout is half the remaining duration
		if dur-less < dur/2 {
			return context.WithTimeout(ctx, dur/2)
		}
		return context.WithTimeout(ctx, dur-less)
	}
	return context.WithCancel(ctx)
}
