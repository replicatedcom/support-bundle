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
	// RawScrubber, if defined, rewrites the raw data to to remove sensitive data
	RawScrubber func([]byte) []byte
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

	descriptionResults := []*types.Result{}
	if raw {
		descriptionResults = append(descriptionResults, &types.Result{Description: task.RawPath})
	}
	if jsonify {
		descriptionResults = append(descriptionResults, &types.Result{Description: task.JSONPath})
	}
	if human {
		descriptionResults = append(descriptionResults, &types.Result{Description: task.HumanPath})
	}

	if !(raw || jsonify || human) {
		return results
	}

	if task.Producer == nil {
		err := errors.New("no data source defined for task")

		return resultsWithErr(err, descriptionResults)
	}

	if task.Timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, task.Timeout)
		defer cancel()
	}

	readers, err := task.Producer(ctx)
	if err != nil {
		return resultsWithErr(err, descriptionResults)
	}

	var resultsMut sync.Mutex
	var readerGroup sync.WaitGroup
	readerGroup.Add(len(readers))

	for name, reader := range readers {
		go func(name string, reader io.Reader) {
			if closer, ok := reader.(io.Closer); ok {
				defer closeLogErr(closer)
			}

			if task.RawScrubber != nil {
				scrubbedReader, scrubbedWriter := io.Pipe()
				go filterStreams(reader, scrubbedWriter, task.RawScrubber)
				reader = scrubbedReader
			}

			rawResult := types.Result{}
			jsonResult := types.Result{}
			humanResult := types.Result{}

			// first write to one file
			if raw {
				writeResult(ctx, rootDir, task.RawPath+name, &rawResult, reader)
			} else if jsonify {
				writeResult(ctx, rootDir, task.JSONPath+name, &jsonResult, reader)
			} else if human {
				writeResult(ctx, rootDir, task.HumanPath+name, &humanResult, reader)
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
