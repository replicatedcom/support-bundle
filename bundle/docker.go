package bundle

import (
	"sync"
	"time"

	"github.com/replicatedcom/support-bundle/metrics"
	"github.com/replicatedcom/support-bundle/types"
)

func dockerMetrics(dataCh chan types.Data, completeCh chan bool, resultsCh chan types.Result, timeout time.Duration, args []string) error {
	var wg sync.WaitGroup

	itemCompleteCh := make(chan bool)

	var tasks = []Task{
		Task{
			Description: "Docker ps",
			ExecFunc:    metrics.Dockerps,
			Timeout:     time.Duration(time.Second * 15),
		},

		Task{
			Description: "Docker info",
			ExecFunc:    metrics.DockerInfo,
			Timeout:     time.Duration(time.Second * 15),
		},
	}
	wg.Add(len(tasks))

	go func() {
		for {
			select {
			case <-itemCompleteCh:
				wg.Done()
				break
			}
		}
	}()

	for _, task := range tasks {
		_ = task.ExecFunc(dataCh, itemCompleteCh, resultsCh, task.Timeout, task.Args)
	}

	completeCh <- true
	return nil
}
