package pipeline

import (
	"context"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

type Jsonifier struct {
}

func (p *Jsonifier) Pipe(ctx context.Context, i interface{}) (interface{}, error) {
	if s, ok := i.(StructuredSource); ok {
		return p.pipe(ctx, s)
	}
	return nil, errors.New("Jsonifier requires StructuredSource")
}

func (p *Jsonifier) pipe(ctx context.Context, s StructuredSource) (JSONSource, error) {
	pr, pw := io.Pipe()
	go func() {
		enc := json.NewEncoder(pw)
		enc.SetIndent("", "  ")
		pw.CloseWithError(enc.Encode(s))
	}()
	return JSONSource(pr), nil
}
