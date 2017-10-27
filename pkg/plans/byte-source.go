package plans

import (
	"context"
	"errors"
	"time"

	"github.com/replicatedcom/support-bundle/pkg/types"
)

// ByteSource is a Task that gets its data source as a []byte.
type ByteSource struct {
	// Producer provides the seed data for this task
	Producer func(context.Context) ([]byte, error)
	// RawScrubber, if defined, rewrites the raw data to to remove sensitive data
	RawScrubber func([]byte) []byte
	// Parser, if defined, structures the raw data for json and human sinks
	Parser func([]byte) (interface{}, error)
	// StructuredScrubber, if defined, rewrites the structured data to remove sensitive data
	//StructedScrubber func(interface{}) (interface{}, error)
	// Template, if defined, renders structured data in a human-readable format
	Template string
	// If RawPath is defined it will get a copy of the raw data []byte
	RawPath string
	// If JSONPath is defined and Parser is defined, it will get a jsonified
	// copy of the strucutred data. If JSONPath is defined and Parser is not,
	// it will get a copy of the raw data.
	JSONPath string
	// If HumanPath is defined and Parser and Template are defined, it will get
	// the output of the template rendered with the structure context. If
	// HumanPath and Parser are defined but Template is not, it will get a YAML
	// copy of the structured data. If HumanPath is defined and Parser is not,
	// it will get a copy of the raw data.
	HumanPath string
	// If Timeout is defined, it will be used rather than the context provided
	// to Exec.
	Timeout time.Duration
}

func (task *ByteSource) Exec(ctx context.Context, rootDir string) []*types.Result {
	parser := task.Parser != nil
	rawScrubber := task.RawScrubber != nil

	raw := task.RawPath != ""

	jsonify := task.JSONPath != ""
	jsonParsed := jsonify && parser
	jsonRaw := jsonify && !jsonParsed

	human := task.HumanPath != ""
	humanTemplated := human && parser && task.Template != ""
	humanYAML := human && parser && !humanTemplated
	humanRaw := human && !humanTemplated && !humanYAML

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

	if rawScrubber {
		data = task.RawScrubber(data)
	}

	if raw {
		write(rootDir, task.RawPath, data, rawResult)
	}

	var structured interface{}
	if parser {
		structured, err = task.Parser(data)
		if err != nil {
			jsonResult.Error = err

			if humanTemplated || humanYAML {
				humanResult.Error = err
				// nothing left to do
				return results
			}
		} else {
			writeJSON(rootDir, task.JSONPath, structured, jsonResult)
		}
	}

	if jsonRaw {
		write(rootDir, task.JSONPath, data, jsonResult)
	}

	if humanTemplated {
		writeTemplate(rootDir, task.HumanPath, task.Template, structured, humanResult)
	}

	if humanYAML {
		writeYAML(rootDir, task.HumanPath, structured, humanResult)
	}

	if humanRaw {
		write(rootDir, task.HumanPath, data, humanResult)
	}

	return results
}
