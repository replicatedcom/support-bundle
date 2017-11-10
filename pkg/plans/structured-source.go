package plans

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/replicatedcom/support-bundle/pkg/types"
)

// StructuredSource is a Task that gets its data as a structured object ready to
// be jsonified or used as a context in a human template.
type StructuredSource struct {
	// Producer provides the seed data for this task
	Producer func(context.Context) (interface{}, error)
	// RawScrubber, if defined, rewrites the raw data to to remove sensitive data
	RawScrubber func([]byte) []byte
	// Template, if defined, renders the task's data in a human-readable format
	Template string
	// RawPath, if defined, gets the tasks's data as json.
	RawPath string
	// JSONPath, if defined, gets a jsonified copy of the task's data.
	JSONPath string
	// If HumanPath is defined and a Template is defined, it will get the
	// output of the template rendered with the task's data as context. If
	// HumanPath is defined and a Template is not, it will get the data as
	// YAML.
	HumanPath string
	// If Timeout is defined, it will be used rather than the context provided
	// to Exec.
	Timeout time.Duration
}

func (task *StructuredSource) Exec(ctx context.Context, rootDir string) []*types.Result {
	rawScrubber := task.RawScrubber != nil

	raw := task.RawPath != ""
	jsonify := task.JSONPath != ""
	human := task.HumanPath != ""
	humanTemplated := human && task.Template != ""
	humanYAML := human && !humanTemplated

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

	rawResult := &types.Result{}
	jsonResult := &types.Result{}
	humanResult := &types.Result{}

	if raw {
		results = append(results, rawResult)
	}
	if jsonify {
		results = append(results, jsonResult)
	}
	if human {
		results = append(results, humanResult)
	}

	if task.Timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, task.Timeout)
		defer cancel()
	}

	data, err := task.Producer(ctx)
	if err != nil {
		return resultsWithErr(err, results)
	}

	if jsonify {
		if rawScrubber {
			buf := bytes.NewBuffer(nil)
			n, err := writeJSON(ctx, buf, data)
			if err != nil {
				jsonResult.Error = err
			}
			if n > 0 {
				data := task.RawScrubber(buf.Bytes())
				writeResultBytes(ctx, rootDir, task.JSONPath, jsonResult, data)
			}
		} else {
			writeResultJSON(ctx, rootDir, task.JSONPath, jsonResult, data)
		}
	}

	if raw {
		if rawScrubber {
			buf := bytes.NewBuffer(nil)
			n, err := writeJSON(ctx, buf, data)
			if err != nil {
				rawResult.Error = err
			}
			if n > 0 {
				data := task.RawScrubber(buf.Bytes())
				writeResultBytes(ctx, rootDir, task.RawPath, rawResult, data)
			}
		} else {
			writeResultJSON(ctx, rootDir, task.RawPath, rawResult, data)
		}
	}

	if humanTemplated {
		if rawScrubber {
			buf := bytes.NewBuffer(nil)
			n, err := writeTemplate(ctx, buf, task.Template, data)
			if err != nil {
				humanResult.Error = err
			}
			if n > 0 {
				data := task.RawScrubber(buf.Bytes())
				writeResultBytes(ctx, rootDir, task.HumanPath, humanResult, data)
			}
		} else {
			writeResultTemplate(ctx, rootDir, task.HumanPath, humanResult, task.Template, data)
		}
	}

	if humanYAML {
		if rawScrubber {
			buf := bytes.NewBuffer(nil)
			n, err := writeYAML(ctx, buf, data)
			if err != nil {
				humanResult.Error = err
			}
			if n > 0 {
				data := task.RawScrubber(buf.Bytes())
				writeResultBytes(ctx, rootDir, task.HumanPath, humanResult, data)
			}
		} else {
			writeResultYAML(ctx, rootDir, task.HumanPath, humanResult, data)
		}
	}

	return results
}
