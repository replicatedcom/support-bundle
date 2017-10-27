package plans

import (
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
	raw := task.RawPath != ""
	jsonify := task.JSONPath != ""
	human := task.HumanPath != ""
	humanTemplated := human && task.Template != ""
	humanYAML := human && !humanTemplated

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

	if task.Timeout != 0 {
		ctx, _ = context.WithTimeout(ctx, task.Timeout)
	}

	data, err := task.Producer(ctx)
	if err != nil {
		return resultsWithErr(err, results)
	}

	if jsonify {
		writeJSON(rootDir, task.JSONPath, data, jsonResult)
	}

	if raw {
		writeJSON(rootDir, task.RawPath, data, rawResult)
	}

	if humanTemplated {
		writeTemplate(rootDir, task.HumanPath, task.Template, data, humanResult)
	}

	if humanYAML {
		writeYAML(rootDir, task.HumanPath, data, humanResult)
	}

	return results
}
