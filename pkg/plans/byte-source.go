package plans

import (
	"bytes"
	"context"
	"io"
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
	Parser func(io.Reader) (interface{}, error)
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
	s := StreamsSource{
		RawScrubber: task.RawScrubber,
		Parser:      task.Parser,
		Template:    task.Template,
		RawPath:     task.RawPath,
		JSONPath:    task.JSONPath,
		HumanPath:   task.HumanPath,
		Timeout:     task.Timeout,
	}
	if task.Producer != nil {
		s.Producer = func(context.Context) (map[string]io.Reader, error) {
			data, err := task.Producer(ctx)
			if err != nil {
				return nil, err
			}
			return map[string]io.Reader{"": bytes.NewReader(data)}, nil
		}
	}
	return s.Exec(ctx, rootDir)
}
