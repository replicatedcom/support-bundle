package plans

import (
	"context"
	"errors"
	"io"
	"os"
	"time"

	"github.com/replicatedcom/support-bundle/types"
)

// StreamSource is a Task that gets its data as an io.Reader
type StreamSource struct {
	// Producer provides the seed data for this task as an io.Reader
	Producer func(context.Context) (io.Reader, error)
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

	rawResult := &types.Result{}
	jsonResult := &types.Result{}
	humanResult := &types.Result{}

	results := []*types.Result{}
	if raw {
		results = append(results, rawResult)
	}
	if jsonify {
		results = append(results, jsonResult)
	}
	if human {
		results = append(results, humanResult)
	}

	if task.Producer == nil {
		err := errors.New("no data source defined for task")
		return resultsWithErr(err, results)
	}

	if !(raw || jsonify || human) {
		return results
	}

	useCtx := ctx
	if task.Timeout != 0 {
		var cancel context.CancelFunc
		useCtx, cancel = context.WithTimeout(useCtx, task.Timeout)
		defer cancel()
	}

	data, err := task.Producer(useCtx)
	if err != nil {
		return resultsWithErr(err, results)
	}
	if closer, ok := data.(io.Closer); ok {
		defer closeLogErr(closer)
	}

	// first write to one file
	if raw {
		ioCopyContext(ctx, rootDir, task.RawPath, data, rawResult)
	} else if jsonify {
		ioCopyContext(ctx, rootDir, task.JSONPath, data, jsonResult)
	} else if human {
		ioCopyContext(ctx, rootDir, task.HumanPath, data, humanResult)
	}

	// then link to any other requested paths
	if raw && jsonify {
		os.Link(task.RawPath, task.HumanPath)
	}
	if raw && human {
		os.Link(task.RawPath, task.HumanPath)
	}
	if jsonify && human {
		os.Link(task.JSONPath, task.HumanPath)
	}

	return results
}
