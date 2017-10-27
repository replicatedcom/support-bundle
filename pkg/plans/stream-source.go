package plans

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"time"

	"github.com/replicatedcom/support-bundle/pkg/types"
)

// StreamSource is a Task that gets its data as an io.Reader
type StreamSource struct {
	// Producer provides the seed data for this task as an io.Reader
	Producer func(context.Context) (io.Reader, error)
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

func (task *StreamSource) Exec(ctx context.Context, rootDir string) []*types.Result {
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

	data, err := task.Producer(ctx)
	if err != nil {
		return resultsWithErr(err, results)
	}
	if closer, ok := data.(io.Closer); ok {
		defer closeLogErr(closer)
	}

	if task.RawScrubber != nil {
		scrubbedReader, scrubbedWriter := io.Pipe()
		go func(readFrom io.Reader, writeTo *io.PipeWriter) {
			lineScanner := bufio.NewScanner(readFrom)
			byteWriter := bufio.NewWriter(writeTo)
			for lineScanner.Scan() {
				line := lineScanner.Bytes()
				line = task.RawScrubber(line)

				// TODO: actually deal with errors besides 'buffer full'
				n, _ := byteWriter.Write(line)
				for n < len(line) {
					line = line[n:]
					n, _ = byteWriter.Write(line)
				}
			}
			writeTo.Close()
		}(data, scrubbedWriter)
		data = scrubbedReader
	}

	rawResult := types.Result{}
	jsonResult := types.Result{}
	humanResult := types.Result{}

	// first write to one file
	if raw {
		ioCopyContext(ctx, rootDir, task.RawPath, data, &rawResult)
	} else if jsonify {
		ioCopyContext(ctx, rootDir, task.JSONPath, data, &jsonResult)
	} else if human {
		ioCopyContext(ctx, rootDir, task.HumanPath, data, &humanResult)
	}

	// then link to any other requested paths
	if raw && jsonify {
		jsonResult = rawResult
		if rawResult.Path != "" {
			os.Link(task.RawPath, task.JSONPath)
			jsonResult.Path = task.JSONPath
		}
	}
	if raw && human {
		humanResult = rawResult
		if rawResult.Path != "" {
			os.Link(task.RawPath, task.HumanPath)
			humanResult.Path = task.HumanPath
		}
	}
	if jsonify && human {
		humanResult = jsonResult
		if jsonResult.Path != "" {
			os.Link(task.JSONPath, task.HumanPath)
			humanResult.Path = task.HumanPath
		}
	}

	if raw {
		results = append(results, &rawResult)
	}
	if jsonify {
		results = append(results, &jsonResult)
	}
	if human {
		results = append(results, &humanResult)
	}

	return results
}
