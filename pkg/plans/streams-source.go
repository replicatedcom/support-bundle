package plans

import (
	"context"
	"errors"
	"io"
	"os"
	"sync"
	"time"

	"github.com/replicatedcom/support-bundle/pkg/types"
)

// StreamsSource is a Task that gets its data as an io.Reader
type StreamsSource struct {
	// Producer provides the seed data for this task as an io.Reader array
	// Names of the sources are provided as a string array
	Producer func(context.Context) (map[string]io.Reader, error)
	// Template, if defined, renders structured data in a human-readable format
	Template string
	// If RawPath is defined it will get a copy of the data
	RawPath string
	// If JSONPath is defined it will get a copy of the data
	JSONPath string
	// If HumanPath is defined it will get a copy of the data
	HumanPath string
	// If Timeout is defined, it will be used rather than the context provided
	// to Exec.
	Timeout time.Duration
}

func (task *StreamsSource) Exec(ctx context.Context, rootDir string) []*types.Result {
	raw := task.RawPath != ""
	jsonify := task.JSONPath != ""
	human := task.HumanPath != ""

	results := []*types.Result{}

	if !(raw || jsonify || human) {
		return results
	}

	if task.Producer == nil {
		err := errors.New("no data source defined for task")
		if raw {
			results = append(results, &types.Result{Description: task.RawPath})
		}
		if jsonify {
			results = append(results, &types.Result{Description: task.JSONPath})
		}
		if human {
			results = append(results, &types.Result{Description: task.HumanPath})
		}
		return resultsWithErr(err, results)
	}

	if task.Timeout != 0 {
		ctx, _ = context.WithTimeout(ctx, task.Timeout)
	}

	readers, err := task.Producer(ctx)
	if err != nil {
		return resultsWithErr(err, results)
	}

	var resultsMut sync.Mutex
	var readerGroup sync.WaitGroup
	readerGroup.Add(len(readers))

	for name, reader := range readers {
		go func(name string, reader io.Reader) {
			if closer, ok := reader.(io.Closer); ok {
				defer closeLogErr(closer)
			}

			rawResult := types.Result{}
			jsonResult := types.Result{}
			humanResult := types.Result{}

			// first write to one file
			if raw {
				ioCopyContext(ctx, rootDir, task.RawPath+name, reader, &rawResult)
			} else if jsonify {
				ioCopyContext(ctx, rootDir, task.JSONPath+name, reader, &jsonResult)
			} else if human {
				ioCopyContext(ctx, rootDir, task.HumanPath+name, reader, &humanResult)
			}

			// then link to any other requested paths
			if raw && jsonify {
				jsonResult = rawResult
				if rawResult.Path != "" {
					os.Link(task.RawPath+name, task.JSONPath+name)
					jsonResult.Path = task.JSONPath + name
				}
			}
			if raw && human {
				humanResult = rawResult
				if rawResult.Path != "" {
					os.Link(task.RawPath+name, task.HumanPath+name)
					humanResult.Path = task.HumanPath + name
				}
			}
			if jsonify && human {
				humanResult = jsonResult
				if jsonResult.Path != "" {
					os.Link(task.JSONPath+name, task.HumanPath+name)
					humanResult.Path = task.HumanPath + name
				}
			}

			resultsMut.Lock()
			if raw {
				results = append(results, &rawResult)
			}
			if jsonify {
				results = append(results, &jsonResult)
			}
			if human {
				results = append(results, &humanResult)
			}
			resultsMut.Unlock()

			readerGroup.Done()
		}(name, reader)
	}

	readerGroup.Wait()
	return results
}