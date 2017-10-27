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
	Producer func(context.Context) ([]io.Reader, []string, error)
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

	if task.Producer == nil {
		err := errors.New("no data source defined for task")
		return resultsWithErr(err, results)
	}

	if !(raw || jsonify || human) {
		return results
	}

	if task.Timeout != 0 {
		ctx, _ = context.WithTimeout(ctx, task.Timeout)
	}

	readers, names, err := task.Producer(ctx)
	if err != nil {
		return resultsWithErr(err, results)
	}

	var resultsMut sync.Mutex
	var readerGroup sync.WaitGroup
	readerGroup.Add(len(readers))

	for i, reader := range readers {
		go func(i int, reader io.Reader) {
			if closer, ok := reader.(io.Closer); ok {
				defer closeLogErr(closer)
			}

			rawResult := types.Result{}
			jsonResult := types.Result{}
			humanResult := types.Result{}

			// first write to one file
			if raw {
				ioCopyContext(ctx, rootDir, task.RawPath+names[i], reader, &rawResult)
			} else if jsonify {
				ioCopyContext(ctx, rootDir, task.JSONPath+names[i], reader, &jsonResult)
			} else if human {
				ioCopyContext(ctx, rootDir, task.HumanPath+names[i], reader, &humanResult)
			}

			// then link to any other requested paths
			if raw && jsonify {
				os.Link(task.RawPath+names[i], task.JSONPath+names[i])
				jsonResult = rawResult
				if jsonResult.Path != "" {
					jsonResult.Path = task.JSONPath + names[i]
				}
			}
			if raw && human {
				os.Link(task.RawPath+names[i], task.HumanPath+names[i])
				humanResult = rawResult
				if humanResult.Path != "" {
					humanResult.Path = task.HumanPath + names[i]
				}
			}
			if jsonify && human {
				os.Link(task.JSONPath+names[i], task.HumanPath+names[i])
				humanResult = jsonResult
				if humanResult.Path != "" {
					humanResult.Path = task.HumanPath + names[i]
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
		}(i, reader)
	}

	readerGroup.Wait()
	return results
}
