package plans

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/pkg/errors"

	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

// StructuredSource is a Task that gets its data as a structured object ready to
// be jsonified or used as a context in a human template.
type StructuredSource struct {
	Spec types.Spec

	// Producer provides the seed data for this task
	Producer func(context.Context) (interface{}, error)
	// RawScrubber, if defined, rewrites the raw data to to remove sensitive data
	RawScrubber func([]byte) []byte
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

func (task *StructuredSource) GetSpec() types.Spec {
	return task.Spec
}

func (task *StructuredSource) Exec(ctx context.Context, rootDir string) []*types.Result {
	s := StreamsSource{
		Spec:        task.Spec,
		RawScrubber: task.RawScrubber,
		Template:    task.Template,
		RawPath:     task.RawPath,
		JSONPath:    task.JSONPath,
		HumanPath:   task.HumanPath,
		Timeout:     task.Timeout,
	}
	if task.Producer != nil {
		structured, structuredErr := task.Producer(ctx)
		s.Parser = func(r io.Reader) (interface{}, error) {
			dec := json.NewDecoder(r)
			var structured interface{}
			return structured, dec.Decode(&structured)
		}
		s.Producer = func(context.Context) (map[string]io.Reader, error) {
			if structuredErr != nil {
				return nil, structuredErr
			}
			data, err := json.Marshal(structured)
			if err != nil {
				return nil, errors.Wrap(err, "json.Marshal")
			}
			return map[string]io.Reader{"": bytes.NewReader(data)}, nil
		}
	}
	return s.Exec(ctx, rootDir)
}
