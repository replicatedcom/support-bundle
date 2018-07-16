package plans

import (
	"context"
	"io"
	"time"

	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

// StreamSource is a Task that gets its data as an io.Reader
type StreamSource struct {
	Spec types.Spec

	// Producer provides the seed data for this task as an io.Reader
	Producer func(context.Context) (io.Reader, error)
	// StreamFormat describe stream format returned by Producer.  Only "" and "tar" are supported.
	StreamFormat string
	// RawScrubber, if defined, rewrites the raw data to to remove sensitive data
	RawScrubber func([]byte) []byte
	// Parser, if defined, structures the raw data for json and human sinks
	Parser func(io.Reader) (interface{}, error)
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

func (task *StreamSource) GetSpec() types.Spec {
	return task.Spec
}

func (task *StreamSource) Exec(ctx context.Context, rootDir string) []*types.Result {
	s := StreamsSource{
		Spec:         task.Spec,
		StreamFormat: task.StreamFormat,
		RawScrubber:  task.RawScrubber,
		Parser:       task.Parser,
		Template:     task.Template,
		RawPath:      task.RawPath,
		JSONPath:     task.JSONPath,
		HumanPath:    task.HumanPath,
		Timeout:      task.Timeout,
	}
	s.Producer = func(context.Context) (map[string]io.Reader, error) {
		r, err := task.Producer(ctx)
		if err != nil {
			return nil, err
		}
		return map[string]io.Reader{"": r}, nil
	}
	return s.Exec(ctx, rootDir)
}
