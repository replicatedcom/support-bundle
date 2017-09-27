package bundle

import (
	"sync"
	"time"

	"github.com/replicatedcom/support-bundle/metrics"
	"github.com/replicatedcom/support-bundle/types"
)

func systemLogFiles(dataCh chan types.Data, completeCh chan bool, resultsCh chan types.Result, timeout time.Duration) error {
	completeCh <- true
	return nil
}

func systemMetrics(dataCh chan types.Data, completeCh chan bool, resultsCh chan types.Result, timeout time.Duration) error {
	var wg sync.WaitGroup

	itemCompleteCh := make(chan bool)

	var tasks = []Task{
		Task{
			Description: "System Uptime",
			ExecFunc:    metrics.Uptime,
			Timeout:     time.Duration(time.Second * 15),
		},

		Task{
			Description: "System Load Average",
			ExecFunc:    metrics.LoadAvg,
			Timeout:     time.Duration(time.Second * 15),
		},

		Task{
			Description: "System Hostname",
			ExecFunc:    metrics.Hostname,
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
		_ = task.ExecFunc(dataCh, itemCompleteCh, resultsCh, task.Timeout)
	}

	completeCh <- true
	return nil
}
